package upgradeagent

import (
	"fmt"
	"time"
)

// KindType Kind type for all structures
type KindType string

// constants represent values of Kind field
const (
	NodeKind        KindType = "node"
	AppTemplateKind          = "apptemplate"
	DeploymentKind           = "deployment"
)

// KindSubType SubType data type
type KindSubType string

// constants represent values for SubType
const (
	DockerCompose KindSubType = "docker-compose"
)

// HealthStatus indicates HealthCheck result
type HealthStatus string

// possible values of HealthStatus
const (
	Undefined = "undefined"
	Healthy   = "healthy"
	Unhealthy = "unhealthy"
)

// Node Identifies an individual VM uniquely, based on its node uuid.
type Node struct {
	Kind      KindType
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeUUID  string
}

// AppTemplate A docker compose (or similar) spec with additional metadata that identifies a particular application by name.
// The combination of name (service or project name) and version is a unique identifier for an app template.
type AppTemplate struct {
	Kind      KindType
	SubType   KindSubType
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Spec      string
}

// Deployment  An instance of an AppTemplate deployed on a Node.
type Deployment struct {
	Kind      KindType
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeUUID  string
}

// ContainerSpec Represents solution information from received from Docker-compose specification
type ContainerSpec struct {
	Name  string
	Image string
	//TBD
}

// UpstreamResponseUpgradeInfo structure represent response from Upstream Service
type UpstreamResponseUpgradeInfo struct {
	Name              string
	Spec              string
	DeleteVolumes     bool
	ThresholdTimeS    int
	HealthCheckStatus HealthStatus
	HealthCheckCmds   []struct {
		ContainerName string
		Cmd           string
	}
}

// LocalSolutionInfo represents information about solutions running on the host.
type LocalSolutionInfo struct {
	name           string
	services       []string
	deploymentType string
}

// AddService add service to string array
func (info *LocalSolutionInfo) AddService(serviceName string) {
	info.services = append(info.services, serviceName)
}

// AddServices add services array to string array
func (info *LocalSolutionInfo) AddServices(services []string) {

	for _, item := range services {
		info.AddService(item)
	}
}

// GetServices returns services related to running solution
func (info LocalSolutionInfo) GetServices() []string {
	return info.services
}

// GetName returns solution name
func (info LocalSolutionInfo) GetName() string {
	return info.name
}

// SetName set solition name
func (info *LocalSolutionInfo) SetName(solutionName string) {
	info.name = solutionName
}

// GetDeploymentKind return deployment kind for solution
func (info LocalSolutionInfo) GetDeploymentKind() string {
	return info.deploymentType
}

// SetDeploymentKind sets deployment kind fo solution
func (info *LocalSolutionInfo) SetDeploymentKind(deploymentKind string) {
	info.deploymentType = deploymentKind
}

// SolutionNameNotFound struct present error when agent couldn't find solution name by container name
type SolutionNameNotFound struct {
	When time.Time
	What string
}

func (e SolutionNameNotFound) Error() string {
	return fmt.Sprintf("%v at %v", e.When, e.What)
}

// NewNode create and setup a new instance of Node structure
func NewNode() *Node {
	return &Node{
		Kind:     NodeKind,
		UUID:     GenUUID(),
		NodeUUID: GenNodeUUID(),
	}
}

// NewAppTemplate create and setup a new instance of AppTemplate structure
func NewAppTemplate() *AppTemplate {
	return &AppTemplate{
		Kind: AppTemplateKind,
		UUID: GenUUID(),
	}
}

// NewDeployment create and setup a new instance of Deployment structure
func NewDeployment() *Deployment {
	return &Deployment{
		Kind: DeploymentKind,
		UUID: GenUUID(),
	}
}

// NewUpstreamResponseUpgradeInfo returns new instance of UpstreamResponseUpgradeInfo data structure
func NewUpstreamResponseUpgradeInfo() *UpstreamResponseUpgradeInfo {
	return &UpstreamResponseUpgradeInfo{
		HealthCheckStatus: Undefined,
	}
}

// NewLocalSolutionInfo returns new instance of LocalSolutionInfo data structure
func NewLocalSolutionInfo() *LocalSolutionInfo {
	return &LocalSolutionInfo{
		deploymentType: UndefinedDeployment,
	}
}
