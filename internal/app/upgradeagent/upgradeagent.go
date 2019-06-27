package upgradeagent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

var (
	client DockerClient
)

// Main - Upgrade Agent entry point
func Main(context *cli.Context) error {
	client = NewClient(false)

	if context.GlobalBool(RunOnceName) {
		return runOnce(context)
	}
	return schedulePeriodicUpgrades(context)
}

func runOnce(context *cli.Context) error {
	log.Infoln("[+]runOnce()")
	containers, rc := client.ListContainers()
	if nil != rc {
		log.Error(rc)
		return rc
	}
	if len(containers) == 0 {
		log.Info("There are no launched containers on the host")
		return nil
	}
	upgradeInfoArray, rc := getLaunchedSolutionsList(context, &containers)
	if nil != rc {
		log.Error(rc)
		return rc
	}
	rc = doUpgradeSolutions(&upgradeInfoArray, &containers)
	log.Infoln("[-]runOnce()")
	return rc
}

func schedulePeriodicUpgrades(context *cli.Context) error {
	log.Infoln("[+]schedulePeriodicUpgrades()")
	{
		gocron.Every(uint64(context.GlobalInt64(UpgradeIntervalName))).Seconds().Do(runOnce, context)
		<-gocron.Start()
	}
	log.Infoln("[-]schedulePeriodicUpgrades()")
	return nil
}

func getLaunchedSolutionsList(context *cli.Context, containers *[]Container) ([]UpstreamResponseUpgradeInfo, error) {
	var (
		solutions []string
		rc        []UpstreamResponseUpgradeInfo
		err       error
	)
	getURLPath := context.GlobalString(UpstreamName) + UpstreamGetUpgrade

	for _, current := range *containers {
		currentSolution, _, err := GetSolutionAndServiceName(current.Name())
		if nil != err {
			log.Error(err)
			continue
		}
		if Contains(&solutions, &currentSolution) {
			continue
		}
		solutions = append(solutions, currentSolution)
		currentPath := getURLPath + currentSolution
		resp, err := http.Get(currentPath)
		if nil != err {
			log.Error(err)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if nil != err {
			log.Error(err)
			continue
		}

		toUpgrade := NewUpstreamResponseUpgradeInfo()
		err = json.Unmarshal([]byte(body), &toUpgrade)

		if nil != err {
			log.Error(err)
			continue
		}

		rc = append(rc, *toUpgrade)
	}
	return rc, err
}

func doUpgradeSolutions(upgradeInfoList *[]UpstreamResponseUpgradeInfo, containers *[]Container) error {
	var (
		rc      error
		toCheck []UpstreamResponseUpgradeInfo
	)

	for _, item := range *upgradeInfoList {
		if !isNewVersion(&item, containers) {
			log.Infof("Solution [%s] is up-to-date.", item.Name)
			continue
		}
		for _, container := range *containers {
			name, _, _ := GetSolutionAndServiceName(container.Name())
			if item.Name == name {
				err := client.StopContainer(container, DefaultTimeoutS*time.Second)
				if nil != err {
					log.Error(err)
				}
			}
		}
		err := client.LaunchSolution(&item)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Solution [%s] is successful launched", item.Name)
		toCheck = append(toCheck, item)
	}
	doHealthChek(&toCheck)

	return rc
}

// isNewVersion check if there are new version available.
func isNewVersion(upgradeInfo *UpstreamResponseUpgradeInfo, containers *[]Container) bool {
	var rc bool

	rc = false
	containersSpec := make(map[string]ContainerSpec)
	specMap := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(upgradeInfo.Spec), specMap)
	if nil != err {
		log.Error(err)
		return false
	}
	if val, exists := specMap[DockerComposeServicesName]; exists {
		servicesMap := val.(map[interface{}]interface{})
		for service := range servicesMap {
			var (
				serviceImageName string
				found            bool
				ok               bool
			)
			details := servicesMap[service].(map[interface{}]interface{})
			found = false
			for item := range details {
				if DockerComposeImageName == item {
					val, itemFound := details[item]
					if itemFound {
						serviceImageName, ok = val.(string)
						if ok {
							found = true
							break
						}
					}
				}
			}
			if found {

				containersSpec[service.(string)] = ContainerSpec{Name: service.(string), Image: serviceImageName}
			}
		}
		for _, container := range *containers {
			solutionName, serviceName, _ := GetSolutionAndServiceName(container.Name())
			val, exist := containersSpec[serviceName]
			if (exist) && (upgradeInfo.Name == solutionName) {
				if container.ImageName() != val.Image {
					rc = true
					break
				}
			}
		}
	}
	return rc
}

func doContainersCheck(solutionInfo *UpstreamResponseUpgradeInfo) {
	timeout := time.After(time.Duration(solutionInfo.ThresholdTimeS) * time.Second)

	for solutionInfo.HealthCheckStatus == Undefined {
		select {
		case <-timeout:
			solutionInfo.HealthCheckStatus = Unhealthy
			log.Infof("Solution [%s] is Unhealthy", solutionInfo.Name)
			return
		default:
			{
				checkContainersCompletedCount := 0

				for _, healthChekCmd := range solutionInfo.HealthCheckCmds {
					container, err := client.GetContainerByName(solutionInfo.Name, healthChekCmd.ContainerName)
					if err != nil {
						log.Error(err)
						break
					}
					config, err := client.InspectContainer(container)
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
					exitCode, err := client.ExecContainer(container, healthChekCmd.Cmd)
					if err != nil {
						log.Error(err)
						break
					}
					if exitCode == 0 {
						log.Infof("Container %s has passed healthcheck command [%s], exit code is [%d]", healthChekCmd.ContainerName, healthChekCmd.Cmd, exitCode)
						checkContainersCompletedCount++
					}
				}
				if checkContainersCompletedCount == len(solutionInfo.HealthCheckCmds) {
					log.Infof("Solution [%s] is healthy", solutionInfo.Name)
					solutionInfo.HealthCheckStatus = Healthy
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// doHealthCheck do solutions healthcheck afeter upgrade is completed
func doHealthChek(toCheck *[]UpstreamResponseUpgradeInfo) {
	for _, item := range *toCheck {
		doContainersCheck(&item)
	}
}
