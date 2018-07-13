package gridworker

type limitPoolFunc func() interface{}

type limitPool struct {
	pool      chan interface{}
	limit     int
	newObject limitPoolFunc
}

// newLimitPool creates a new pool
func newLimitPool(limit int, n limitPoolFunc) *limitPool {
	pool := &limitPool{
		pool:      make(chan interface{}, limit),
		limit:     limit,
		newObject: n,
	}

	pool.fill()

	return pool
}

// Get takes the object from the pool.
func (p *limitPool) get() interface{} {
	return <-p.pool
}

// Fills the pool channel
func (p *limitPool) fill() {
	if len(p.pool) == 0 {
		for i := 0; i < p.limit; i++ {
			p.put(p.newObject())
		}
	}
}

// Put returns the object to the pool.
func (p *limitPool) put(c interface{}) {
	select {
	case p.pool <- c:
	default:
	}
}
