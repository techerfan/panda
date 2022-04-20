package main

import (
	"fmt"

	"github.com/techerfan/panda"
)

func main() {

	app := panda.NewApp(panda.Config{
		DoNotShowLogs:     false,
		CommunicationType: panda.JSON,
	})

	app.NewConnection(func(client *panda.Client) {

		client.OnMessage(func(msg string) {
			fmt.Println("onmessage:", msg)
		})
	})

	app.Serve()
}
