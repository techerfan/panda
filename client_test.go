package panda

import (
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
