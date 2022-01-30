package panda

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/techerfan/panda/logger"
)

func TestNeWClientIds(t *testing.T) {
	client1 := newClient(&websocket.Conn{}, logger.New())
	client2 := newClient(&websocket.Conn{}, logger.New())

	if client1.id == client2.id {
		t.Error("Ids are same")
	}
}
