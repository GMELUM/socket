package user

import (
	"encoding/json"

	"github.com/gmelum/socket/entity/connect"
)

type User struct {
	events []func(event string, data json.RawMessage) (err error)
	conn   *connect.Connect
}

func (user *User) Init(conn *connect.Connect) (err error) {
	user.conn = conn
	return nil
}

func (user *User) HandlerEvents(event string, data json.RawMessage) (err error) {
	for _, callback := range user.events {
		err := callback(event, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) Events(callback func(event string, data json.RawMessage) (err error)) {
	user.events = append(user.events, callback)
}

func (user *User) Send(event string, data interface{}) (err error) {
	return user.conn.Send(0, event, data)
}
