package socket

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gmelum/socket/entity/connect"

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
				soc.onRequest,
				soc.onHost,
				soc.onHeader,
				soc.onBeforeUpgrade,

				&soc.eventsCors,
				&soc.eventsRequest,
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

			time.Sleep(time.Second * 20)

			conn.Event(func(ev netpoll.Event) {

				if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
					once.Do(handlerClose)
					return
				}

				if conn.Status != "OPEN" {
					once.Do(handlerClose)
				}

				innerPool.Submit(func() {
					println("message")
					time.Sleep(time.Second * 5)
				})

			})
		})

		// outerPool.Submit(func() {
		// 	for _, callback := range soc.eventsConnect {
		// 		callback(conn)
		// 	}
		// })

	})

	select {}

}
