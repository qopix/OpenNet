package transport

import (
	"context"
	"net"
)

type Transport interface {
	Name() string
	Dial(ctx context.Context, address string) (net.Conn, error)
}

type Direct struct{}

func (Direct) Name() string {
	return "direct"
}

func (Direct) Dial(ctx context.Context, address string) (net.Conn, error) {
	var dialer net.Dialer

	return dialer.DialContext(ctx, "tcp", address)
}