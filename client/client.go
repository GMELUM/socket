package client

import (
	"bytes"
	"net"
	"sync"

	"github.com/gmelum/socket/utils/pool"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

type Client struct {
	ID         string
	mutexRead  sync.Mutex
	mutexWrite sync.Mutex

	Conn net.Conn
	pool *pool.Pool

	events []func(...interface{})

	Close func()
}

func New(
	conn net.Conn,
	pool *pool.Pool,
) *Client {

	id := uuid.New()

	return &Client{
		ID:   id.String(),
		Conn: conn,
		pool: pool,
	}

}

func (client *Client) Read() error {

	client.mutexRead.Lock()
	defer client.mutexRead.Unlock()

	header, reader, err := wsutil.NextReader(client.Conn, ws.StateServerSide)
	if err != nil {
		return err
	}

	if header.OpCode.IsControl() {
		handler := wsutil.ControlFrameHandler(client.Conn, ws.StateServerSide)
		return handler(header, reader)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return err
	}

	// decoded := decoding(buf.Bytes())

	// if client.callbackEvents != nil {
	// 	client.callbackEvents(decoded.Type, decoded.Value)
	// }

	return nil

}
