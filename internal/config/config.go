package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	ListenAddress string `json:"listen_address"`
	PanelAddress  string `json:"panel_address"`
}

func Default() Config {
	return Config{
		ListenAddress: "127.0.0.1:9000",
		PanelAddress:  "127.0.0.1:8765",
	}
}

func Save(path string, cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}