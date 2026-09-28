package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/qopix/OpenNet/internal/key"
	"github.com/qopix/OpenNet/internal/protocol"
)

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/opennet/handshake", handshake)

	server := &http.Server{
		Addr:              ":8765",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Println("OpenNet v0.7 node")
	fmt.Println()
	fmt.Println("Listening on:")
	fmt.Println("  http://0.0.0.0:8765")
	fmt.Println()
	fmt.Println("Endpoint:")
	fmt.Println("  /opennet/handshake")
	fmt.Println()

	return server.ListenAndServe()
}

func handshake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	body, err := io.ReadAll(
		http.MaxBytesReader(w, r.Body, 4096),
	)

	if err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)
		return
	}

	var request protocol.Handshake

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(
			w,
			"invalid handshake",
			http.StatusBadRequest,
		)
		return
	}

	if request.Version != protocol.Version {
		http.Error(
			w,
			"unsupported protocol version",
			http.StatusBadRequest,
		)
		return
	}

	if strings.TrimSpace(request.ClientID) == "" {
		http.Error(
			w,
			"missing client id",
			http.StatusUnauthorized,
		)
		return
	}

	response, err := protocol.EncodeResponse()

	if err != nil {
		http.Error(
			w,
			"internal error",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(response)
}

func Connect(value string) error {
	parsed, err := key.Parse(value)

	if err != nil {
		return err
	}

	fmt.Println("OpenNet connection")
	fmt.Println()
	fmt.Println("Node:", parsed.Endpoint)
	fmt.Println("Client:", parsed.ClientID)
	fmt.Println()
	fmt.Println("Sending handshake...")

	payload := protocol.Handshake{
		Version:  protocol.Version,
		ClientID: parsed.ClientID,
	}

	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()

	resp, err := client.Post(
		strings.TrimRight(parsed.Endpoint, "/")+"/opennet/handshake",
		"application/json",
		bytes.NewReader(data),
	)

	latency := time.Since(start)

	if err != nil {
		return fmt.Errorf(
			"node unreachable: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"node rejected handshake: HTTP %d",
			resp.StatusCode,
		)
	}

	var response protocol.Response

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	fmt.Printf(
		"Handshake OK (%dms)\n",
		latency.Milliseconds(),
	)

	fmt.Println("Protocol:", response.Version)
	fmt.Println("Node status:", response.Status)

	return nil
}
