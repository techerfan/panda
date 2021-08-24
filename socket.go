package panda

import "github.com/gorilla/websocket"

type Socket struct {
	conn *websocket.Conn
}
