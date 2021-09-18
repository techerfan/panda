package panda

import (
	"github.com/gorilla/websocket"
)

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

	go channel.listener()

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

// sends message to clients which subscribed on the 'pande-client' side.
func (ch *channel) sendMessageToClients(message string) {
	for _, cl := range ch.clients {
		cl.lock.Lock()
		defer cl.lock.Unlock()
		msg := &messageStruct{
			Message: message,
			Channel: ch.name,
			MsgType: Raw,
		}
		cl.conn.WriteMessage(websocket.TextMessage, msg.marshal())
	}
}

// send messages to subscribers.
// subscribers are the listeners that are defined in
// the server side.
func (ch *channel) sendMessageToSubscribers(message string) {
	for _, sub := range ch.subscribers {
		go func(sub *subscriber) {
			sub.lock.Lock()
			defer sub.lock.Unlock()
			if sub.isOpen {
				sub.newMessage <- message
			}
		}(sub)
	}
}

func (ch *channel) destroy() {
	ch = nil
}

// it listens on 'msgSender' channel which is used in order to
// handle channel's new messages.
func (ch *channel) listener() {
	for msg := range ch.msgSender {
		ch.onNewMessage(msg)
	}
}
