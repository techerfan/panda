package panda

import "github.com/gorilla/websocket"

type channel struct {
	name    string
	clients []*Client
}

func NewChannel(name string) *channel {
	channel := &channel{
		name: name,
	}
	return channel
}

// func (c *channel) broadcast() {

// }

func (c *channel) addClient(cl *Client) {
	c.clients = append(c.clients, cl)
}

func (c *channel) removeClient(cl *Client) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	for i, el := range c.clients {
		if el == cl {
			c.clients = append(c.clients[:i], c.clients[i+1:]...)
		}
	}
}

func (c *channel) sendMessage(msg *messageStruct) {
	for _, cl := range c.clients {
		cl.lock.Lock()
		defer cl.lock.Unlock()
		cl.conn.WriteMessage(websocket.BinaryMessage, msg.marshal())
	}
}

func (c *channel) destroy() {
	c = nil
}
