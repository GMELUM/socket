package connect

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/alitto/pond"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"

	"github.com/mailru/easygo/netpoll"
)

type Connect struct {
	ID string

	Host    string
	Uri     string
	Headers map[string]string

	mutexRead  sync.Mutex
	mutexWrite sync.Mutex

	conn   net.Conn
	pool   *pond.WorkerPool
	desc   *netpoll.Desc
	poller netpoll.Poller
}

func New(

	listener net.Listener,
	poller netpoll.Poller,
	pool *pond.WorkerPool,

	readBufferSize int,
	writeBufferSize int,
	protocol func([]byte) bool,
	protocolCustom func([]byte) (string, bool),
	extension func(httphead.Option) bool,
	extensionCustom func([]byte, []httphead.Option) ([]httphead.Option, bool),
	negotiate func(httphead.Option) (httphead.Option, error),
	header ws.HandshakeHeader,

	eventsCors *[]func(origin string) (err error),
	eventsConnect *[]func(cn *Connect) (err error),

) (*Connect, error) {

	connect := &Connect{
		ID:      uuid.New().String(),
		poller:  poller,
		pool:    pool,
		Headers: make(map[string]string),
	}

	upgrader := &ws.Upgrader{
		ReadBufferSize:  readBufferSize,
		WriteBufferSize: writeBufferSize,
		Protocol:        protocol,
		ProtocolCustom:  protocolCustom,
		Extension:       extension,
		ExtensionCustom: extensionCustom,
		Negotiate:       negotiate,
		Header:          header,
		OnRequest: func(uri []byte) error {
			connect.Uri = string(uri)
			return nil
		},
		OnHost: func(host []byte) error {
			connect.Host = string(host)
			return nil
		},
		OnHeader: func(key, data []byte) (err error) {
			connect.Headers[strings.ToLower(string(key))] = string(data)
			if bytes.Equal(key, []byte("Origin")) {
				for _, callback := range *eventsCors {
					err := callback(string(data))
					if err != nil {
						return err
					}
				}
			}
			return nil
		},
		OnBeforeUpgrade: func() (header ws.HandshakeHeader, err error) {
			if _, ok := connect.Headers["origin"]; !ok {
				return nil, errors.New("origin is not defined")
			}

			for _, callback := range *eventsConnect {
				err := callback(connect)
				if err != nil {
					return nil, err
				}
			}

			return header, nil
		},
	}

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}

	_, err = upgrader.Upgrade(conn)
	if err != nil {
		return nil, err
	}

	connect.conn = conn
	connect.desc = netpoll.Must(netpoll.HandleRead(conn))

	return connect, nil

}

func (conn *Connect) Event(callback func(ev netpoll.Event)) {
	conn.poller.Start(conn.desc, callback)
}

func (conn *Connect) Close() {
	conn.poller.Stop(conn.desc)
	conn.conn.Close()
	conn.desc.Close()
}

func (conn *Connect) Send(id int, event string, data interface{}) error {

	conn.mutexWrite.Lock()
	defer conn.mutexWrite.Unlock()

	encoded := encoding(Encoding{
		ID:    id,
		Type:  event,
		Value: data,
	})

	buffer := new(bytes.Buffer)
	writer := wsutil.NewWriter(buffer, ws.StateServerSide, ws.OpBinary)
	writer.Write(*encoded)

	if err := writer.Flush(); err != nil {
		return err
	}

	conn.pool.Submit(func() {
		conn.conn.Write(buffer.Bytes())
	})

	return nil

}

func (conn *Connect) Read() (*Decoding, error) {

	conn.mutexRead.Lock()
	defer conn.mutexRead.Unlock()

	header, reader, err := wsutil.NextReader(conn.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}

	if header.OpCode.IsControl() {
		handler := wsutil.ControlFrameHandler(conn.conn, ws.StateServerSide)
		return nil, handler(header, reader)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return decoding(buf.Bytes())

}
