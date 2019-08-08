package upgradeagent

type kubernetesDeployment struct {
	upgradeInfo    *UpstreamResponseUpgradeInfo
	upstreamClient UpstreamClient
}

// NewKubernetesDeployment creates a new instance of kubernetes deployment kind.
func NewKubernetesDeployment(upstreamServiceClient UpstreamClient) DeploymentClient {
	return &kubernetesDeployment{upgradeInfo: NewUpstreamResponseUpgradeInfo(), upstreamClient: upstreamServiceClient}
}

// UpgradeCheck do check if there is a new version of the current solution available
func (d *kubernetesDeployment) UpgradeCheck(localSolutionInfo LocalSolutionInfo) bool {
	if upgradeInfo, isNewVersion := d.upstreamClient.RequestUpgrade(localSolutionInfo); isNewVersion {
		d.upgradeInfo = upgradeInfo
		return true
	}
	return false
}

// Upgrade does upgrade the current solution
func (d *kubernetesDeployment) Upgrade() error {
	return nil
}

// Rollback does rollback the current solution to the previous state in case both upgrade or health check is fails
func (d *kubernetesDeployment) Rollback() {
	// TBD
}

// HealthCheck does a health check of the current solution after the upgrade
func (d *kubernetesDeployment) HealthCheck() bool {
	d.upgradeInfo.HealthCheckStatus = Unhealthy
	return false
}
