package gridworker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// StartServer will start up a socket server that will
// allow other distributed workers to connect to it
func (w *DistributedWorker) StartServer(address string) {
	http.HandleFunc("/socket", w.onSocket)
	http.ListenAndServe(address, nil)
}

// ConnectToServer will connect to another distributed worker
// at a specific address
func (w *DistributedWorker) ConnectToServer(address string) {
	n, _, err := websocket.DefaultDialer.Dial(address, nil)

	if err != nil {
		log.Fatal("connect error:", err)
	}

	ctx := w.newConnection(n)
	go ctx.listen()
	go ctx.listenForWrite()
	go ctx.announceAuth()
}

// onSocket is called while a new socket connection is made
func (w *DistributedWorker) onSocket(wr http.ResponseWriter, r *http.Request) {
	c, err := w.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	ctx := w.newConnection(c)
	go ctx.listenForWrite()
	ctx.listen()
}

// connection is a representation of a distributedWorker connection
type connection struct {
	// conn is the web socket connection
	conn *websocket.Conn

	// distributedWorker is the DistributedWorker object the connection is
	// related to
	distributedWorker *DistributedWorker

	// writeChan is where new messages are sent before being sent over the socket
	writeChan chan *Message

	// announcedAuth determines if the connection has sent an auth message
	annoucedAuth bool
}

// newConnection is called when a new socket connection is made,
// it generates a new connection object
func (w *DistributedWorker) newConnection(conn *websocket.Conn) *connection {
	c := &connection{
		conn:              conn,
		distributedWorker: w,
		writeChan:         make(chan *Message, 25),
	}

	return c
}

// listen watches for new messages coming from the socket
// and then processes the messages in a loop
func (c *connection) listen() {
	for {
		err := c.listenForMesssage()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection Closed")
				c.conn.Close()
				break
			}

			err := c.handleError()
			if err != nil {
				fmt.Println("Error when attempting to handle error")
			}

			c.conn.Close()
			break
		}
	}
}

// listenForWrite will watch the writeChann for outgoing messages
// and send them synchronously
func (c *connection) listenForWrite() {
	for {
		m := <-c.writeChan
		b, _ := json.Marshal(m)
		c.conn.WriteMessage(websocket.BinaryMessage, b)
	}
}

// announceAuth sends out an auth message to the socket
func (c *connection) announceAuth() {
	if c.annoucedAuth {
		return
	}

	m := c.distributedWorker.NewMessage()
	m.Command = cmdAUTH
	m.SetString("id", c.distributedWorker.id)
	c.annoucedAuth = true
	c.writeChan <- m
}

// handleError is called when an error occurs in processing
// a socket response
func (c *connection) handleError() (err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("Unknown Error")
		}
	}()

	return err
}

// listenForMessage waits for messages to come in over the socket
// it does not handle errors and will instead hand those off to
// the listen function
func (c *connection) listenForMesssage() (err error) {
	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}

		go c.processPacket(b)
	}
}

// processPacket gets the byte data from the socket and converts it
// to a message
func (c *connection) processPacket(b []byte) {
	m := c.distributedWorker.NewMessage()
	json.Unmarshal(b, m)

	c.distributedWorker.handleMessageForConnection(m, c)
}
