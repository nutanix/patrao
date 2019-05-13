package core

import (
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
