package upgradeagent

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type dockerComposeDeployment struct {
	upgradeInfo       *UpstreamResponseUpgradeInfo
	upstreamClient    UpstreamClient
	context           *cli.Context
	dockerClient      DockerClient
	localSolutionInfo *LocalSolutionInfo
}

// NewDockerComposeDeployment creates a new instance of docker-compose deployment kind.
func NewDockerComposeDeployment(ctx *cli.Context, upstreamServiceClient UpstreamClient, dockerCli DockerClient, solutionInfo *LocalSolutionInfo) DeploymentClient {
	return &dockerComposeDeployment{
		upgradeInfo:       NewUpstreamResponseUpgradeInfo(),
		upstreamClient:    upstreamServiceClient,
		context:           ctx,
		dockerClient:      dockerCli,
		localSolutionInfo: solutionInfo,
	}
}

// CheckUpgrade do check if there is a new version of the current solution available
func (d *dockerComposeDeployment) CheckUpgrade() bool {
	var (
		upgradeInfo *UpstreamResponseUpgradeInfo
		err         error
		containers  []Container
	)
	if upgradeInfo, err = d.upstreamClient.RequestUpgrade(*d.localSolutionInfo); err != nil {
		log.Error(err)
		return false
	}
	if containers, err = d.dockerClient.ListContainers(); err != nil {
		log.Error(err)
		return false
	}
	if isNewVersion(upgradeInfo, containers) == false {
		log.Infof("Solution [%s] is up-to-date.", upgradeInfo.Name)
		return false
	}
	d.upgradeInfo = upgradeInfo
	return true
}

// Upgrade does upgrade the current solution
func (d *dockerComposeDeployment) DoUpgrade() error {
	var (
		containers []Container
		err        error
	)
	if containers, err = d.dockerClient.ListContainers(); err != nil {
		return fmt.Errorf("DockerComposeDeployment::DoUpgrade [%v]", err)
	}
	for _, container := range containers {
		if name, found := container.GetProjectName(); found && d.localSolutionInfo.name == name {
			if err = d.dockerClient.StopContainer(container, DefaultTimeoutS*time.Second); err != nil {
				return fmt.Errorf("DockerComposeDeployment::DoUpgrade [%v]", err)
			}
		}
	}
	if err = d.LaunchSolution(); err != nil {
		return fmt.Errorf("DockerComposeDeployment::DoUpgrade [%v]", err)
	}
	log.Infof("Solution [%s] is successful launched", d.localSolutionInfo.name)
	return nil
}

