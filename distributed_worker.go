package gridworker

import (
	"sync"

	"github.com/gorilla/websocket"
)

// DistributedWorker is used to represent a worker that
// is spread between multiple machines or processes
type DistributedWorker struct {
	id          string
	processPool *processPool
	refMap      *sync.Map

	hasLocalWorkers bool

	taskRegistry map[string]*Task

	workerQueue     chan *Worker
	connectionQueue chan *connection

	upgrader websocket.Upgrader
}

// NewDistributedWorker creates a new distributed worker object
// without any child workers
func NewDistributedWorker() *DistributedWorker {
	w := DistributedWorker{}
	w.processPool = newProcessPool()
	w.id = w.processPool.guidPool.get().(string)
	w.refMap = &sync.Map{}
	w.workerQueue = make(chan *Worker, maxWorkers)
	w.connectionQueue = make(chan *connection, maxConnections)
	w.taskRegistry = map[string]*Task{}
	w.upgrader = websocket.Upgrader{}
	return &w
}

// NewDistributedWorkerWithLocalWorkers creates a new distributed worker
// object with a set number of child workers. The child workers perform
// tasks locally.
func NewDistributedWorkerWithLocalWorkers(numberOfLocalWorkers int) *DistributedWorker {
	w := NewDistributedWorker()

	for i := 0; i < numberOfLocalWorkers; i++ {
		w.workerQueue <- w.newWorker()
	}

	w.hasLocalWorkers = true

	return w
}

// NewMessage creates a new message to send to the worker
func (d DistributedWorker) NewMessage() *Message {
	message := d.processPool.messagePool.NewMessage()
	message.Arguments = d.processPool.messagePool.getMessageContentItem()
	return message
}

// Run sends a message to the DistributedWorker task queue
func (d DistributedWorker) Run(m *Message) *Reciept {
	r := d.processPool.taskRecieptPool.newTaskReciept()

	m.ReferenceID = r.ReferenceID

	d.registerTaskRecieptForRef(r, m.ReferenceID)
	d.send(m)

	return r
}

// RegisterTask adds a task to the task registry, to allow it to be run
func (d *DistributedWorker) RegisterTask(t *Task) {
	d.taskRegistry[t.ID] = t
}

// newWorker creates a new local (non distrubted) worker to perform work
func (d *DistributedWorker) newWorker() *Worker {
	w := Worker{}
	w.taskRegistry = d.taskRegistry
	w.processPool = d.processPool

	w.TaskQueue = make(chan *Context, maxTasks)

	go w.listenToTaskQueue()

	return &w
}

// getConnection gets a connection out of the connection queue
func (d *DistributedWorker) getConnection() *connection {
	c := <-d.connectionQueue
	d.connectionQueue <- c

	return c
}

// send will send a message object to a connection in the pool
func (d *DistributedWorker) send(m *Message) {
	d.getConnection().writeChan <- m
}

// getLocalWorker gets a worker for the workerQueue
func (d *DistributedWorker) getLocalWorker() *Worker {
	return <-d.workerQueue
}

// listenToReceipt waits for a response from a worker reciept and sends it to
// the requester's connection
func (d *DistributedWorker) listenToReceipt(r *Reciept, localWorker *Worker) {
	for {
		m := <-r.Response
		d.send(m)

		if m.Done {
			break
		}
	}

	d.workerQueue <- localWorker
}

// registerTaskRecieptForRef stores a ref in the refmap so that it's response channel
// can be reassociated once the task is complete
func (d *DistributedWorker) registerTaskRecieptForRef(r *Reciept, ref string) {
	d.refMap.Store(ref, r)
}

// handleMessageForConnection processes messages from connections and checks if it is
// the auth command. In which case it will handle it differently than regular messages
// otherwise it gets sent to the handleMessage function
func (d *DistributedWorker) handleMessageForConnection(m *Message, c *connection) {
	if m.Command == cmdAUTH {
		d.handleAuthMessageForConnection(m, c)
		return
	}

	d.handleMessage(m)
}

// handleMessage processes and dispatches messages from the sockets
func (d *DistributedWorker) handleMessage(m *Message) {
	if m.Command == cmdRSP {
		r, ok := d.refMap.Load(m.ReferenceID)

		if ok {
			r.(*Reciept).Response <- m
			d.refMap.Delete(m.ReferenceID)
		}

		return
	}

	if d.hasLocalWorkers == true {
		localWorker := d.getLocalWorker()
		r := localWorker.Run(m)
		d.listenToReceipt(r, localWorker)
	}
}

// handleAuthMessageForConnection handles actual auth commands
func (d *DistributedWorker) handleAuthMessageForConnection(m *Message, c *connection) {
	d.id = m.GetString("id")
	d.connectionQueue <- c
	c.announceAuth()
}
