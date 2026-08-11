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

// Path renvoie l'emplacement de config.json, a cote de l'executable.
func Path() (string, error) {
	dir, err := executableDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Save ecrit la configuration a cote de l'executable en 0600 (le jeton de
// presence ne doit etre lisible que par l'utilisateur). Utilise par le flux de
// liaison au premier lancement.
func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}

	if cfg.PollIntervalSeconds <= 0 {
		cfg.PollIntervalSeconds = 30
	}

	payload, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, payload, 0o600)
}

func Load() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}

	payload, err := os.ReadFile(path)
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
