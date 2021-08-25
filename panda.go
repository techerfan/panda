package panda

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const DefaultWebSocketPath = "/ws"

var Upgrader = websocket.Upgrader{
	ReadBufferSize:    0,
	WriteBufferSize:   0,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type CommunicationType int

const (
	JSON CommunicationType = iota
	BINARY
	// not implemented yet
	XML
)

type App struct {
	config Config
	client []*client
}

type Config struct {
	ServerAddress     string
	WebSocketPath     string
	CommunicationType CommunicationType
}

func NewApp(config Config) *App {
	app := &App{
		config: config,
	}

	if config.WebSocketPath == "" {
		app.config.WebSocketPath = DefaultWebSocketPath
	}

	return app
}

func (a *App) serveWs(rw http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	a.client = append(a.client, &client{
		conn: conn,
	})
}

func (a *App) Serve() {
	http.HandleFunc(a.config.WebSocketPath, func(rw http.ResponseWriter, r *http.Request) {
		a.serveWs(rw, r)
	})
	http.ListenAndServe(a.config.ServerAddress, nil)
}

func (a *App) Channel() {

}
