#!/usr/bin/env bash
set -euo pipefail

echo "================================="
echo " OpenNet v1.0 Bootstrap"
echo "================================="

if [ ! -d ".git" ]; then
    echo "[ERROR] Run this from the OpenNet repository root."
    exit 1
fi

TMP="$(mktemp -d)"

[ -f README.md ] && cp README.md "$TMP/README.md"
[ -f LICENSE ] && cp LICENSE "$TMP/LICENSE"

find . -mindepth 1 -maxdepth 1 \
    ! -name ".git" \
    ! -name "README.md" \
    ! -name "LICENSE" \
    ! -name "bootstrap" \
    -exec rm -rf {} +

mkdir -p \
    cmd/opennet \
    internal/config \
    internal/protocol \
    internal/transport \
    internal/socks5 \
    data/cache \
    worker/src

cat > go.mod <<'EOF'
module github.com/qopix/OpenNet

go 1.23

require nhooyr.io/websocket v1.8.17
EOF

cat > internal/config/config.go <<'EOF'
package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	WorkerURL string `json:"worker_url"`
	Token     string `json:"token"`
	SOCKSHost string `json:"socks_host"`
	SOCKSPort int    `json:"socks_port"`
}

func Load(path string) (Config, error) {
	var c Config

	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}

	if c.WorkerURL == "" || c.WorkerURL == "wss://YOUR-WORKER.workers.dev" {
		return c, errors.New("worker_url is not configured")
	}

	if c.SOCKSHost == "" {
		c.SOCKSHost = "127.0.0.1"
	}

	if c.SOCKSPort == 0 {
		c.SOCKSPort = 1080
	}

	return c, nil
}
EOF

cat > internal/protocol/protocol.go <<'EOF'
package protocol

const (
	Version byte = 1
	OpConnect byte = 1
	OpData byte = 2
	OpClose byte = 3
)
EOF

cat > internal/transport/transport.go <<'EOF'
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
EOF

cat > internal/socks5/server.go <<'EOF'
package socks5

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"

	"github.com/qopix/OpenNet/internal/config"
	"github.com/qopix/OpenNet/internal/transport"
	"nhooyr.io/websocket"
)

type Server struct {
	cfg       config.Config
	transport *transport.Transport
}

func New(cfg config.Config) *Server {
	return &Server{
		cfg:       cfg,
		transport: transport.New(cfg),
	}
}

