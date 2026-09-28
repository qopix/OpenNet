package config

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	WorkerURL string `json:"worker_url"`
	Token     string `json:"token"`
	SOCKSHost string `json:"socks_host"`
	SOCKSPort int    `json:"socks_port"`
}

func Load(path string) (Config, error) {
	var c Config

	data, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}

	if c.WorkerURL == "" || c.WorkerURL == "wss://YOUR-WORKER.workers.dev" {
		return c, errors.New("worker_url is not configured")
	}

	if c.SOCKSHost == "" {
		c.SOCKSHost = "127.0.0.1"
	}

	if c.SOCKSPort == 0 {
		c.SOCKSPort = 1080
	}

	return c, nil
}
