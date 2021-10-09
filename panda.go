package panda

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/techerfan/panda/logger"
)

const (
	DefaultWebSocketPath = "/ws"
	DefaultLogsHeader    = "Panda"
	DefaultServerAddress = ":8000"
)

type CommunicationType int

const (
	JSON CommunicationType = iota
	// not implemented yet
	BINARY
	XML
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:    0,
	WriteBufferSize:   0,
	EnableCompression: true,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type App struct {
	config  Config
	clients []*Client
	newConn chan *Client
	// to check if app listens on new connection
	isListening bool
	// to stop apps from listening on new connections
	stopListening chan bool
}

type Config struct {
	ServerAddress     string
	WebSocketPath     string
	CommunicationType CommunicationType
	// to choose if module print logs or not
	NotShowLogs bool
	// a name that will be showed in logs between [] like [Panda]
	Logsheader string
	Logger     logger.LoggerInterface
}

func NewApp(config ...Config) *App {
	app := &App{
		config:        Config{},
		newConn:       make(chan *Client),
		stopListening: make(chan bool),
	}

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.WebSocketPath == "" {
		app.config.WebSocketPath = DefaultWebSocketPath
	}

	if app.config.Logsheader == "" {
		app.config.Logsheader = DefaultLogsHeader
	}

	if app.config.ServerAddress == "" {
		app.config.ServerAddress = DefaultServerAddress
	}

	if app.config.Logger == nil {
		app.config.Logger = logger.GetLogger()
	} else {
		logger.SetLogger(app.config.Logger)
	}

	// app.initializeLogger()

	return app
}

// func (a *App) initializeLogger() {
// 	l := logger.GetLogger()
// 	l.SetName(a.config.Logsheader)
// 	l.SetShowLogs(!a.config.NotShowLogs)
// }

func (a *App) serveWs(rw http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(rw, r, nil)
	if err != nil {
		logger.GetLogger().Log(logger.Error, err.Error())
		return
	}

	newCl := newClient(conn)

	// whenever a new client joins, we will send it over newConn channel
	// but app must listens on new connections.
	// we did this because if nobody listens on channel, Go will exit
	// the program by code 1.
	if a.isListening {
		a.newConn <- newCl
	}
}

func (a *App) removeClient(c *Client) {
	for i, cl := range a.clients {
		if cl == c {
			a.clients = append(a.clients[:i], a.clients[i+1:]...)
			break
		}
	}
}

func (a *App) Serve() {
	http.HandleFunc(a.config.WebSocketPath, func(rw http.ResponseWriter, r *http.Request) {
		a.serveWs(rw, r)
	})
	logger.GetLogger().Log(logger.Info, "WebSocket Server is up on: "+a.config.ServerAddress)
	if err := http.ListenAndServe(a.config.ServerAddress, nil); err != nil {
		logger.GetLogger().Log(logger.Error, err.Error())
	}
}

func (a *App) Broadcast(channelName string, message string) {
	getChannelsInstance().getChannelByName(channelName).sendMessageToClients(message)
}

func (a *App) Send(message string) {
	for _, cl := range a.clients {
		go func(c *Client) {
			c.lock.Lock()
			defer c.lock.Unlock()
			err := c.conn.WriteMessage(websocket.TextMessage, newMessage("", message, Raw).marshal())
			if err != nil {
				a.removeClient(c)
			}
		}(cl)
	}
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
				app.clients = append(app.clients, newConn)
				callback(newConn)
			case <-app.stopListening:
				return
			}
		}
	}(a)
}
