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

// Base Command Names
const (
	cmdRSP  string = "$RESP$"
	cmdAUTH        = "$AUTH$"
)

// Message Constants

// The buffer size for the message pool
const messagePoolLimit = 25

// The buffer size for the message content pool
const messageContentPoolLimit = 25

// The size of the context pool
const contextPoolLimit = 25

// The size of the guid pool
const guidPoolLimit = 25

// maxTasks is the maximum number of the task queue
const maxTasks = 1024

// The maximum number of TCP connections that one worker may have
const maxConnections = 1024

// The maximum number of workers that a distributed worker may have
const maxWorkers = 1024

// The maximum number of outstanding messages allowed for a single worker
const maximumOutstandingMessages = 512
