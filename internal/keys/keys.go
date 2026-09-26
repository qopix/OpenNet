package keys

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

type KeyPair struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}

func Generate() (*KeyPair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}

	return &KeyPair{
		PublicKey:  pub,
		PrivateKey: priv,
	}, nil
}

func Save(path string, pair *KeyPair) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data := base64.StdEncoding.EncodeToString(pair.PrivateKey)

	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		return err
	}

	return nil
}

func Load(path string) (*KeyPair, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}

	if len(raw) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key")
	}

	priv := ed25519.PrivateKey(raw)
	pub := priv.Public().(ed25519.PublicKey)

	return &KeyPair{
		PublicKey:  pub,
		PrivateKey: priv,
	}, nil
}