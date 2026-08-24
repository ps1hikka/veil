package session

import (
	"net"

	"gateway/internal/config"
	"gateway/internal/vless"
)

type Session struct {
	Conn   net.Conn
	Header *vless.Header
	Client config.Client
}

func New(conn net.Conn, h *vless.Header, c config.Client) *Session {
	return &Session{
		Conn:   conn,
		Header: h,
		Client: c,
	}
}
