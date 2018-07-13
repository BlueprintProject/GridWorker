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
