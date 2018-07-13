package gridworker

// refillPoolFunc is the function that is run to create a new object
type refillPoolFunc func() interface{}

// refillPool is used where object reuse isn't possible, but allocs
// still need to be avoided. It essentially creates a buffer of objects
// that can readily accessed and then creates a new object in a background
// go routine to replenish the buffer
type refillPool struct {
	// pool is a channel that acts as the object buffer
	pool chan interface{}

	// refillChan is the channel that is used to tell the
	// pool that an new object is needed
	refillChan chan bool

	// limit tells the pool how many objects need to exist in the buffer
	limit int

	// newObject is the function that is run to create a new object
	newObject refillPoolFunc
}

// newRefillPool creates a new pool
func newRefillPool(limit int, n refillPoolFunc) *refillPool {
	pool := &refillPool{
		pool:       make(chan interface{}, limit),
		refillChan: make(chan bool, 0),
		limit:      limit,
		newObject:  n,
	}

	pool.fill()
	go pool.refillLoop()

	return pool
}

// Get takes the object from the pool.
func (p *refillPool) get() interface{} {
	var c interface{}

	select {
	case c = <-p.pool:
		p.replenish()
	default:
		c = p.newObject()
	}

	return c
}

// fill will fill the pool channel
func (p *refillPool) fill() {
	for i := len(p.pool); i < p.limit; i++ {
		p.put(p.newObject())
	}
}

// refill loops on the refillChan and will call the fill function
// when it is looped
func (p *refillPool) refillLoop() {
	for <-p.refillChan {
		p.fill()
	}
}

// replenish will send a message to the refillChan which refillLoop
// is watching. However if refillLoop is already executing a fill
// it will not wait and just drop the request
func (p *refillPool) replenish() {
	select {
	case p.refillChan <- true:
	default:
	}
}

// put returns the object to the pool.
func (p *refillPool) put(c interface{}) {
	select {
	case p.pool <- c:
	default:
	}
}
