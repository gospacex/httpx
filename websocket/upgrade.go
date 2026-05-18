package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Upgrader struct {
	ReadBufferSize  int
	WriteBufferSize int
	CheckOrigin    func(*http.Request) bool
}

var DefaultUpgrader = &Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  u.ReadBufferSize,
		WriteBufferSize: u.WriteBufferSize,
		CheckOrigin:     u.CheckOrigin,
	}
	return upgrader.Upgrade(w, r, nil)
}

func UpgradeHTTP(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return DefaultUpgrader.Upgrade(w, r)
}