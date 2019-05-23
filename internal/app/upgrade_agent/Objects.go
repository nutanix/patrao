package core

import (
	"fmt"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	log "github.com/sirupsen/logrus"
)

// KindTYPE Kind type for all structures
type KindTYPE string

// constants represent values of Kind field
const (
	NodeKind        KindTYPE = "node"
	AppTemplateKind          = "apptemplate"
	DeploymentKind           = "deployment"
)

// SubTYPE SubType data type
type SubTYPE string

// constants represent values for SubType
const (
	DockerCompose SubTYPE = "docker-compose"
)

// Node Identifies an individual VM uniquely, based on its node uuid.
type Node struct {
	Kind      KindTYPE
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeUUID  string
}

// AppTemplate A docker compose (or similar) spec with additional metadata that identifies a particular application by name.
// The combination of name (service or project name) and version is a unique identifier for an app template.
type AppTemplate struct {
	Kind      KindTYPE
	SubType   SubTYPE
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Spec      string
}

// Deployment  An instance of an AppTemplate deployed on a Node.
type Deployment struct {
	Kind      KindTYPE
	UUID      string
	CreatedAt time.Time
	UpdatedAt time.Time
	NodeUUID  string
}

/*
// CurrentContainerVersion structure present container infor for request to upstream api
type CurrentContainerVersion struct {
	ID      string
	CREATED string
	NAME    string
	IMAGE   string
}

// NewContainerVersion structure present information about new version of running container
type NewContainerVersion struct {
	ID            string
	NAME          string
	IMAGE         string
	DeleteVolumes bool `json:"DELETE_VOLUMES"`
}
*/

// UpstreamResponseUpgradeInfo structure represent response from Upstream Service
type UpstreamResponseUpgradeInfo struct {
	Name          string
	Spec          string
	DeleteVolumes string
}

// SolutionNameNotFound struct present error when agent couldn't find solution name by container name
type SolutionNameNotFound struct {
	When time.Time
	What string
}

func (e SolutionNameNotFound) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

// NewNode create and setup a new instance of Node structure
func NewNode() (obj *Node) {
	obj = new(Node)

	obj.Kind = NodeKind
	obj.UUID = genUUID()
	obj.NodeUUID = genNodeUUID()

	return
}

// NewAppTemplate create and setup a new instance of AppTemplate structure
func NewAppTemplate() (obj *AppTemplate) {
	obj = new(AppTemplate)

	obj.Kind = AppTemplateKind
	obj.UUID = genUUID()

	return
}

// NewDeployment create and setup a new instance of Deployment structure
func NewDeployment() (obj *Deployment) {
	obj = new(Deployment)

	obj.Kind = DeploymentKind
	obj.UUID = genUUID()
	obj.NodeUUID = genNodeUUID()

	return
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
