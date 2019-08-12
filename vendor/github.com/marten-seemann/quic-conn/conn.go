package quicconn

import (
	"net"
	"time"

	quic "github.com/lucas-clemente/quic-go"
)

type conn struct {
	session quic.Session

	receiveStream quic.Stream
	sendStream    quic.Stream
}

func newConn(sess quic.Session) (*conn, error) {
	stream, err := sess.OpenStream()
	if err != nil {
		return nil, err
	}
	return &conn{
		session:    sess,
		sendStream: stream,
	}, nil
}

func (c *conn) Read(b []byte) (int, error) {
	if c.receiveStream == nil {
		var err error
		c.receiveStream, err = c.session.AcceptStream()
		// TODO: check stream id
		if err != nil {
			return 0, err
		}
		// quic.Stream.Close() closes the stream for writing
		err = c.receiveStream.Close()
		if err != nil {
			return 0, err
		}
	}

	return c.receiveStream.Read(b)
}

func (c *conn) Write(b []byte) (int, error) {
	return c.sendStream.Write(b)
}

// LocalAddr returns the local network address.
// needed to fulfill the net.Conn interface
func (c *conn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *conn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

func (c *conn) Close() error {
	return c.session.Close()
}

func (c *conn) SetDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *conn) SetWriteDeadline(t time.Time) error {
	return nil
}

var _ net.Conn = &conn{}
