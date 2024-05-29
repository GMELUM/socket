package main

import (
	"github.com/gmelum/socket"
	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/entity/context"
	// "github.com/gmelum/socket/errors"
)

func main() {

	// Create a new websocket server instance
	server := socket.New(socket.Options{

		// Limiting the processing of new connections
		OuterWorkers: 10,
		OuterTasks:   100,

		// We limit the processing of messages from
		// connected users
		InnerWorkers: 1000,
		InnerTasks:   10000,
	})

	// Receiving and checking origin from the websocket
	// connection request headers. Called before upgrading to websocket
	server.Cors(func(origin string) (err error) {
		println("origin:", string(origin))
		// if string(origin) != "example.com" {
		// 	return errors.Forbidden()
		// }
		return nil
	})

	// We receive a connection when a user connects
	server.Connect(func(cn *connect.Connect) (err error) {
		println("connect:", cn.ID)
		return nil
	})

	// Receiving a connection when a user disconnects
	server.Disconnect(func(cn *connect.Connect) {
		println("disconnect:", cn.ID)
	})

	// Getting a connection when the user resets
	server.Reject(func(cn *connect.Connect) {
		if cn != nil {
			println("reject:", cn.ID)
		}
	})

	// First Handling. Processing middleware executed before the main request processing
	server.On("user.get", func(ctx *context.Context) (err error) {
		return nil
	})

	// Second Handling. Processing middleware executed before the main request processing
	server.On("user.get", func(ctx *context.Context) (err error) {
		return nil
	})

	// Handling events sent from the user
	server.Event("user.get", func(ctx *context.Context) (err error) {

		data := map[string]interface{}{
			"id":         123,
			"first_name": "Artur",
			"last_name":  "Getman",
			"level":      99,
		}

		return ctx.Answer(data)

	})

	// Start listening on the port 18300
	err := server.Listen(18300)
	if err != nil {
		panic(err)
	}

}
