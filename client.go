package panda

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/techerfan/panda/logger"
)

var idCounter = uint32(makeRandomInt(3))

type Client struct {
	ctx                context.Context
	app                *App
	conn               *websocket.Conn
	lock               *sync.Mutex
	id                 string
	stopListening      chan bool
	isListening        bool
	newMessage         chan string
	subscribedChannels []*channel
	listeners          map[string]chan string
	ticket             string
	logger             logger.Logger
}

func newClient(
	ctx context.Context,
	app *App,
	logger logger.Logger,
	conn *websocket.Conn,
	ticket string,
) *Client {
	client := &Client{
		ctx:           ctx,
		app:           app,
		conn:          conn,
		lock:          &sync.Mutex{},
		id:            makeId(),
		stopListening: make(chan bool),
		newMessage:    make(chan string),
		listeners:     make(map[string]chan string),
		ticket:        ticket,
		logger:        logger,
	}

	go client.reader()

	closeHandlerInstance := conn.CloseHandler()
	conn.SetCloseHandler(func(code int, text string) error {
		close(client.stopListening)
		client.closeHandler()
		client = nil
		return closeHandlerInstance(code, text)
	})

	return client
}

func (c *Client) OnMessage(callback func(msg string)) {
	c.isListening = true
	go func() {
		for {
			select {
			case msg := <-c.newMessage:
				callback(msg)
			case <-c.stopListening:
				c.isListening = false
				return
			}
		}
	}()
}

func (c *Client) On(channelName string, callback func(msg string)) {
	listenerChan := make(chan string)

	c.listeners[channelName] = listenerChan

	for message := range listenerChan {
		callback(message)
	}
}

func (c *Client) Send(message string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	msg, err := newMessage("", message, Raw).marshal()
	if err != nil {
		c.logger.Error(err.Error())
		return
	}
	err = c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			// Because Destroy uses the same lock as this method
			// does, therefore it should be called as a separate
			// goroutine so this method can end and unlock the lock.
			// Otherwise we will have a livelock.
			go c.Destroy()
		}
		c.logger.Error(err.Error())
	}
}

func (c *Client) Publish(channel string, message string) {
	go func() {
		ch := getChannelsInstance(c.logger).getChannelByName(channel)
		ch.msgSender <- message
	}()
}

func (c *Client) GetTicket() string {
	return c.ticket
}

func (c *Client) Destroy() error {
	defer func() {
		if recover() != nil {
			c.logger.Error("an error occured while destroying a client")
		}
	}()
	c.lock.Lock()
	defer c.lock.Unlock()
	close(c.stopListening)
	// because 'closeHandler' method sets client to nil, we
	// should close the connection before we lose it.
	err := c.conn.Close()
	c.closeHandler()
	return err
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) Context() context.Context {
	return c.ctx
}

func (c *Client) reader() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		messageStruct, err := unmarshalMsg(msg)
		if err != nil {
			c.logger.Error(err.Error())
		}

		if messageStruct != nil {
			switch messageStruct.MsgType {
			case Subscribe:
				c.subscribeToChannel(messageStruct.Channel)
			case Unsubscribe:
				c.unsubscribeToChannel(messageStruct.Channel)
			case Raw:
				c.receiveRawMsg(messageStruct)
			}
		}
	}
}

func (c *Client) subscribeToChannel(channelName string) {
	ch := getChannelsInstance(c.logger).getChannelByName(channelName)
	ch.addClient(c)
	c.subscribedChannels = append(c.subscribedChannels, ch)
}

func (c *Client) unsubscribeToChannel(channelName string) {
	ch := getChannelsInstance(c.logger).getChannelByName(channelName)
	ch.removeClient(c)
	for i, channel := range c.subscribedChannels {
		if ch == channel {
			c.subscribedChannels = append(c.subscribedChannels[:i], c.subscribedChannels[i+1:]...)
		}
	}
}

func (c *Client) receiveRawMsg(msg *messageStruct) {
	if msg.Channel != "" {
		if ch, ok := c.listeners[msg.Channel]; ok {
			ch <- msg.Message
		}
	} else {
		if c.isListening {
			c.newMessage <- msg.Message
		}
	}
}

func (c *Client) closeHandler() {
	for _, ch := range c.subscribedChannels {
		ch.removeClient(c)
	}
	c.app.removeClient(c)
	c = nil
}

// generates a random unsinged integer (32 bit).
func makeRandomInt(bytesLen int) uint64 {
	b := make([]byte, bytesLen)
	if _, err := rand.Reader.Read(b); err != nil {
		panic(fmt.Errorf("[Client]: cannot generate random number: %v", err))
	}

	// making the number regarding bytes length:
	var randomInt uint64 = 0
	for i := 1; i < bytesLen; i++ {
		randomInt = randomInt | uint64(b[i-1])<<(8*(bytesLen-i))
	}
	return randomInt
}

// makes id using mongoDB standard.
// follow this link: https://docs.mongodb.com/manual/reference/method/ObjectId/
func makeId() string {
	var id [12]byte
	t := time.Now().Unix()
	// 4 bytes timestamp value
	binary.BigEndian.PutUint32(id[:], uint32(t))

	// 5 bytes random value
	randomNum := makeRandomInt(5)
	id[4] = byte(randomNum >> 32)
	id[5] = byte(randomNum >> 24)
	id[6] = byte(randomNum >> 16)
	id[7] = byte(randomNum >> 8)
	id[8] = byte(randomNum)

	// 3 bytes incrementing counter, initialized to a random value
	i := atomic.AddUint32(&idCounter, 1)
	id[9] = byte(i >> 16)
	id[10] = byte(i >> 8)
	id[11] = byte(i)

	return fmt.Sprintf("%x", id)
}
