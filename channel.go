package panda

import "github.com/gorilla/websocket"

type channel struct {
	name      string
	clients   []*Client
	msgSender chan string
}

func NewChannel(name string) *channel {
	channel := &channel{
		name:      name,
		msgSender: make(chan string),
	}
	return channel
}

// func (c *channel) broadcast() {

// }

func (ch *channel) addClient(cl *Client) {
	ch.clients = append(ch.clients, cl)
}

func (ch *channel) removeClient(cl *Client) {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	for i, el := range ch.clients {
		if el == cl {
			ch.clients = append(ch.clients[:i], ch.clients[i+1:]...)
		}
	}
}

func (ch *channel) sendMessage(message string) {
	for _, cl := range ch.clients {
		cl.lock.Lock()
		defer cl.lock.Unlock()
		msg := &messageStruct{
			Message: message,
			Channel: ch.name,
			MsgType: Raw,
		}
		cl.conn.WriteMessage(websocket.BinaryMessage, msg.marshal())
	}
}

func (ch *channel) destroy() {
	ch = nil
}
