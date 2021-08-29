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
	config  Config
	clients []*Client
	newConn chan Client
	// to check if app listens on new connection
	isListening bool
	// to stop apps from listening on new connections
	stopListening chan bool
}

type Config struct {
	ServerAddress     string
	WebSocketPath     string
	CommunicationType CommunicationType
}

func NewApp(config Config) *App {
	app := &App{
		config:        config,
		newConn:       make(chan Client),
		stopListening: make(chan bool),
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

	newCl := newClient(conn)

	// whenever a new client joins, we will send it over newConn channel
	// but app must listens on new connections.
	// we did this because if nobody listens on channel, Go will exit
	// the program by code 1.
	if a.isListening {
		a.newConn <- *newCl
	}
}

func (a *App) Serve() {
	http.HandleFunc(a.config.WebSocketPath, func(rw http.ResponseWriter, r *http.Request) {
		a.serveWs(rw, r)
	})
	http.ListenAndServe(a.config.ServerAddress, nil)
}

func (a *App) Send(channelName string, message string) {
	getChannelsInstance().getChannelByName(channelName).sendMessage(message)
}

func (a *App) NewConnection(callback func(client *Client)) {
	// it is not possible to have multiple listeners. so that we stop
	// other listeners (if any exist) and then make a new one.
	if a.isListening {
		a.stopListening <- true
		a.isListening = true
	}
	go func(app *App) {
		app.isListening = true
		for {
			select {
			case newConn := <-app.newConn:
				app.clients = append(app.clients, &newConn)
				callback(&newConn)
			case <-app.stopListening:
				return
			}
		}
	}(a)
}
