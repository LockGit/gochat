package quick

import (
	"context"
	"crypto/tls"
	"net"

	quic "github.com/lucas-clemente/quic-go"
)

// Dial creates a new QUIC connection
// it returns once the connection is established and secured with forward-secure keys
func Dial(addr string, tlsConfig *tls.Config, quicConfig *quic.Config) (net.Conn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	// DialAddr returns once a forward-secure connection is established
	quicSession, err := quic.Dial(udpConn, udpAddr, addr, tlsConfig, quicConfig)
	if err != nil {
		return nil, err
	}

	stream, err := quicSession.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}

	return &Conn{
		conn:    udpConn,
		session: quicSession,
		stream:  stream,
	}, nil
}
