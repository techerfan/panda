package panda

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn          *websocket.Conn
	lock          *sync.Mutex
	id            string
	stopListening chan bool
}

var idCounter = uint32(makeRandomInt(3))

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

func newClient(conn *websocket.Conn) *Client {
	client := &Client{
		conn:          conn,
		id:            makeId(),
		stopListening: make(chan bool),
		newMessage:    make(chan string),
	}

	go client.reader()

	closeHandlerInstance := conn.CloseHandler()

	conn.SetCloseHandler(func(code int, text string) error {
		close(client.stopListening)
		return closeHandlerInstance(code, text)
	})

	return client
}

func (c *Client) reader() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		messageStruct := unmarshalMsg(msg)
		if messageStruct != nil {
			switch messageStruct.MsgType {
			case Subscribe:
				c.subscribe(messageStruct.Channel)
			case Unsubscribe:
				c.unsubscribe(messageStruct.Channel)
			case Raw:
				c.receiveRawMsg(messageStruct)
			}
		}
	}
}

func (c *Client) subscribe(channelName string) {
	getChannelsInstance().getChannelByName(channelName).addClient(c)
}

func (c *Client) unsubscribe(channelName string) {
	getChannelsInstance().getChannelByName(channelName).removeClient(c)
}

func (c *Client) receiveRawMsg(msg *messageStruct) {
	if msg.Channel != "" {
		ch := getChannelsInstance().getChannelByName(msg.Channel)
		ch.msgSender <- msg.Message
	} else {
		c.newMessage <- msg.Message
	}
}

	}
}

func OnMessage(callback func(string)) {

}

func (c *Client) Listen(channelName string, callback func(msg string)) {
	ch := getChannelsInstance().getChannelByName(channelName)

	go func(ch *channel, c *Client) {
		for {
			select {
			case msg := <-ch.msgSender:
				callback(msg)
			case <-c.stopListening:
				return
			}
		}
	}(ch, c)
}
