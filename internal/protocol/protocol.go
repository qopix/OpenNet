package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	Version byte = 1

	TypeHello byte = 1
	TypeData  byte = 2
	TypeClose byte = 3
)

const HeaderSize = 6

func Write(w io.Writer, typ byte, payload []byte) error {
	if len(payload) > 16*1024*1024 {
		return fmt.Errorf("payload too large")
	}

	header := make([]byte, HeaderSize)

	header[0] = Version
	header[1] = typ

	binary.BigEndian.PutUint32(
		header[2:],
		uint32(len(payload)),
	)

	if _, err := w.Write(header); err != nil {
		return err
	}

	_, err := w.Write(payload)

	return err
}

func Read(r io.Reader) (byte, []byte, error) {
	header := make([]byte, HeaderSize)

	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, err
	}

	if header[0] != Version {
		return 0, nil, fmt.Errorf("unsupported protocol version")
	}

	length := binary.BigEndian.Uint32(header[2:])

	if length > 16*1024*1024 {
		return 0, nil, fmt.Errorf("payload too large")
	}

	payload := make([]byte, length)

	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, err
	}

	return header[1], payload, nil
}