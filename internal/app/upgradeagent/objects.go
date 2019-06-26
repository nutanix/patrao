package upgradeagent

import (
	"fmt"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
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
	Name            string
	Spec            string
	DeleteVolumes   string
	ThresholdTimeS  int
	HealthCheckCmds []struct {
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
		UUID:     genUUID(),
		NodeUUID: genNodeUUID(),
	}
}

// NewAppTemplate create and setup a new instance of AppTemplate structure
func NewAppTemplate() *AppTemplate {
	return &AppTemplate{
		Kind: AppTemplateKind,
		UUID: genUUID(),
	}
}

// NewDeployment create and setup a new instance of Deployment structure
func NewDeployment() *Deployment {
	return &Deployment{
		Kind: DeploymentKind,
		UUID: genUUID(),
	}
}

func genUUID() string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return u4.String()
}

func genNodeUUID() string {
	// TBD
	return "node-uuid"
}
