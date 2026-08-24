package proxy

import (
	"io"
	"net"

	"gateway/internal/session"
)

func PipeToReality(sess *session.Session, sockPath string) error {
	unix, err := net.Dial("unix", sockPath)
	if err != nil {
		return err
	}
	defer unix.Close()
	defer sess.Conn.Close()

	errCh := make(chan error, 2)

	go func() {
		_, err := io.Copy(unix, sess.Conn)
		errCh <- err
	}()

	go func() {
		_, err := io.Copy(sess.Conn, unix)
		errCh <- err
	}()

	return <-errCh
}
