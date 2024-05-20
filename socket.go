package socket

import (
	"net"

	"github.com/gmelum/socket/entity/connect"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

type socket struct {
	listener net.Listener

	outerWorkers int
	outerTasks   int

	innerWorkers int
	innerTasks   int

	eventsConnect    []func(cn *connect.Connect) (err error)
	eventsDisconnect []func(cn *connect.Connect)
	eventsReject     []func(cn *connect.Connect)
	eventsCors       []func(origin string) (err error)
	eventsRequest    []func(uri string) (err error)

	readBufferSize  int
	writeBufferSize int
	protocol        func([]byte) bool
	protocolCustom  func([]byte) (string, bool)
	extension       func(httphead.Option) bool
	extensionCustom func([]byte, []httphead.Option) ([]httphead.Option, bool)
	negotiate       func(httphead.Option) (httphead.Option, error)
	header          ws.HandshakeHeader
	onRequest       func(uri []byte) error
	onHost          func(host []byte) error
	onHeader        func(key, value []byte) error
	onBeforeUpgrade func() (header ws.HandshakeHeader, err error)
}
