package web

import (
	"encoding/json"
	"net/http"

	"github.com/qopix/OpenNet/internal/keys"
)

type Server struct {
	Address string
}

type KeyResponse struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func New(address string) *Server {
	return &Server{
		Address: address,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/api/key", keyHandler)

	return http.ListenAndServe(s.Address, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/index.html")
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pair, err := keys.Generate()
	if err != nil {
		http.Error(w, "failed to generate key", http.StatusInternalServerError)
		return
	}

	response := KeyResponse{
		PublicKey:  pair.PublicKey,
		PrivateKey: pair.PrivateKey,
	}

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
}