package key

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	clientIDSize = 16
	secretSize   = 32
)

func random(size int) ([]byte, error) {
	data := make([]byte, size)

	if _, err := rand.Read(data); err != nil {
		return nil, err
	}

	return data, nil
}

func Create(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)

	u, err := url.Parse(endpoint)
	if err != nil {
		return "", errors.New("invalid endpoint")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("endpoint must use http or https")
	}

	if u.Host == "" {
		return "", errors.New("endpoint host is missing")
	}

	id, err := random(clientIDSize)
	if err != nil {
		return "", err
	}

	secret, err := random(secretSize)
	if err != nil {
		return "", err
	}

	clientID := base64.RawURLEncoding.EncodeToString(id)
	token := base64.RawURLEncoding.EncodeToString(secret)

	return fmt.Sprintf(
		"opennet://%s/%s/%s",
		u.String(),
		clientID,
		token,
	), nil
}

type Parsed struct {
	Endpoint string
	ClientID string
	Secret   []byte
}

func Parse(value string) (*Parsed, error) {
	if !strings.HasPrefix(value, "opennet://") {
		return nil, errors.New("invalid OpenNet key")
	}

	raw := strings.TrimPrefix(value, "opennet://")

	parts := strings.Split(raw, "/")

	if len(parts) < 4 {
		return nil, errors.New("invalid OpenNet key format")
	}

	scheme := parts[0]
	host := parts[1]
	clientID := parts[2]
	token := parts[3]

	if scheme != "http:" && scheme != "https:" {
		return nil, errors.New("invalid endpoint scheme")
	}

	endpoint := scheme + "//" + host

	id, err := base64.RawURLEncoding.DecodeString(clientID)
	if err != nil || len(id) != clientIDSize {
		return nil, errors.New("invalid client id")
	}

	secret, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(secret) != secretSize {
		return nil, errors.New("invalid secret")
	}

	return &Parsed{
		Endpoint: endpoint,
		ClientID: clientID,
		Secret:   secret,
	}, nil
}

func CreateAndShow() error {
	fmt.Print("Node endpoint (http/https): ")

	var endpoint string

	if _, err := fmt.Scanln(&endpoint); err != nil {
		return err
	}

	value, err := Create(endpoint)
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("OpenNet client key:")
	fmt.Println()
	fmt.Println(value)
	fmt.Println()
	fmt.Println("Connect:")
	fmt.Println("  opennet connect <key>")

	return nil
}
