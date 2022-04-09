package panda

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/techerfan/panda/logger"
)

func TestNeWClientIds(t *testing.T) {
	client1 := newClient(logger.New(), &websocket.Conn{}, "")
	client2 := newClient(logger.New(), &websocket.Conn{}, "")

	if client1.id == client2.id {
		t.Error("Ids are same")
	}
}

func TestDestroy(t *testing.T) {
	t.Run("close a closed channel", func(t *testing.T) {
		cl := &Client{
			lock:          &sync.Mutex{},
			stopListening: make(chan bool),
			conn:          &websocket.Conn{},
		}
		close(cl.stopListening)
		err := cl.Destroy()
		if err != nil {
			t.Error(err)
		}
	})
}
