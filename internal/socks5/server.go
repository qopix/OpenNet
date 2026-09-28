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
