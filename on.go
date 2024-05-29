package socket

import "github.com/gmelum/socket/entity/context"

func (soc *socket) On(path string, callback func(ctx *context.Context) (err error)) {
	soc.middlewares[path] = append(soc.middlewares[path], callback)
}
