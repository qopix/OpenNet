package web

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/qopix/OpenNet/internal/keys"
)

type Server struct {
	Address string
	KeyPath string
}

type KeyResponse struct {
	PublicKey string `json:"public_key"`
}

func New(address, keyPath string) *Server {
	return &Server{
		Address: address,
		KeyPath: keyPath,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/api/key", s.generateKey)

	return http.ListenAndServe(s.Address, mux)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, "web/index.html")
}

func (s *Server) generateKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pair, err := keys.Generate()

	if err != nil {
		http.Error(w, "failed to generate key", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(filepath.Dir(s.KeyPath), 0700); err != nil {
		http.Error(w, "failed to create key directory", http.StatusInternalServerError)
		return
	}

	if err := keys.Save(s.KeyPath, pair); err != nil {
		http.Error(w, "failed to save key", http.StatusInternalServerError)
		return
	}

	response := KeyResponse{
		PublicKey: base64.StdEncoding.EncodeToString(pair.PublicKey),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
