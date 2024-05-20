package socket

import (
	"testing"

	"github.com/gmelum/socket/entity/connect"
)

func TestNew(t *testing.T) {

	socket := New(Options{})

	socket.Cors(func(origin []byte) (err error) {
		return nil
	})

	socket.Connect(func(ch *connect.Connect) {
		println(ch)
	})

	socket.Disconnect(func(ch *connect.Connect) {
		println(ch)
	})

	socket.Reject(func(ch *connect.Connect) {
		println(ch)
	})

	err := socket.Listen(18300)
	if err != nil {
		panic(err)
	}

}