// CheckHealth does a health check of the current solution after the upgrade
func (d *dockerComposeDeployment) CheckHealth() bool {
	timeout := time.After(time.Duration(d.upgradeInfo.ThresholdTimeS) * time.Second)
	for d.upgradeInfo.HealthCheckStatus == Undefined {
		select {
		case <-timeout:
			d.upgradeInfo.HealthCheckStatus = Unhealthy
			log.Infof("Solution [%s] is Unhealthy", d.upgradeInfo.Name)
			return false
		default:
			{
				checkContainersCompletedCount := 0
				for _, healthChekCmd := range d.upgradeInfo.HealthCheckCmds {
					container, err := d.dockerClient.GetContainerByName(d.upgradeInfo.Name, healthChekCmd.ContainerName)
					if err != nil {
						log.Error(err)
						break
					}
					config, err := d.dockerClient.InspectContainer(container)
					if err != nil {
						log.Error(err)
						break
					}
					if !config.State.Running {
						log.Infof("Container %s is NOT Running state.", healthChekCmd.ContainerName)
						continue
					}
					if config.State.Health != nil {
						log.Infof("Container %s has embedded healthchek.", healthChekCmd.ContainerName)
						if config.State.Health.Status != types.Healthy {
							log.Infof("Container %s have healthy information. The current status is [%s]", healthChekCmd.ContainerName, config.State.Health.Status)
							continue
						}
					} else {
						log.Infof("Container %s has NOT embedded healthchek. Skip this step.", healthChekCmd.ContainerName)
					}
					exitCode, err := d.dockerClient.ExecContainer(container, healthChekCmd.Cmd)
					if err != nil {
						log.Error(err)
						break
					}
					if exitCode == 0 {
						log.Infof("Container %s has passed healthcheck command [%s], exit code is [%d]", healthChekCmd.ContainerName, healthChekCmd.Cmd, exitCode)
						checkContainersCompletedCount++
					}
				}
				if checkContainersCompletedCount == len(d.upgradeInfo.HealthCheckCmds) {
					log.Infof("Solution [%s] is healthy", d.upgradeInfo.Name)
					d.upgradeInfo.HealthCheckStatus = Healthy
					return true
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	return false
}

// GetLocalSolutionInfo returns pointer to LocalSolutionInfo data structure
func (d *dockerComposeDeployment) GetLocalSolutionInfo() *LocalSolutionInfo {
	return d.localSolutionInfo
}

// isNewVersion check if there are new version available.
func isNewVersion(upgradeInfo *UpstreamResponseUpgradeInfo, containers []Container) bool {
	var (
		found bool
	)
	containersSpec := make(map[string]ContainerSpec)
	specMap := make(map[string]interface{})
	rc := false
	if err := yaml.Unmarshal([]byte(upgradeInfo.Spec), specMap); err != nil {
		log.Error(err)
		return false
	}
	if val, exists := specMap[DockerComposeServicesName]; exists {
		servicesMap := val.(map[interface{}]interface{})
		for service := range servicesMap {
			var (
				serviceImageName string
				ok               bool
			)
			details := servicesMap[service].(map[interface{}]interface{})
			found = false
			for item := range details {
				if DockerComposeImageName == item {
					if val, itemFound := details[item]; itemFound {
						if serviceImageName, ok = val.(string); ok {
							found = true
							containersSpec[service.(string)] = ContainerSpec{Name: service.(string), Image: serviceImageName}
							break
						}
					}
				}
			}
		}
		for _, container := range containers {
			var (
				solutionName string
				serviceName  string
			)
			if solutionName, found = container.GetProjectName(); !found {
				continue
			}
			if serviceName, found = container.GetServiceName(); !found {
				continue
			}
			if val, exist := containersSpec[serviceName]; (exist) && (upgradeInfo.Name == solutionName) {
				if container.ImageName() != val.Image {
					rc = true
					break
				}
			}
		}
	}
	return rc
}

// LaunchSolution launch solution based on the received docker-compose specification
func (d *dockerComposeDeployment) LaunchSolution() error {
	if _, isFileExist := os.Stat(d.upgradeInfo.Name); !os.IsNotExist(isFileExist) {
		os.RemoveAll(d.upgradeInfo.Name)
	}

	if err := os.Mkdir(d.upgradeInfo.Name, os.ModePerm); err != nil {
		return fmt.Errorf("DockerComposeDeployment::LaunchSolution() [%v]", err)
	}
	defer os.Remove(d.upgradeInfo.Name)
	dockerComposeFileName := path.Join(d.upgradeInfo.Name, DockerComposeFileName)
	f, err := os.Create(dockerComposeFileName)
	if err != nil {
		return fmt.Errorf("DockerComposeDeployment::LaunchSolution() [%v]", err)
	}
	defer func() {
		f.Close()
		os.Remove(dockerComposeFileName)
	}()
	if _, err = f.Write([]byte(d.upgradeInfo.Spec)); err != nil {
		return fmt.Errorf("DockerComposeDeployment::LaunchSolution() [%v]", err)
	}
	log.Infof("Launching solution [%s]", d.upgradeInfo.Name)
	ex, _ := os.Executable()
	rootPath := filepath.Dir(ex)
	cmd := exec.Command(DockerComposeCommand, "-f", path.Join(rootPath, dockerComposeFileName), "up", "-d")
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("DockerComposeDeployment::LaunchSolution() [%v]", err)
	}
	return nil
}

// DoRollback does rollback the current solution to the previous state in case both upgrade or health check is fails
func (d *dockerComposeDeployment) DoRollback() {
	// TBD
}
