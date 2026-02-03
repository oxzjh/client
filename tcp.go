package client

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/oxzjh/stream"
)

type tcp struct {
	sync.Mutex
	conn net.Conn
	opts *options
}

func (t *tcp) Read() ([]byte, error) {
	return t.opts.parser.Parse(t.conn)
}

func (t *tcp) Write(data []byte) error {
	t.Lock()
	defer t.Unlock()
	_, err := t.conn.Write(t.opts.maker.Make(len(data)))
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
		maker:  stream.NewMaker(0x92),
		parser: stream.NewSimpleParser(4),
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
