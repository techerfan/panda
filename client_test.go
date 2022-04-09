package panda

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

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
