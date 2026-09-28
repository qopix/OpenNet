package transport

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/qopix/OpenNet/internal/config"
	"nhooyr.io/websocket"
)

type Transport struct {
	cfg config.Config
}

func New(cfg config.Config) *Transport {
	return &Transport{cfg: cfg}
}

func (t *Transport) Connect(ctx context.Context) (*websocket.Conn, error) {
	u, err := url.Parse(t.cfg.WorkerURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "ws" && u.Scheme != "wss" {
		return nil, fmt.Errorf("worker URL must use ws:// or wss://")
	}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	headers := make(map[string][]string)

	if t.cfg.Token != "" {
		headers["Authorization"] = []string{
			"Bearer " + t.cfg.Token,
		}
	}

	conn, _, err := websocket.Dial(
		ctx,
		t.cfg.WorkerURL,
		&websocket.DialOptions{
			HTTPHeader: headers,
		},
	)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
