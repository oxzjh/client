package client

import (
	"encoding/json"
	"net"
)

type udp struct {
	conn   net.Conn
	buffer []byte
}

func (u *udp) Read() ([]byte, error) {
	n, err := u.conn.Read(u.buffer)
	return u.buffer[:n], err
}

func (u *udp) Write(data []byte) error {
	_, err := u.conn.Write(data)
	return err
}

func (u *udp) WriteJson(v any) error {
	data, _ := json.Marshal(v)
	return u.Write(data)
}

func (u *udp) Close() error {
	return u.conn.Close()
}

func NewUDP(addr string, bufferSize uint16) (IClient, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}
	return &udp{conn, make([]byte, bufferSize)}, nil
}
