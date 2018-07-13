package gridworker

import (
	uuid "github.com/nu7hatch/gouuid"
)

// processPool is the object that stores all of our pools
// It may be shared between many different objects and is
// completely concurrent safe
type processPool struct {
	messagePool     *messagePool
	contextPool     *contextPool
	taskRecieptPool *taskRecieptPool
	guidPool        *refillPool
}

// newProcessPool creates a newProcessPool
// This only should be called during startup as it is where all
// of the allocs are happening. If this is being done anywhere else
// we're not using allocs effectively and it is likely your bottleneck
func newProcessPool() *processPool {
	p := &processPool{}

	p.newMessagePool()
	p.newContextPool()
	p.newTaskRecieptPool()

	p.guidPool = newRefillPool(guidPoolLimit, func() interface{} {
		u, _ := uuid.NewV4()
		return u.String()
	})

	return p
}
