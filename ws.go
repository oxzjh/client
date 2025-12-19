package client

import (
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ws struct {
	sync.Mutex
	conn *websocket.Conn
}

func (w *ws) Read() ([]byte, error) {
	_, p, err := w.conn.ReadMessage()
	return p, err
}

func (w *ws) Write(data []byte) error {
	w.Lock()
	defer w.Unlock()
	return w.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (w *ws) Close() error {
	return w.conn.Close()
}

func NewWS(addr string, header http.Header, proxy string) (IClient, error) {
	dialer := websocket.DefaultDialer
	if proxy != "" {
		dialer = &websocket.Dialer{
			Proxy: func(r *http.Request) (*url.URL, error) {
				return url.Parse(proxy)
			},
			HandshakeTimeout: 45 * time.Second,
		}
	}
	conn, _, err := dialer.Dial(addr, header)
	if err != nil {
		return nil, err
	}
	return &ws{conn: conn}, nil
}
