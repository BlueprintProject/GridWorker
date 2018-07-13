package gridworker

// messageContent is the base object inside of a message
type messageContent map[string]interface{}

// Message is the object that communications are sent through
type Message struct {
	Command string `json:"cmd,omitempty"`
	Done    bool   `json:"fin,omitempty"`

	Arguments messageContent `json:"arg,omitempty"`

	ReferenceID string `json:"ref,omitempty"`

	messagePool *messagePool
}

// isCtrl determines if a message is a command message or a standard control message
func (m *Message) isCtrl() bool {
	return m.Command == cmdAUTH || m.Command == cmdRSP
}

// messagePool is the general data representation for the message pool
type messagePool struct {
	pool        *refillPool
	contentPool *refillPool
}

// newMessagePool creates a new message pool object
func (p *processPool) newMessagePool() {
	mp := &messagePool{}

	mp.pool = newMessageSyncPool(mp)
	mp.contentPool = newMessageContentSyncPool()

	p.messagePool = mp
}

// NewMessageSyncPool gives us a sync pool for messages to prevent allocs
// on the fly
func newMessageSyncPool(p *messagePool) *refillPool {
	return newRefillPool(messagePoolLimit, func() interface{} {
		return &Message{
			messagePool: p,
			Done:        true,
		}
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
	message.Arguments = c.processPool.messagePool.getMessageContentItem()
	message.Command = cmdRSP
	message.ReferenceID = c.input.ReferenceID

	return message
}

// SetMore tells the rest fo the stack that there are more messages coming
// This is used for multi message worker responses
func (m *Message) SetMore(more bool) {
	m.Done = !more
}
