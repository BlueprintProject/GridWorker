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
