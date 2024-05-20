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
	outerQueue   int

	innerWorkers int
	innerQueue   int

	eventsConnect    []func(cn *connect.Connect)
	eventsDisconnect []func(cn *connect.Connect)
	eventReject      []func(cn *connect.Connect)
	eventsAuth       []func(url []byte) (err error)
	eventsCors       []func(origin []byte) (err error)

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
