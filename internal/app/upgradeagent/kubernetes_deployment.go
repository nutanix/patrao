package upgradeagent

import (
	"github.com/urfave/cli"
)

type kubernetesDeployment struct {
	upgradeInfo       *UpstreamResponseUpgradeInfo
	upstreamClient    UpstreamClient
	context           *cli.Context
	dockerClient      DockerClient
	localSolutionInfo *LocalSolutionInfo
}

// NewKubernetesDeployment creates a new instance of kubernetes deployment kind.
func NewKubernetesDeployment(ctx *cli.Context, upstreamServiceClient UpstreamClient, dockerCli DockerClient, solutionInfo *LocalSolutionInfo) DeploymentClient {
	return &kubernetesDeployment{
		upgradeInfo:       NewUpstreamResponseUpgradeInfo(),
		upstreamClient:    upstreamServiceClient,
		context:           ctx,
		dockerClient:      dockerCli,
		localSolutionInfo: solutionInfo,
	}
}

// GetLocalSolutionInfo returns pointer to LocalSolutionInfo data structure
func (d *kubernetesDeployment) GetLocalSolutionInfo() *LocalSolutionInfo {
	return d.localSolutionInfo
}

// CheckUpgrade do check if there is a new version of the current solution available
func (d *kubernetesDeployment) CheckUpgrade() bool {
	return false
}

// DoUpgrade does upgrade the current solution
func (d *kubernetesDeployment) DoUpgrade() error {
	return nil
}

// DoRollback does rollback the current solution to the previous state in case both upgrade or health check is fails
func (d *kubernetesDeployment) DoRollback() {
	// TBD
}

// CheckHealth does a health check of the current solution after the upgrade
func (d *kubernetesDeployment) CheckHealth() bool {
	d.upgradeInfo.HealthCheckStatus = Unhealthy
	return false
}

// LaunchSolution starts solution using UpstreamResponseUpgradeInfo data structure
func (d *kubernetesDeployment) LaunchSolution() error {
	return nil
}
