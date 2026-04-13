package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	PresenceToken       string `json:"presenceToken"`
	StatusURL           string `json:"statusUrl"`
	PollIntervalSeconds int    `json:"pollIntervalSeconds"`
}

func executableDir() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(executablePath), nil
}

func Load() (Config, error) {
	dir, err := executableDir()
	if err != nil {
		return Config{}, err
	}

	payload, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(payload, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.PollIntervalSeconds <= 0 {
		cfg.PollIntervalSeconds = 30
	}

	if cfg.PresenceToken == "" {
		return Config{}, errors.New("missing presenceToken in config.json")
	}
	if cfg.StatusURL == "" {
		return Config{}, errors.New("missing statusUrl in config.json")
	}

	return cfg, nil
}
