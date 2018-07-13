package gridworker

// Base Command Names
const (
	cmdRSP  string = "$RESP$"
	cmdAUTH        = "$AUTH$"
)

// Message Constants

// The buffer size for the message pool
const messagePoolLimit = 25

// The buffer size for the message content pool
const messageContentPoolLimit = 25

// The size of the context pool
const contextPoolLimit = 25

// The size of the guid pool
const guidPoolLimit = 25

// maxTasks is the maximum number of the task queue
const maxTasks = 1024

// The maximum number of TCP connections that one worker may have
const maxConnections = 1024

// The maximum number of workers that a distributed worker may have
const maxWorkers = 1024

// The maximum number of outstanding messages allowed for a single worker
const maximumOutstandingMessages = 512
