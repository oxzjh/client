package client

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/oxzjh/server"
)

type tcp struct {
	sync.Mutex
	conn net.Conn
	opts *options
}

func (t *tcp) Read() ([]byte, error) {
	return t.opts.parser(t.conn)
}

func (t *tcp) Write(data []byte) error {
	t.Lock()
	defer t.Unlock()
	_, err := t.conn.Write(t.opts.maker(len(data)))
	if err == nil {
		_, err = t.conn.Write(data)
	}
	return err
}

func (t *tcp) WriteJson(v any) error {
	data, _ := json.Marshal(v)
	return t.Write(data)
}

func (t *tcp) Close() error {
	return t.conn.Close()
}

func NewTCP(addr string, opts ...Option) (IClient, error) {
	os := &options{
		maker:  server.MakeStream,
		parser: server.ParseStream4,
	}
	for _, opt := range opts {
		opt(os)
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcp{conn: conn, opts: os}, nil
}
