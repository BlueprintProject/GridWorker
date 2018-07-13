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

type remoteWorker struct {
	connectionQueue   chan *connection
	distributedWorker *DistributedWorker

	inQueue             bool
	outstandingMessages int
}

func (d *DistributedWorker) newRemoteWorker() *remoteWorker {
	r := remoteWorker{}
	r.distributedWorker = d
	r.connectionQueue = make(chan *connection, maxConnections)
	return &r
}

// getConnection gets a connection out of the connection queue
func (r *remoteWorker) getConnection() *connection {
	c := <-r.connectionQueue
	r.connectionQueue <- c

	return c
}

// send will send a message object to a connection in the pool
func (r *remoteWorker) send(m *Message) {
	isCmd := !m.isCtrl()

	r.getConnection().writeChan <- m

	if isCmd {
		r.outstandingMessages++

		if !r.inQueue && r.outstandingMessages < maximumOutstandingMessages {
			r.inQueue = true
			r.distributedWorker.remoteWorkerQueue <- r
		}
	}
}
