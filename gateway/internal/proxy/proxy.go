package proxy

import (
	"fmt"
	"io"
	"net"
	"time"

	"gateway/internal/session"
)

func Handle(sess *session.Session) error {
	target := fmt.Sprintf("%s:%d", sess.Header.Address, sess.Header.Port)

	dest, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		return fmt.Errorf("dial %s: %w", target, err)
	}
	defer dest.Close()

	errCh := make(chan error, 2)

	go func() {
		_, err := io.Copy(dest, sess.Conn)
		errCh <- err
	}()

	go func() {
		_, err := io.Copy(sess.Conn, dest)
		errCh <- err
	}()

	return <-errCh
}
