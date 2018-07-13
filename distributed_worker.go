/*
 *   This file is part of GridWorker.
 *
 *   Copyright (c) 2018 Mocha Industries, LLC.
 *   All rights reserved.
 *
 *   GridWorker is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   GridWorker is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with GridWorker.  If not, see <https://www.gnu.org/licenses/>.
 */

package gridworker

import (
	"sync"
	"time"

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

	workerQueue       chan *Worker
	remoteWorkerQueue chan *remoteWorker

	remoteWorkerRegistry map[string]*remoteWorker

	upgrader websocket.Upgrader
}

// NewDistributedWorker creates a new distributed worker object
// without any child workers
func NewDistributedWorker() *DistributedWorker {
	w := DistributedWorker{}
	w.processPool = newProcessPool()
	w.id = w.processPool.guidPool.get().(string)
	w.refMap = &sync.Map{}
	w.remoteWorkerRegistry = map[string]*remoteWorker{}
	w.workerQueue = make(chan *Worker, maxWorkers)
	w.remoteWorkerQueue = make(chan *remoteWorker, maxWorkers)
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
	rw := d.getRemoteWorker()
	r.remoteWorker = rw

	rw.send(m)

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

// getLocalWorker gets a worker for the workerQueue
func (d *DistributedWorker) getLocalWorker() *Worker {
	return <-d.workerQueue
}

// getRemoteWorker gets a remote worker for the workerQueue
func (d *DistributedWorker) getRemoteWorker() *remoteWorker {
	r := <-d.remoteWorkerQueue
	r.inQueue = false
	return r
}

// send will send a message to a remote worker
func (d *DistributedWorker) send(m *Message) {
	d.getRemoteWorker().send(m)
}

// listenToReceipt waits for a response from a worker reciept and sends it to
// the requester's connection
func (d *DistributedWorker) listenToReceipt(r *Reciept, localWorker *Worker) {
	for {
		m := <-r.Response
		r.remoteWorker.send(m)

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

	d.handleMessage(m, c.remoteWorker)
}

// handleMessage processes and dispatches messages from the sockets
func (d *DistributedWorker) handleMessage(m *Message, rem *remoteWorker) {
	if m.Command == cmdRSP {
		rem.outstandingMessages--
		if !rem.inQueue && rem.outstandingMessages < maximumOutstandingMessages {
			rem.inQueue = true
			d.remoteWorkerQueue <- rem
		}

		rI, ok := d.refMap.Load(m.ReferenceID)

		if !ok {
			return
		}

		r := rI.(*Reciept)

		r.Response <- m

		if m.Done {
			d.refMap.Delete(m.ReferenceID)
		}

		return
	}

	if d.hasLocalWorkers == true {
		localWorker := d.getLocalWorker()
		r := localWorker.Run(m)
		r.remoteWorker = rem
		d.listenToReceipt(r, localWorker)
	}
}

// handleAuthMessageForConnection handles actual auth commands
func (d *DistributedWorker) handleAuthMessageForConnection(m *Message, c *connection) {
	id := m.GetString("id")
	r := d.remoteWorkerRegistry[id]

	if r == nil {
		r = d.newRemoteWorker()
		d.remoteWorkerRegistry[id] = r
	}

	r.connectionQueue <- c
	c.remoteWorker = r
	c.announceAuth()

	go func(d *DistributedWorker, r *remoteWorker) {
		time.Sleep(1 * time.Second)
		d.remoteWorkerQueue <- r
	}(d, r)
}
