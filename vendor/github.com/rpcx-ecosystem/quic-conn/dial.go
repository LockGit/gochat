package quicconn

import (
	"crypto/tls"
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

var quicListen = quic.Listen

// Listen creates a QUIC listener on the given network interface
func Listen(network, laddr string, tlsConfig *tls.Config) (net.Listener, error) {
	udpAddr, err := net.ResolveUDPAddr(network, laddr)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: network, Source: nil, Addr: nil, Err: err}
	}
	conn, err := net.ListenUDP(network, udpAddr)
	if err != nil {
		return nil, err
	}

	ln, err := quicListen(conn, tlsConfig, nil)
	if err != nil {
		return nil, err
	}
	return &server{
		quicServer: ln,
	}, nil
}

// Dial creates a new QUIC connection
// it returns once the connection is established and secured with forward-secure keys
func Dial(addr string, tlsConfig *tls.Config) (net.Conn, error) {
	// DialAddr returns once a forward-secure connection is established
	quicSession, err := quic.DialAddr(addr, tlsConfig, nil)
	if err != nil {
		return nil, err
	}

	sendStream, err := quicSession.OpenStream()
	if err != nil {
		return nil, err
	}

	return &conn{
		session:    quicSession,
		sendStream: sendStream,
	}, nil
}
