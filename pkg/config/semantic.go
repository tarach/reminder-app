package config

import (
	"fmt"
	"strconv"
	"time"
)

func validateSemantic(cfg *Config) error {
	seen := make(map[string]int)
	for i, r := range cfg.Reminders {
		if j, ok := seen[r.Name]; ok {
			return fmt.Errorf("reminder names must be unique: %q appears at index %d and %d", r.Name, j+1, i+1)
		}
		seen[r.Name] = i
	}
	for i := range cfg.Reminders {
		r := &cfg.Reminders[i]
		if r.Counter == nil {
			continue
		}
		err := validateCounterValues(r, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateCounterValues(r *Reminder, idx int) error {
	c := r.Counter
	validators := map[string]func(string) bool{
		"duration": validDurationValue,
		"number":   validNumberValue,
	}
	ok := validators[c.Format](c.CurrentValue)
	if !ok {
		return fmt.Errorf("reminder %d (%q): counter.currentValue %q is not a valid %s", idx+1, r.Name, c.CurrentValue, c.Format)
	}
	ok = validators[c.Format](c.InitialValue)
	if !ok {
		return fmt.Errorf("reminder %d (%q): counter.initialValue %q is not a valid %s", idx+1, r.Name, c.InitialValue, c.Format)
	}
	return nil
}

func validDurationValue(s string) bool {
	_, err := time.ParseDuration(s)
	if err == nil {
		return true
	}
	return validDurationClock(s)
}

func validDurationClock(s string) bool {
	if len(s) < 4 {
		return false
	}
	parts := splitClock(s)
	if len(parts) != 2 {
		return false
	}
	m, err1 := strconv.Atoi(parts[0])
	n, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return false
	}
	return m >= 0 && n >= 0 && n < 60
}

func splitClock(s string) []string {
	for i, c := range s {
		if c == ':' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}

func validNumberValue(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}
