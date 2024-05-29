package context

import (
	"encoding/json"

	"github.com/gmelum/socket/entity/connect"
)

type Context struct {
	ID   int             `json:"id"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`

	Conn *connect.Connect `json:"-"`
}

func (ctx *Context) Answer(data interface{}) error {
	err := ctx.Conn.Send(ctx.ID, ctx.Type, data)
	if err != nil {
		// ctx.Conn.Close()
		return err
	}
	return nil
}
