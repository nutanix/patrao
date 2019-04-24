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
