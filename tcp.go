package client

import "net"

type tcp struct {
	conn   net.Conn
	buffer []byte
}

func (t *tcp) Read() ([]byte, error) {
	n, err := t.conn.Read(t.buffer)
	return t.buffer[:n], err
}

func (t *tcp) Write(data []byte) error {
	_, err := t.conn.Write(data)
	return err
}

func (t *tcp) Close() error {
	return t.conn.Close()
}

func NewTCP(addr string, bufferSize uint16) (IClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcp{conn, make([]byte, bufferSize)}, nil
}
