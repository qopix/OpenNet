package transport

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

type Transport interface {
	Name() string
	Dial(context.Context, string) (net.Conn, error)
}

type TCP struct {
	Timeout time.Duration
}

func (t TCP) Name() string {
	return "tcp"
}

func (t TCP) Dial(ctx context.Context, address string) (net.Conn, error) {
	timeout := t.Timeout

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	return dialer.DialContext(ctx, "tcp", address)
}

type TLS struct {
	Config  *tls.Config
	Timeout time.Duration
}

func (t TLS) Name() string {
	return "tls"
}

func (t TLS) Dial(ctx context.Context, address string) (net.Conn, error) {
	timeout := t.Timeout

	if timeout == 0 {
		timeout = 10 * time.Second
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	config := t.Config

	if config == nil {
		config = &tls.Config{}
	}

	return tls.DialWithDialer(
		dialer,
		"tcp",
		address,
		config,
	)
}
