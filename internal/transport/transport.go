package transport

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

func DialTLS(ctx context.Context, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	return tls.DialWithDialer(
		dialer,
		"tcp",
		address,
		&tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	)
}
