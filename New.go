package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func New(opt Options) *socket {

	return &socket{

		outerWorkers: opt.OuterWorkers,
		outerQueue:   opt.OuterQueue,

		innerWorkers: opt.InnerWorkers,
		innerQueue:   opt.InnerQueue,

		// workers:        opt.Workers,
		// queue:          opt.Queue,
		// spawn:          opt.Spawn,
		// maxConn:        opt.MaxConn,
		// handshakeDelay: opt.HandshakeDelay,

		readBufferSize:  opt.ReadBufferSize,
		writeBufferSize: opt.WriteBufferSize,
		protocol:        opt.Protocol,
		protocolCustom:  opt.ProtocolCustom,
		extension:       opt.Extension,
		extensionCustom: opt.ExtensionCustom,
		negotiate:       opt.Negotiate,
		header:          opt.Header,
		onRequest:       opt.OnRequest,
		onHost:          opt.OnHost,
		onHeader:        opt.OnHeader,
		onBeforeUpgrade: opt.OnBeforeUpgrade,

		eventsConnect:    []func(ch *connect.Connect){},
		eventsDisconnect: []func(ch *connect.Connect){},
		eventsAuth:       []func(url []byte) (err error){},
		eventsCors:       []func(origin []byte) (err error){},
	}
}
