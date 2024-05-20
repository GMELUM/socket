package connect

import (
	"net"
	"sync"

	"github.com/gmelum/socket/utils/pool"

	"github.com/gobwas/ws"
	"github.com/google/uuid"

	"github.com/mailru/easygo/netpoll"
)

type Connect struct {
	ID     string
	Status string

	mutexRead  sync.Mutex
	mutexWrite sync.Mutex

	conn   net.Conn
	pool   *pool.Pool
	desc   *netpoll.Desc
	poller netpoll.Poller
}

func New(
	listener net.Listener,
	upgrader *ws.Upgrader,
	poller netpoll.Poller,
	pool *pool.Pool,
) (*Connect, error) {

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	_, err = upgrader.Upgrade(conn)
	if err != nil {
		return nil, err
	}

	desc := netpoll.Must(netpoll.HandleRead(conn))

	return &Connect{
		ID:     uuid.New().String(),
		Status: "HANDSHAKE",
		conn:   conn,
		desc:   desc,
		poller: poller,
		pool:   pool,
	}, nil

}

func (conn *Connect) Event(callback func(ev netpoll.Event)) {
	conn.poller.Start(conn.desc, callback)
}

func (conn *Connect) Close() {
	conn.poller.Stop(conn.desc)
	conn.conn.Close()
	conn.desc.Close()
}
