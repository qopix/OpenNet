package web

import (
	"encoding/json"
	"net/http"

	"github.com/qopix/OpenNet/internal/keys"
)

type Response struct {
	Key string `json:"key"`
}

func Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/key", keyHandler)

	return http.ListenAndServe("127.0.0.1:8765", mux)
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key, err := keys.Generate()
	if err != nil {
		http.Error(w, "key generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Response{Key: key})
}