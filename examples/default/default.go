package main

import (
	"github.com/gmelum/socket"
	"github.com/gmelum/socket/entity/connect"
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
		return nil
	})

	// Get the connection uri. Called before upgrading to websocket
	server.Request(func(uri string) (err error) {
		println("uri:", string(uri))
		return nil
	})

	// We receive a connection when a user connects
	server.Connect(func(ch *connect.Connect) (err error) {
		println("connect:", ch.ID)
		return nil
	})

	// Receiving a connection when a user disconnects
	server.Disconnect(func(ch *connect.Connect) {
		println("disconnect:", ch.ID)
	})

	// Getting a connection when the user resets
	server.Reject(func(ch *connect.Connect) {
		println("reject:", ch.ID)
	})

	// Start listening on the port 18300
	err := server.Listen(18300)
	if err != nil {
		panic(err)
	}

}
