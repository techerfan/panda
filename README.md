# Panda
Panda is a library for event-based communications via **WebSocket**.

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
4. **`NotShowLogs`**: It is a boolean. If it is `true`, the module will not print logs and if it is `false`, The logger will work and you will be able to see logs. The default is `false`.
5. **`LogsHeader`**: It is a `string` item. The logger will add it to the beginning of each log.
The default is `Panda`.

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

})
```

### Client Methods
1. `OnMessage`: Listens to the messages that client sends them to server:
```golang
client.OnMessage(func (msg string) {
  // do whatever you want with the message...
})
```
2. `On`: Listens to the messages that are exchanged over a specified channel. You can decide whether send them to clients or do something else:
```golang
client.On("chat_message", func(msg string) {
	// do sth with the message...
})
```
4. `Send`: To Send a message to client: 
```golang
client.Send("your message")
``` 
5. `Publish`: To publish a message over a specified channel:
```golang
client.Publish("channel_name", "your message")
```

## License 
Licensed under the [MIT License](/LICENSE).
