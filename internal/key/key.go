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
	clientIDBytes = 16
	secretBytes   = 32
)

func randomBytes(size int) ([]byte, error) {
	data := make([]byte, size)

	if _, err := rand.Read(data); err != nil {
		return nil, err
	}

	return data, nil
}

func Create(endpoint string) (string, error) {
	endpoint = strings.TrimSpace(endpoint)

	if endpoint == "" {
		return "", errors.New("endpoint is empty")
	}

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

	clientID, err := randomBytes(clientIDBytes)

	if err != nil {
		return "", fmt.Errorf("generate client id: %w", err)
	}

	secret, err := randomBytes(secretBytes)

	if err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}

	id := base64.RawURLEncoding.EncodeToString(clientID)
	token := base64.RawURLEncoding.EncodeToString(secret)

	return fmt.Sprintf(
		"opennet://%s/%s/%s",
		u.Host,
		id,
		token,
	), nil
}

func Parse(value string) (string, string, string, error) {
	value = strings.TrimSpace(value)

	if !strings.HasPrefix(value, "opennet://") {
		return "", "", "", errors.New("invalid OpenNet key")
	}

	raw := strings.TrimPrefix(value, "opennet://")

	parts := strings.Split(raw, "/")

	if len(parts) != 3 {
		return "", "", "", errors.New("invalid OpenNet key format")
	}

	host := parts[0]
	clientID := parts[1]
	token := parts[2]

	if host == "" || clientID == "" || token == "" {
		return "", "", "", errors.New("invalid OpenNet key")
	}

	if _, err := base64.RawURLEncoding.DecodeString(clientID); err != nil {
		return "", "", "", errors.New("invalid client id")
	}

	secret, err := base64.RawURLEncoding.DecodeString(token)

	if err != nil || len(secret) != secretBytes {
		return "", "", "", errors.New("invalid secret")
	}

	return host, clientID, token, nil
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
	fmt.Println("Connect with:")
	fmt.Println("  opennet connect <key>")

	return nil
}
