package socket

import (
	"github.com/gmelum/socket/entity/connect"
)

func New(opt Options) *socket {
	return &socket{

		outerWorkers: opt.OuterWorkers,
		outerTasks:   opt.OuterTasks,

		innerWorkers: opt.InnerWorkers,
		innerTasks:   opt.InnerTasks,

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

		eventsConnect:    []func(ch *connect.Connect) (err error){},
		eventsDisconnect: []func(ch *connect.Connect){},
		eventsCors:       []func(origin string) (err error){},
		eventsRequest:    []func(uri string) (err error){},
	}
}
