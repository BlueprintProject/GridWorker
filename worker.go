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

// Worker is a base object that is able to do one task at a time
type Worker struct {
	taskRegistry map[string]*Task
	processPool  *processPool

	TaskQueue chan *Context
}

// NewWorker creates a new worker object
func NewWorker() *Worker {
	worker := Worker{}
	worker.taskRegistry = map[string]*Task{}
	worker.processPool = newProcessPool()
	worker.TaskQueue = make(chan *Context, maxTasks)

	go worker.listenToTaskQueue()

	return &worker
}

// NewMessage creates a new message object to be sent to the worker
func (w *Worker) NewMessage() *Message {
	message := w.processPool.messagePool.NewMessage()
	message.arguments = w.processPool.messagePool.getMessageContentItem()
	return message
}

// Run adds the passed message to the task queue
func (w *Worker) Run(m *Message) *Reciept {
	c := w.processPool.contextPool.newContext()
	r := w.processPool.taskRecieptPool.newTaskReciept()

	c.input = m
	c.reciept = r

	w.TaskQueue <- c

	return r
}

// RegisterTask adds a task to the task registry so that it can be called
func (w *Worker) RegisterTask(t *Task) {
	w.taskRegistry[t.ID] = t
}

// listenToTaskQueue watches for new tasks and performs them in a loop
// only one task can be performed at a time for each worker
func (w *Worker) listenToTaskQueue() {
	for {
		c := <-w.TaskQueue

		task := w.taskRegistry[c.input.Command]

		go c.reciept.listenForResponse(c)
		task.Action(c)
	}
}
