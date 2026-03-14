package config

import (
	"os"
	"path/filepath"
)

func validateRuntime(cfg *Config) []Warning {
	var warnings []Warning
	for i, r := range cfg.Reminders {
		if r.Alarm == nil || r.Alarm.Sound == nil || r.Alarm.Sound.File == "" {
			continue
		}
		path := r.Alarm.Sound.File
		if !filepath.IsAbs(path) {
			path = filepath.Clean(path)
		}
		_, err := os.Stat(path)
		if err != nil {
			warnings = append(warnings, Warning{
				ReminderIndex: i,
				ReminderName:  r.Name,
				Message:       "sound file not found or not readable: " + r.Alarm.Sound.File,
				DisableSound:  true,
			})
			continue
		}
		f, err := os.Open(path)
		if err != nil {
			warnings = append(warnings, Warning{
				ReminderIndex: i,
				ReminderName:  r.Name,
				Message:       "sound file not readable: " + r.Alarm.Sound.File,
				DisableSound:  true,
			})
			continue
		}
		_ = f.Close()
	}
	return warnings
}
