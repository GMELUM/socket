package socket

import (
	"fmt"

	"github.com/gmelum/socket/entity/context"
)

func (soc *socket) Event(path string, callback func(ctx *context.Context) (err error)) {
	if _, ok := soc.events[path]; ok {
		panic(fmt.Sprintf("event \"%v\" already exists", path))
	}
	soc.events[path] = callback
}
