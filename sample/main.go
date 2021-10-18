package main

import (
	"fmt"

	"github.com/techerfan/panda"
)

func main() {

	app := panda.NewApp(panda.Config{
		NotShowLogs:       false,
		CommunicationType: panda.JSON,
	})

	app.NewConnection(func(client *panda.Client) {

		client.OnMessage(func(msg string) {
			fmt.Println("onmessage:", msg)
		})

		client.Subscribe("chat_message", func(msg string) {
			fmt.Println("chat_message:", msg)
		})

		client.Subscribe("goodbye", func(msg string) {

		})
	})

	app.Serve()
}
