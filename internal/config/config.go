package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir  string
	DBPath   string
	Port     string
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dataDir := filepath.Join(home, ".tabgraph")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}
	return &Config{
		DataDir: dataDir,
		DBPath:  filepath.Join(dataDir, "data.db"),
		Port:    "8080",
	}, nil
}
