package upgradeagent

type dockerComposeDeployment struct {
	upgradeInfo    *UpstreamResponseUpgradeInfo
	upstreamClient UpstreamClient
}

// NewDockerComposeDeployment creates a new instance of docker-compose deployment kind.
func NewDockerComposeDeployment(upstreamServiceClient UpstreamClient) DeploymentClient {
	return &dockerComposeDeployment{upgradeInfo: NewUpstreamResponseUpgradeInfo(), upstreamClient: upstreamServiceClient}
}

// UpgradeCheck do check if there is a new version of the current solution available
func (d *dockerComposeDeployment) UpgradeCheck(localSolutionInfo LocalSolutionInfo) bool {
	if upgradeInfo, isNewVersion := d.upstreamClient.RequestUpgrade(localSolutionInfo); isNewVersion {
		d.upgradeInfo = upgradeInfo
		return true
	}
	return false
}

// Upgrade does upgrade the current solution
func (d *dockerComposeDeployment) Upgrade() error {
	return nil
}

// Rollback does rollback the current solution to the previous state in case both upgrade or health check is fails
func (d *dockerComposeDeployment) Rollback() {
	// TBD
}

// HealthCheck does a health check of the current solution after the upgrade
func (d *dockerComposeDeployment) HealthCheck() bool {
	d.upgradeInfo.HealthCheckStatus = Unhealthy
	return false
}
