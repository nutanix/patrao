package core

// HostName is command line parameter name
const HostName = "host"

// HostUsage is description of HostName
const HostUsage = "daemon socket to connect to docker"

// HostValue is default value for HostName
const HostValue = "unix:///var/run/docker.sock"

// HostEnvVar is env variable name for HostName
const HostEnvVar = "PATRAO_DOCKER_HOST"

// RunOnceName is name for bool command line parameter
const RunOnceName = "run-once"

// RunOnceUsage is description of RunOnceName
const RunOnceUsage = "Run once now and exit"

// DefaultStopSignal is default stop signal
const DefaultStopSignal = "SIGTERM"

// PatraoAgentContainerName is container name of Patrao Agent
const PatraoAgentContainerName = "/patrao_agent"

// UpstreamName is command line parameter name
const UpstreamName = "upstreamHost"

// UpstreamUsage is description of UpstreamName
const UpstreamUsage = "upstream host name"

// UpstreamValue is default value for UpstreamName
const UpstreamValue = "http://localhost:1080"

// UpstreamEnvVar is env variable for UpstreamName
const UpstreamEnvVar = "PATRAO_UPSTREAM_HOST"

// DockerComposeFileName is defaut file name for *.yml scripts
const DockerComposeFileName = "docker-compose.yml"

// DockerComposeCommand is command name of compose service
const DockerComposeCommand = "docker-compose"

// UpstreamGetUpgrade is template for http get request to Upstream Service
const UpstreamGetUpgrade = "/v1/node/test_node_id/request/upgrade/"

// UpgradeIntervalName is command line parameter name
const UpgradeIntervalName = "upgradeInterval"

// UpgradeIntervalUsage is description of UpgradeIntervalName
const UpgradeIntervalUsage = "upgrade interval in seconds (default is 1 hour)"

// UpgradeIntervalValue is default value for UpgradeIntervalName in seconds (3600 seconds is 1 hour default value)
const UpgradeIntervalValue = "3600"

// UpgradeIntervalValueEnvVar is env variable for UpgradeIntervalName
const UpgradeIntervalValueEnvVar = "PATRAO_UPGRADE_INTERVAL_S"

// DockerComposeImageName image section name. related to docker-compose.yml
const DockerComposeImageName = "image"

// DockerComposeServicesName services section name. related to docker-compose.yml
const DockerComposeServicesName = "services"
