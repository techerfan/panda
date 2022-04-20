# Panda
Panda is a library for event-based communications via **WebSocket**.

## Client Package:
[This](https://github.com/techerfan/panda-client) is `panda-client` written in TypeScript. You can use it in your modern Web Apps along with `Panda`.

### installation
```
go get -u github.com/techerfan/panda
```

## How to use?
First you need to create a new `App`. 
In order to create a new App, you have this option to whether pass configuration or not. Configuration consists of:
1. **`ServerAddress`**: It is the address to the server. The default is `:8000`.
2. **`WebSocketPath`**: The path of Web Socket. The default is `/ws`.
3. **`CommunicationType`**: You can choose the method of sending your data via Web Socket. The default is `JSON` and unfortunately, `XML` and `Binary` are not implemented yet.
4. **`DoNotShowLogs`**: It is a boolean. If it is `true`, the module will not print logs and if it is `false`, The logger will work and you will be able to see logs. The default is `false`.
5. **`LogsHeader`**: It is a `string` item. The logger will add it to the beginning of each log.
The default is `Panda`.
6. **`AuthenticationHandler`**: This handler validates client's connection. If it is nil, package will consider that authentication is not needed and let the client to establish the connection. It takes a token as input and returns a boolean in order to specify whether continue or not and a time that shows when the connection should be destroyed.
7. **`TicketTokenExpirationHandler`**: This handler decides what to do when a client's ticket is expired. If it is nil, there will be no default behavior.
8. **`Logger`**: You can use your own logger if it follows [this](logger/logger.go) interface.

⚠️ If you want to authenticate your clients by `AuthenticationHandler` you need to generate a **ticket** by yourself and add it to WebSocket URL as a query (e.g. http://localhost:8000/ws?ticket=MY_TICKET). This way you will receive the client's ticket by `AuthenticationHandler` when it tries to establish connection in order to validate the ticket.

### Initalizing 
```golang
import "github.com/techerfan/panda"

app := panda.NewApp(panda.Config{
  ServerAddress: ":8080",
  WebSocketPath: "/ws",
  CommunicationType: panda.JSON,
  NotShowLogs: flase,
  LogsHeader: "My App",
})

// or you can stick to the defaults and simply do this:

app := panda.NewApp()
```

## New Connection

Panda lets you know, by the `NewConnection` method, whenever a client connects to the server:

```golang
app.NewConnection(func(client *panda.Client) {
  // handle the client here...
})
```

### Client Methods
1. `OnMessage`: Listens to the messages which client sends to the server:
```golang
client.OnMessage(func (msg string) {
  // do whatever you want with the message...
})
```
2. `On`: Listens to the messages that are exchanged over a specified channel. You can decide whether send them to the clients or do something else:
```golang
client.On("chat_message", func(msg string) {
	// do sth with the message...
})
```
4. `Send`: To Send a message to the client: 
```golang
client.Send("your message")
``` 
5. `Publish`: To publish a message over a specified channel:
```golang
client.Publish("channel_name", "your message")
```
7. `GetTicket`: Each client is authenticated via a ticket. You can get this ticket by this method:
```golang
var ticket string
ticket = client.GetTicket()
```
6. `Destroy`: To destroy a client's connection.
```golang
err := client.Destroy()
```

## License 
Licensed under the [MIT License](/LICENSE).
