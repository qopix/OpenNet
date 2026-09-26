package transport

import (
	"context"
	"crypto/tls"
	"net"
)

type Transport interface {
	Name() string
	Dial(context.Context, string) (net.Conn, error)
}

type TLS struct {
	Config *tls.Config
}

func (t TLS) Name() string {
	return "tls"
}

func (t TLS) Dial(ctx context.Context, address string) (net.Conn, error) {
	dialer := &net.Dialer{}

	return tls.DialWithDialer(
		dialer,
		"tcp",
		address,
		t.Config,
	)
}