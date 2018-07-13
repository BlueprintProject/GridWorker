package gridworker

// TaskAction is the function that will be run when a task is executed
type TaskAction func(c *Context)

// Task is the base object that is registered in the task registry
// Every single command is a task
type Task struct {
	ID     string
	Action TaskAction
}
