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

// UpstreamGetUpgrade is template for http get request to Upstream Service
const UpstreamGetUpgrade = "/v1/node/test_node_id/request/upgrade/"

// UpstreamEnvVar is env variable for UpstreamName
const UpstreamEnvVar = "PATRAO_UPSTREAM_HOST"
