package socket

import (
	"bytes"
	"fmt"
	"net"
	"sync"

	"github.com/gmelum/socket/entity/connect"
	"github.com/gmelum/socket/utils/pool"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
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

	upgrader := &ws.Upgrader{
		ReadBufferSize:  soc.readBufferSize,
		WriteBufferSize: soc.writeBufferSize,
		Protocol:        soc.protocol,
		ProtocolCustom:  soc.protocolCustom,
		Extension:       soc.extension,
		ExtensionCustom: soc.extensionCustom,
		Negotiate:       soc.negotiate,
		Header:          soc.header,
		OnRequest:       soc.onRequest,
		OnHost:          soc.onHost,
		OnHeader: func(key, data []byte) (err error) {
			if bytes.Equal(key, []byte("Origin")) {
				for _, clb := range soc.eventsCors {
					clb(data)
				}
				return
			}
			if soc.onHeader != nil {
				return soc.onHeader(key, data)
			}
			return nil
		},
		OnBeforeUpgrade: soc.onBeforeUpgrade,
	}

	outerPool := pool.New(soc.outerWorkers, soc.outerQueue)
	innerPool := pool.New(soc.innerWorkers, soc.innerQueue)

	poller.Start(acceptDesc, func(e netpoll.Event) {

		once := new(sync.Once)
		poller.Resume(acceptDesc)

		conn, err := connect.New(
			listener,
			upgrader,
			poller,
			innerPool,
		)
		if err != nil {
			for _, clb := range soc.eventReject {
				clb(conn)
			}
			return
		}

		if outerPool.NotReady() {
			for _, clb := range soc.eventReject {
				clb(conn)
			}
			return
		}

		conn.Event(func(ev netpoll.Event) {

			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				once.Do(conn.Close)
				return
			}

			if conn.Status == "OPEN" {
				err := innerPool.Schedule(func() {

				})
				if err != nil {
					once.Do(conn.Close)
					return
				}
			}

			once.Do(conn.Close)

		})

		for _, clb := range soc.eventsConnect {
			clb(conn)
		}

		// if soc.maxConn > 0 && soc.maxConn <= soc.UserList.Length() {
		// 	s.reject(listener)
		// 	return
		// }

		// conn, err := listener.Accept()
		// if err != nil {
		// 	return
		// }

		// _, err = upgrader.Upgrade(conn)
		// if err != nil {
		// 	return
		// }

		// desc := netpoll.Must(netpoll.HandleRead(conn))

		// user := client.New(conn, pool)
		// user.Close = func() {
		// 	poller.Stop(desc)
		// 	user.Conn.Close()
		// 	desc.Close()
		// 	for _, clb := range soc.eventsDisconnect {
		// 		clb(user)
		// 	}
		// }

		// poller.Start(desc, func(ev netpoll.Event) {
		// 	if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
		// 		once.Do(user.Close)
		// 		return
		// 	}
		// 	pool.Schedule(func() {
		// 		// err := user.Read()
		// 	})
		// })

	})

	select {}

}
