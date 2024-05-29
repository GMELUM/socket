package socket

import (
	"fmt"
	"net"
	"sync"

	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/entity/context"

	"github.com/mailru/easygo/netpoll"

	"github.com/alitto/pond"
)

func (soc *socket) Listen(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	poller, err := netpoll.New(nil)
	if err != nil {
		return err
	}

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		listener, netpoll.EventRead|netpoll.EventOneShot,
	))

	outerPool := pond.New(soc.outerWorkers, soc.outerTasks)
	innerPool := pond.New(soc.innerWorkers, soc.innerTasks)

	poller.Start(acceptDesc, func(e netpoll.Event) {
		outerPool.Submit(func() {

			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Возникла паника:", err)
				}
			}()

			once := new(sync.Once)
			poller.Resume(acceptDesc)

			conn, err := connect.New(

				listener,
				poller,
				innerPool,

				soc.readBufferSize,
				soc.writeBufferSize,
				soc.protocol,
				soc.protocolCustom,
				soc.extension,
				soc.extensionCustom,
				soc.negotiate,
				soc.header,

				&soc.eventsCors,
				&soc.eventsConnect,
			)
			if err != nil {
				for _, callback := range soc.eventsReject {
					callback(conn)
				}
				return
			}

			handlerClose := func() {
				conn.Close()
				for _, callback := range soc.eventsDisconnect {
					callback(conn)
				}
			}

			conn.Event(func(ev netpoll.Event) {

				if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
					once.Do(handlerClose)
					return
				}

				innerPool.Submit(func() {
					msg, err := conn.Read()

					if err == nil {

						if msg.ID == 0 && msg.Type == "ping" {
							conn.Send(0, "pong", "")
							return
						}

						ctx := &context.Context{
							ID:   msg.ID,
							Type: msg.Type,
							Data: msg.Value,
							Conn: conn,
						}

						if list, ok := soc.middlewares[msg.Type]; ok {
							if len(list) > 0 {
								for _, callback := range list {
									err := callback(ctx)
									if err != nil {
										ctx.Conn.Send(ctx.ID, ctx.Type, err.Error())
										return
									}
								}
							}
						}

						if callback, ok := soc.events[msg.Type]; ok {
							err := callback(ctx)
							if err != nil {
								ctx.Conn.Send(ctx.ID, ctx.Type, err.Error())
								return
							}
						}

					}
				})

			})

		})

	})

	select {}

}
