package gridworker

// Context is used to pass information and give responses in
// worker tasks
type Context struct {
	// Response is a chan that the worker responds with
	Response chan *Message

	// input is the base input Message for the worker
	input *Message

	// reciept is the task receipt for the context
	reciept *Reciept

	// closed is used to determine if the task has sent it's final message
	closed bool

	// process pool is a reference to the process pool
	processPool *processPool
}

// contextPool is a struct that allows us to get context objects
// out of a pre alloced pool rather than reallocing every time
type contextPool struct {
	// pool is the actual context pool
	pool *limitPool
}

// newContextPool creates a new contextPool for the processPool
func (p *processPool) newContextPool() {
	p.contextPool = &contextPool{
		pool: newLimitPool(contextPoolLimit, func() interface{} {
			return &Context{
				Response:    make(chan *Message, 1),
				processPool: p,
			}
		}),
	}
}

// newContext gets a new context object from the context pool
func (c *contextPool) newContext() *Context {
	return c.pool.get().(*Context)
}

// release returns the context to the context pool
func (c *Context) release() {
	c.input = nil
	c.closed = false
	c.processPool.contextPool.pool.put(c)
}
