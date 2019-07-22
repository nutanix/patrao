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
