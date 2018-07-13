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
)

// Reciept is the response object you get when a task
// is sent to the task queue. It is where you'll receieve
// the response
type Reciept struct {
	// Response is the channnel where your results will come through
	Response chan *Message

	// ReferenceID is the GUID assigned to the task
	ReferenceID string

	// closed determines if more messages will be coming on the
	// response queue
	closed bool

	// processPool is a reference to the shared process pool
	processPool *processPool

	// remoteWorker is a reference to the remoteWorker the task is
	// being performed on (if any)
	remoteWorker *remoteWorker
}

// taskRecieptPool is a pool used to quickly get TaskReciept objects
type taskRecieptPool struct {
	// pool is the actual pool
	pool *sync.Pool

	// processPool is a referecne to the shared process pool
	processPool *processPool
}

// newTaskRecieptPool creates a new task reciept pool
func (p *processPool) newTaskRecieptPool() {
	p.taskRecieptPool = &taskRecieptPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return &Reciept{
					Response:    make(chan *Message, 1),
					processPool: p,
				}
			},
		},
		processPool: p,
	}
}

// newTaskReciept gets a fresh task reciept out of the pool
func (t *taskRecieptPool) newTaskReciept() *Reciept {
	r := t.pool.Get().(*Reciept)
	r.ReferenceID = t.processPool.guidPool.get().(string)
	return r
}

// release returns the task receipt to the pool
func (t *Reciept) release() {
	if !t.closed {
		m := t.processPool.messagePool.NewMessage()
		t.Response <- m
	}

	t.closed = false
	t.processPool.taskRecieptPool.pool.Put(t)
}

// listenForResponse locks on to a Context object and waits
// for responses to come through. It will then release the reciept
// back to the pool if the response message is the last one
func (t *Reciept) listenForResponse(c *Context) {
	for {
		m := <-c.Response
		t.Response <- m

		if m.Done {
			break
		}
	}

	t.closed = true
	c.release()
	t.release()
}