func (s *Server) ListenAndServe() error {
	addr := net.JoinHostPort(
		s.cfg.SOCKSHost,
		strconv.Itoa(s.cfg.SOCKSPort),
	)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Printf("[OpenNet] SOCKS5 listening on %s\n", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	if err := handshake(conn); err != nil {
		return
	}

	host, port, err := readRequest(conn)
	if err != nil {
		return
	}

	ws, err := s.transport.Connect(context.Background())
	if err != nil {
		writeFailure(conn)
		return
	}

	defer ws.Close(websocket.StatusNormalClosure, "")

	target := fmt.Sprintf("CONNECT %s:%d", host, port)

	if err := ws.Write(
		context.Background(),
		websocket.MessageText,
		[]byte(target),
	); err != nil {
		writeFailure(conn)
		return
	}

	typ, data, err := ws.Read(context.Background())
	if err != nil ||
		typ != websocket.MessageText ||
		string(data) != "OK" {
		writeFailure(conn)
		return
	}

	writeSuccess(conn)

	proxy(conn, ws)
}

func handshake(conn net.Conn) error {
	header := make([]byte, 2)

	if _, err := io.ReadFull(conn, header); err != nil {
		return err
	}

	if header[0] != 5 {
		return errors.New("not SOCKS5")
	}

	methods := make([]byte, int(header[1]))

	if _, err := io.ReadFull(conn, methods); err != nil {
		return err
	}

	_, err := conn.Write([]byte{5, 0})
	return err
}

func readRequest(conn net.Conn) (string, uint16, error) {
	header := make([]byte, 4)

	if _, err := io.ReadFull(conn, header); err != nil {
		return "", 0, err
	}

	if header[0] != 5 || header[1] != 1 {
		return "", 0, errors.New("unsupported SOCKS request")
	}

	var host string

	switch header[3] {
	case 1:
		ip := make([]byte, 4)

		if _, err := io.ReadFull(conn, ip); err != nil {
			return "", 0, err
		}

		host = net.IP(ip).String()

	case 3:
		length := make([]byte, 1)

		if _, err := io.ReadFull(conn, length); err != nil {
			return "", 0, err
		}

		name := make([]byte, int(length[0]))

		if _, err := io.ReadFull(conn, name); err != nil {
			return "", 0, err
		}

		host = string(name)

	case 4:
		ip := make([]byte, 16)

		if _, err := io.ReadFull(conn, ip); err != nil {
			return "", 0, err
		}

		host = net.IP(ip).String()

	default:
		return "", 0, errors.New("unknown address type")
	}

	port := make([]byte, 2)

	if _, err := io.ReadFull(conn, port); err != nil {
		return "", 0, err
	}

	return host, binary.BigEndian.Uint16(port), nil
}

func writeSuccess(conn net.Conn) {
	_, _ = conn.Write([]byte{
		5, 0, 0, 1,
		0, 0, 0, 0,
		0, 0,
	})
}

func writeFailure(conn net.Conn) {
	_, _ = conn.Write([]byte{
		5, 1, 0, 1,
		0, 0, 0, 0,
		0, 0,
	})
}

func proxy(conn net.Conn, ws *websocket.Conn) {
	ctx := context.Background()
	errc := make(chan error, 2)

	go func() {
		buf := make([]byte, 32*1024)

		for {
			n, err := conn.Read(buf)

			if n > 0 {
				if werr := ws.Write(
					ctx,
					websocket.MessageBinary,
					append([]byte(nil), buf[:n]...),
				); werr != nil {
					errc <- werr
					return
				}
			}

			if err != nil {
				errc <- err
				return
			}
		}
	}()

	go func() {
		for {
			typ, data, err := ws.Read(ctx)

			if err != nil {
				errc <- err
				return
			}

			if typ != websocket.MessageBinary {
				continue
			}

			if _, err := conn.Write(data); err != nil {
				errc <- err
				return
			}
		}
	}()

	<-errc
}
EOF

cat > cmd/opennet/main.go <<'EOF'
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/qopix/OpenNet/internal/config"
	"github.com/qopix/OpenNet/internal/socks5"
)

