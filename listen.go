package socket

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/entity/context"

	"github.com/mailru/easygo/netpoll"

	"github.com/alitto/pond"
)

func (soc *socket) Listen(port int) error {

	if soc.mode == "poller" {
		return soc.newPoller(port)
	}

	if soc.mode == "default" {
		return soc.newDefault(port)
	}

	return errors.New("the specified server mode is not supported")

}

func (soc *socket) newDefault(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	// outerPool := pond.New(soc.outerWorkers, soc.outerTasks)
	innerPool := pond.New(soc.innerWorkers, soc.innerTasks)

	for {

		once := new(sync.Once)

		conn, err := connect.New(

			listener,
			nil,
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
			continue
		}

		handlerClose := func() {
			conn.Close()
			for _, callback := range soc.eventsDisconnect {
				callback(conn)
			}
		}

		go func() {
			defer once.Do(handlerClose)

			for {
				msg, err := conn.Read()

				if err == nil {

					if msg.ID == 0 && msg.Type == "ping" {
						conn.Send(0, "pong", "")
						continue
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
									continue
								}
							}
						}
					}

					if callback, ok := soc.events[msg.Type]; ok {
						err := callback(ctx)
						if err != nil {
							ctx.Conn.Send(ctx.ID, ctx.Type, err.Error())
							continue
						}
					}

				}
			}

		}()

	}

}

func (soc *socket) newPoller(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}

	acceptDesc := netpoll.Must(netpoll.HandleListener(
		listener, netpoll.EventRead|netpoll.EventOneShot,
	))

	outerPool := pond.New(soc.outerWorkers, soc.outerTasks)
	innerPool := pond.New(soc.innerWorkers, soc.innerTasks)

	poller, err := netpoll.New(nil)
	if err != nil {
		return err
	}

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
