package panda

import "github.com/gorilla/websocket"

type channel struct {
	name        string
	clients     []*Client
	msgSender   chan string
	subscribers []*subscriber
}

func NewChannel(name string) *channel {
	channel := &channel{
		name:      name,
		msgSender: make(chan string),
	}
	return channel
}

func (ch *channel) onNewMessage(message string) {
	ch.sendMessageToClients(message)
	ch.sendMessageToSubscribers(message)
}

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

func (ch *channel) subscribe(newSubscriber *subscriber) {
	ch.subscribers = append(ch.subscribers, newSubscriber)
}

func (ch *channel) unsubscribe(toBeRemovedSub *subscriber) {
	for i, subscriber := range ch.subscribers {
		if subscriber == toBeRemovedSub {
			ch.subscribers = append(ch.subscribers[:i], ch.subscribers[i+1:]...)
		}
	}
}

func (ch *channel) sendMessageToClients(message string) {
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

func (ch *channel) sendMessageToSubscribers(message string) {
	for _, subscriber := range ch.subscribers {
		go func() {
			subscriber.lock.Lock()
			defer subscriber.lock.Unlock()
			if subscriber.isOpen {
				subscriber.newMessage <- message
			}
		}()
	}
}

func (ch *channel) destroy() {
	ch = nil
}
