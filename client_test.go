package panda

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestNeWClientIds(t *testing.T) {
	client1 := newClient(&websocket.Conn{})
	client2 := newClient(&websocket.Conn{})

	if client1.id == client2.id {
		t.Error("Ids are same")
	}
}
