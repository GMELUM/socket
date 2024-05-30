package socket

import (
	"net"

	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/entity/context"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

type socket struct {
	mode string

	listener net.Listener

	outerWorkers int
	outerTasks   int

	innerWorkers int
	innerTasks   int

	readBufferSize  int
	writeBufferSize int
	protocol        func([]byte) bool
	protocolCustom  func([]byte) (string, bool)
	extension       func(httphead.Option) bool
	extensionCustom func([]byte, []httphead.Option) ([]httphead.Option, bool)
	negotiate       func(httphead.Option) (httphead.Option, error)
	header          ws.HandshakeHeader

	eventsConnect    []func(cn *connect.Connect) (err error)
	eventsDisconnect []func(cn *connect.Connect)
	eventsReject     []func(cn *connect.Connect)
	eventsCors       []func(origin string) (err error)

	middlewares map[string][]func(ctx *context.Context) (err error)
	events      map[string]func(ctx *context.Context) (err error)
}