const configPath = "data/config.json"

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Println("OpenNet v1.0")
			return

		case "status":
			fmt.Println("OpenNet v1.0")
			fmt.Println("SOCKS5: 127.0.0.1:1080")
			return

		case "stop":
			fmt.Println("OpenNet v1.0 runs in foreground.")
			fmt.Println("Stop it with Ctrl+C.")
			return
		}
	}

	fmt.Println("=================================")
	fmt.Println(" OpenNet v1.0")
	fmt.Println("=================================")

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}

	fmt.Println("[OpenNet] Configuration loaded.")
	fmt.Println("[OpenNet] Starting transport...")
	fmt.Println("[OpenNet] Starting SOCKS5...")

	server := socks5.New(cfg)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("[OpenNet] SOCKS5 error: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Println("[OpenNet] Shield: ON")
	fmt.Println("[OpenNet] SOCKS5: 127.0.0.1:1080")
	fmt.Println("OpenNet is running.")
	fmt.Println("Press Ctrl+C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	fmt.Println("\n[OpenNet] Stopped.")
}
EOF

cat > data/config.json <<'EOF'
{
  "worker_url": "wss://YOUR-WORKER.workers.dev",
  "token": "CHANGE-ME",
  "socks_host": "127.0.0.1",
  "socks_port": 1080
}
EOF

cat > worker/wrangler.jsonc <<'EOF'
{
  "$schema": "node_modules/wrangler/config-schema.json",
  "name": "opennet-v1",
  "main": "src/index.js",
  "compatibility_date": "2026-09-28"
}
EOF

cat > worker/src/index.js <<'EOF'
import { connect } from "cloudflare:sockets";

export default {
  async fetch(request, env) {
    if (request.headers.get("Upgrade")?.toLowerCase() !== "websocket") {
      return new Response("OpenNet transport\n", { status: 426 });
    }

    const auth = request.headers.get("Authorization") || "";

    if (!env.OPENNET_TOKEN || auth !== `Bearer ${env.OPENNET_TOKEN}`) {
      return new Response("Unauthorized\n", { status: 401 });
    }

    const pair = new WebSocketPair();
    const client = pair[0];
    const server = pair[1];

    server.accept();

    let socket = null;
    let writer = null;
    let connected = false;

    const closeAll = () => {
      try {
        writer?.releaseLock();
      } catch {}

      try {
        socket?.close();
      } catch {}

      try {
        server.close();
      } catch {}
    };

    server.addEventListener("message", async (event) => {
      try {
        if (typeof event.data === "string") {
          if (!event.data.startsWith("CONNECT ")) {
            server.close(1002, "bad command");
            return;
          }

          const target = event.data.slice(8);
          const separator = target.lastIndexOf(":");

          if (separator <= 0) {
            server.close(1002, "bad target");
            return;
          }

          const hostname = target.slice(0, separator);
          const port = Number(target.slice(separator + 1));

          if (
            !hostname ||
            !Number.isInteger(port) ||
            port < 1 ||
            port > 65535
          ) {
            server.close(1002, "bad target");
            return;
          }

          socket = connect({
            hostname,
            port
          });

          writer = socket.writable.getWriter();
          connected = true;

          server.send("OK");

          const reader = socket.readable.getReader();

          try {
            while (true) {
              const { value, done } = await reader.read();

              if (done) break;

              if (value) {
                server.send(value);
              }
            }
          } finally {
            reader.releaseLock();
            closeAll();
          }

          return;
        }

        if (!connected || !writer) {
          server.close(1002, "not connected");
          return;
        }

        let data;

        if (event.data instanceof ArrayBuffer) {
          data = new Uint8Array(event.data);
        } else if (event.data instanceof Blob) {
          data = new Uint8Array(await event.data.arrayBuffer());
        } else {
          return;
        }

        await writer.write(data);
      } catch (error) {
        console.log(String(error));
        closeAll();
      }
    });

    server.addEventListener("close", closeAll);
    server.addEventListener("error", closeAll);

    return new Response(null, {
      status: 101,
      webSocket: client
    });
  }
};
EOF

cat > worker/README.md <<'EOF'
# OpenNet Worker

Install Wrangler:

npm install -g wrangler

Login:

wrangler login

Create a token:

openssl rand -hex 32

Set the token:

wrangler secret put OPENNET_TOKEN

Deploy:

wrangler deploy
EOF

cat > .gitignore <<'EOF'
opennet
data/cache/
data/config.local.json
*.log
worker/.wrangler/
EOF

[ -f "$TMP/README.md" ] && cp "$TMP/README.md" README.md
[ -f "$TMP/LICENSE" ] && cp "$TMP/LICENSE" LICENSE

rm -rf "$TMP"

gofmt -w cmd internal

echo "[OpenNet] Installing dependencies..."
go mod tidy

echo "[OpenNet] Building..."
go build -o opennet ./cmd/opennet

chmod 700 opennet

echo
echo "================================="
echo " OpenNet v1.0 READY"
echo "================================="
echo
echo "Worker:"
echo "  cd worker"
echo "  wrangler secret put OPENNET_TOKEN"
echo "  wrangler deploy"
echo
echo "Then edit:"
echo "  data/config.json"
echo
echo "Run:"
echo "  ./opennet"
echo
echo "Bootstrap finished."

rm -f -- "$0"
