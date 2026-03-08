package config

import (
	"os"
	"path/filepath"
)

const (
	configFileName   = "config.json"
	homeConfigSubdir = ".reminder-app"
)

// LookupResult holds the result of config path lookup.
type LookupResult struct {
	Path string
	OK   bool
}

// Lookup searches for config in home directory first, then current working directory.
// Home has priority. Returns the first path where the file exists, or empty + OK false.
func Lookup() LookupResult {
	homePath := homeConfigPath()
	_, err := os.Stat(homePath)
	if err == nil {
		return LookupResult{Path: homePath, OK: true}
	}
	cwdPath := filepath.Join(".", configFileName)
	_, err = os.Stat(cwdPath)
	if err == nil {
		return LookupResult{Path: cwdPath, OK: true}
	}
	return LookupResult{OK: false}
}

// HomeConfigPath returns the path where config is stored in the user's home directory.
func HomeConfigPath() string {
	return homeConfigPath()
}

func homeConfigPath() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, homeConfigSubdir, configFileName)
}

// HomeConfigDir returns the directory for home config (e.g. ~/.reminder-app).
func HomeConfigDir() string {
	dir, _ := os.UserHomeDir()
	return filepath.Join(dir, homeConfigSubdir)
}
