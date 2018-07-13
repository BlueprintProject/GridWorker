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
