package socket

import (
	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/entity/context"
)

func New(opt Options) *socket {

	if opt.OuterWorkers == 0 {
		opt.OuterWorkers = 1
	}

	if opt.OuterTasks == 0 {
		opt.OuterTasks = 1
	}

	if opt.InnerWorkers == 0 {
		opt.InnerWorkers = 1
	}

	if opt.InnerTasks == 0 {
		opt.InnerTasks = 1
	}

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

		eventsConnect:    []func(ch *connect.Connect) (err error){},
		eventsDisconnect: []func(ch *connect.Connect){},
		eventsCors:       []func(origin string) (err error){},

		middlewares: make(map[string][]func(ctx *context.Context) (err error)),
		events:      make(map[string]func(ctx *context.Context) (err error)),
	}
}
