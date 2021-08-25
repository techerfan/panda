package panda

import "github.com/gorilla/websocket"

type client struct {
	conn *websocket.Conn
}
