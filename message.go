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

import "github.com/golang/protobuf/proto"

// messageContent is the base object inside of a message
type messageContent map[string]*DynamicType

// Message is the object that communications are sent through
type Message struct {
	Command string

	arguments map[string]*DynamicType

	done        bool
	referenceID string

	messagePool *messagePool
}

// isCtrl determines if a message is a command message or a standard control message
func (m *Message) isCtrl() bool {
	return m.Command == cmdAUTH || m.Command == cmdRSP
}

// messagePool is the general data representation for the message pool
type messagePool struct {
	pool             *refillPool
	messageProtoPool *refillPool
	contentPool      *refillPool
}

// newMessagePool creates a new message pool object
func (p *processPool) newMessagePool() {
	mp := &messagePool{}

	mp.pool = newMessageSyncPool(mp)
	mp.contentPool = newMessageContentSyncPool()
	mp.messageProtoPool = newMessageProtoPool()

	p.messagePool = mp
}

// NewMessageSyncPool gives us a sync pool for messages to prevent allocs
// on the fly
func newMessageSyncPool(p *messagePool) *refillPool {
	return newRefillPool(messagePoolLimit, func() interface{} {
		return &Message{
			messagePool: p,
			done:        true,
		}
	})
}

// NewMessageProtoPool gives us a pool for messages props to prevent allocs
// on the fly
func newMessageProtoPool() *refillPool {
	return newRefillPool(messagePoolLimit, func() interface{} {
		return &MessageProto{}
	})
}

// NewMessageContentSyncPool gives us a sync pool for messageContent to prevent
// allocs on the fly.
// However since maps cnanot be efficiently cleared, the messageContentPool
// never has the same object put in it twice. Instead objects are pre alloc'ed
// when they are requested for the next request, but in a background thread.
// This allows us to handle without an alloc, but means that each content item
// is clean when it starts.
func newMessageContentSyncPool() *refillPool {
	return newRefillPool(messageContentPoolLimit, func() interface{} {
		return messageContent{}
	})
}

// getMessageContentItem gets a new message content object and replaces the object
// in the messageContentPool
func (m *messagePool) getMessageContentItem() messageContent {
	return m.contentPool.get().(messageContent)
}

// NewMessage gets a message struct from the messagePool
func (m *messagePool) NewMessage() *Message {
	message := m.pool.get().(*Message)
	message.messagePool = m
	return message
}

// NewMessage gets a message struct from the messagePool
// and sets it with a content object
func (c *Context) NewMessage() *Message {
	message := c.processPool.messagePool.NewMessage()
	message.arguments = c.processPool.messagePool.getMessageContentItem()
	message.Command = cmdRSP
	message.referenceID = c.input.referenceID

	return message
}

// SetMore tells the rest fo the stack that there are more messages coming
// This is used for multi message worker responses
func (m *Message) SetMore(more bool) {
	m.done = !more
}

// fromProto converts a messageproto type to message
func (m *Message) fromProto(mp *MessageProto) {
	m.Command = mp.GetCommand()
	m.done = mp.GetDone()
	m.referenceID = mp.GetReferenceID()
	m.arguments = mp.GetArguments()
}

func (m *Message) toProto(mp *MessageProto) {
	mp.Command = proto.String(m.Command)
	mp.Done = proto.Bool(m.done)
	mp.ReferenceID = proto.String(m.referenceID)
	mp.Arguments = m.arguments
}
