package panda

import (
	"net/http"
	"time"

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
	// this handler validates client connection. if it was nil,
	// package considers that authentication is not needed and
	// let the client to establish the connection. it takes a
	// token as input returns a boolean in order to whether continue
	// or not and a time that specifies when connection should be destroyed.
	AuthenticationHandler        func(string) (*time.Time, bool)
	TicketTokenExpirationHandler func(client *Client)
	Logger                       logger.Logger
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
		app.config.Logger = logger.New()
	}

	return app
}

func (a *App) serveWs(rw http.ResponseWriter, r *http.Request, destructionTime *time.Time, ticket string) {
	conn, err := Upgrader.Upgrade(rw, r, nil)
	if err != nil {
		a.config.Logger.Error(err.Error())
		return
	}

	newCl := newClient(a.config.Logger, conn, ticket)

	// to close client's connection after the specified time
	// it is optionanl to set destruction time so that developer
	// can use the package without authentication/authorization.
	if destructionTime != nil {
		timer := time.NewTimer(time.Until(*destructionTime))
		go func() {
			<-timer.C
			a.removeClient(newCl)
			if a.config.TicketTokenExpirationHandler != nil {
				a.config.TicketTokenExpirationHandler(newCl)
			}
			err := newCl.Destroy()
			if err != nil {
				a.config.Logger.Error(err.Error())
			}
		}()
	}

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
		var destructionTime *time.Time
		var ticket string
		if a.config.AuthenticationHandler != nil {
			queries := r.URL.Query()
			ticket = queries.Get("ticket")
			if ticket == "" {
				return
			}
			var isTicketOk bool
			destructionTime, isTicketOk = a.config.AuthenticationHandler(ticket)
			if !isTicketOk {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		a.serveWs(rw, r, destructionTime, ticket)
	})
	a.config.Logger.Info("WebSocket Server is up on: " + a.config.ServerAddress)
	if err := http.ListenAndServe(a.config.ServerAddress, nil); err != nil {
		a.config.Logger.Error(err.Error())
	}
}

func (a *App) Broadcast(channelName string, message string) {
	getChannelsInstance(a.config.Logger).getChannelByName(channelName).sendMessageToClients(message)
}

func (a *App) Destroy(channelName string) {
	getChannelsInstance(a.config.Logger).getChannelByName(channelName).destroy()
}

func (a *App) Send(message string) {
	for _, cl := range a.clients {
		go func(c *Client) {
			c.lock.Lock()
			defer c.lock.Unlock()
			msg, err := newMessage("", message, Raw).marshal()
			if err != nil {
				a.config.Logger.Error(err.Error())
			}
			err = c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				a.config.Logger.Error(err.Error())
				a.removeClient(c)
			}
		}(cl)
	}
}

func (a *App) NewConnection(callback func(client *Client)) {
	// it is not possible to have multiple listeners. so that we must stop
	// other listeners (if any exists) and then make a new one.
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
				go func() {
					callback(newConn)
				}()
			case <-app.stopListening:
				return
			}
		}
	}(a)
}
