package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes cfg as JSON to path. Creates parent directory (e.g. ~/.reminder-app) if it does not exist.
func Save(path string, cfg *Config) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}
