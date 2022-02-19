package panda

import (
	"github.com/gorilla/websocket"
	"github.com/techerfan/panda/logger"
)

type channel struct {
	name      string
	clients   []*Client
	msgSender chan string
	logger    logger.Logger
}

func NewChannel(logger logger.Logger, name string) *channel {
	channel := &channel{
		name:      name,
		msgSender: make(chan string),
	}

	go channel.listener()

	return channel
}

func (ch *channel) onNewMessage(message string) {
	ch.sendMessageToClients(message)
	// ch.sendMessageToSubscribers(message)
}

func (ch *channel) addClient(cl *Client) {
	ch.clients = append(ch.clients, cl)
}

func (ch *channel) removeClient(cl *Client) {
	for i, el := range ch.clients {
		if el == cl {
			ch.clients = append(ch.clients[:i], ch.clients[i+1:]...)
			return
		}
	}
}

// sends message to clients which subscribed on the 'pande-client' side.
func (ch *channel) sendMessageToClients(message string) {
	msg, err := (&messageStruct{
		Message: message,
		Channel: ch.name,
		MsgType: Raw,
	}).marshal()
	if err != nil {
		ch.logger.Error(err.Error())
	}
	for _, cl := range ch.clients {
		go func(cl *Client) {
			cl.lock.Lock()
			defer cl.lock.Unlock()
			err := cl.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				cl.logger.Error(err.Error())
			}
		}(cl)
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
