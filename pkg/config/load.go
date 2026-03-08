package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadFile reads the config file at path and returns raw bytes.
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return data, nil
}

// ParseJSON parses raw JSON bytes into Config. Performs JSON validation only:
// valid JSON and root must be an object. Use Validate for full validation.
func ParseJSON(data []byte) (*Config, error) {
	var raw map[string]interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON structure: %w", err)
	}
	return &cfg, nil
}

// LoadFromBytes parses config from bytes (e.g. embedded sample). Same as ParseJSON.
func LoadFromBytes(data []byte) (*Config, error) {
	return ParseJSON(data)
}

// Load reads the file at path, parses JSON, and returns config. Does not validate.
func Load(path string) (*Config, error) {
	data, err := ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseJSON(data)
}

// IsReadable reports whether the path exists and is readable.
func IsReadable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

// IsWritable reports whether the path can be written (file exists and is writable, or parent dir is writable).
func IsWritable(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return fileWritable(path)
	}
	return dirWritable(filepath.Dir(path))
}

func fileWritable(path string) bool {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

func dirWritable(dir string) bool {
	_, err := os.Stat(dir)
	if err != nil {
		return false
	}
	f, err := os.CreateTemp(dir, ".reminder-app-write-test-")
	if err != nil {
		return false
	}
	name := f.Name()
	_ = f.Close()
	_ = os.Remove(name)
	return true
}
