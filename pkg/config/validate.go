package config

import "fmt"

// Validator is a function that validates config and returns an error or nil.
type Validator func(*Config) error

// Validate runs all validation stages. Returns any runtime warnings (e.g. sound file missing) and the first validation error.
func Validate(cfg *Config) (warnings []Warning, err error) {
	validators := []Validator{
		validateRoot,
		validateReminders,
		validateSemantic,
	}
	for _, v := range validators {
		err = v(cfg)
		if err != nil {
			return nil, err
		}
	}
	warnings = validateRuntime(cfg)
	return warnings, nil
}

func validateRoot(cfg *Config) error {
	if cfg.Format == "" {
		return fmt.Errorf("configuration: missing required field 'format'")
	}
	formatOK := cfg.Format == "12h" || cfg.Format == "24h"
	if !formatOK {
		return fmt.Errorf("configuration: 'format' must be \"12h\" or \"24h\", got %q", cfg.Format)
	}
	if cfg.Reminders == nil {
		return fmt.Errorf("configuration: missing required field 'reminders'")
	}
	return nil
}
