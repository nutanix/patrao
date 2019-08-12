package upgradeagent

import (
	"fmt"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Contains returns true in case if there is item in array. otherwise return false
func Contains(names *[]string, item *string) bool {
	for _, it := range *names {
		if it == *item {
			return true
		}
	}
	return false
}

// ParseLabels parse labels map to LocalSolutionInfo data structure
func ParseLabels(context *cli.Context, dockerClient DockerClient, labels map[string]string) (DeploymentClient, error) {
	if value, found := labels[DockerComposeProjectLabel]; found {
		info := NewLocalSolutionInfo()
		info.SetDeploymentKind(DockerComposeDeployment)
		info.SetName(value)
		info.AddServices(labels[DockerComposeServiceLabel])

		return NewDockerComposeDeployment(context, GetUpstreamClient(context), dockerClient, info), nil
	}
	return nil, fmt.Errorf("Cannot read labels [%s]", labels)
}

// GetLocalSolutionList return the list of running solutions
func GetLocalSolutionList(context *cli.Context, dockerClient DockerClient) map[string]DeploymentClient {
	projectMap := make(map[string]DeploymentClient)
	containers, err := dockerClient.ListContainers()
	if err != nil {
		log.Error(err)
		return projectMap
	}
	if containers != nil {
		for _, current := range containers {
			deploymentClient, err := ParseLabels(context, dockerClient, current.Labels())
			if err != nil {
				log.Error(err)
				continue
			}
			if _, ok := projectMap[deploymentClient.GetLocalSolutionInfo().GetName()]; ok {
				projectMap[deploymentClient.GetLocalSolutionInfo().GetName()].GetLocalSolutionInfo().AddServices(deploymentClient.GetLocalSolutionInfo().GetServices()...)
			} else {
				projectMap[deploymentClient.GetLocalSolutionInfo().GetName()] = deploymentClient
			}
		}
	}
	return projectMap
}

// GenUUID generate UUID string
func GenUUID() string {
	u4, _ := uuid.NewV4()
	return u4.String()
}

// GenNodeUUID generate node uuid
func GenNodeUUID() string {
	// TBD
	return "node-uuid"
}

// GetUpstreamClient returns apropriate upstream client instance depend on command line arguments
func GetUpstreamClient(context *cli.Context) UpstreamClient {
	switch upstreamType := context.GlobalString(UpstreamTypeName); upstreamType {
	case UpstreamTypeValue:
		return NewMockUpstreamClient(context)
		/*
			case UpstreamGitHub:
				return NewGitHubUpstreamClient(context)
			case UpstreamAmazonS3Bucket:
				return NewAmazonS3BucketClient(context)
				....
		*/
	default:
		log.Panicf("unsuported upstreamType [%s]", upstreamType)
	}

	return nil
}
