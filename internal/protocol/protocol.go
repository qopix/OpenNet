package protocol

import (
	"encoding/json"
	"time"
)

const Version = 1

type Handshake struct {
	Version  int    `json:"version"`
	ClientID string `json:"client_id"`
}

type Response struct {
	Version int    `json:"version"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func NewResponse() Response {
	return Response{
		Version: Version,
		Status:  "ok",
		Time:    time.Now().UTC().Format(time.RFC3339),
	}
}

func EncodeResponse() ([]byte, error) {
	return json.Marshal(NewResponse())
}
