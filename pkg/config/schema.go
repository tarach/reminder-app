package config

import (
	"fmt"
	"strings"
)

func validateReminders(cfg *Config) error {
	for i := range cfg.Reminders {
		err := validateReminder(cfg, i, &cfg.Reminders[i])
		if err != nil {
			return err
		}
	}
	return nil
}

var reminderTypeHandlers = map[string]func(*Config, int, *Reminder) error{
	"alert":   validateReminderAlert,
	"counter": validateReminderCounter,
}

func validateReminder(cfg *Config, idx int, r *Reminder) error {
	if r.Name == "" {
		return fmt.Errorf("reminder %d: missing or empty 'name'", idx+1)
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("reminder %d (%q): 'name' must not be empty", idx+1, r.Name)
	}
	if r.DaysActive < 0 || r.DaysActive > 127 {
		return fmt.Errorf("reminder %d (%q): 'daysActive' must be between 0 and 127, got %d", idx+1, r.Name, r.DaysActive)
	}
	handler, ok := reminderTypeHandlers[r.Type]
	if !ok {
		return fmt.Errorf("reminder %d (%q): 'type' must be \"alert\" or \"counter\", got %q", idx+1, r.Name, r.Type)
	}
	if err := validateAlarm(r); err != nil {
		return fmt.Errorf("reminder %d (%q): %w", idx+1, r.Name, err)
	}
	return handler(cfg, idx, r)
}

func validateAlarm(r *Reminder) error {
	if r.Alarm == nil {
		return nil
	}
	if r.Alarm.Rules != nil && len(r.Alarm.Rules) == 0 {
		return fmt.Errorf("alarm.rules must not be empty when present")
	}
	for i, rule := range r.Alarm.Rules {
		hasAt := rule.At != ""
		hasEvery := rule.Every != ""
		if hasAt == hasEvery {
			return fmt.Errorf("alarm.rules[%d]: must have exactly one of 'at' or 'every'", i)
		}
		if hasAt && rule.At == "" {
			return fmt.Errorf("alarm.rules[%d]: 'at' must be a non-empty string", i)
		}
		if hasEvery && rule.Every == "" {
			return fmt.Errorf("alarm.rules[%d]: 'every' must be a non-empty string", i)
		}
	}
	if r.Alarm.Sound != nil {
		hasDefault := r.Alarm.Sound.Default != nil
		hasFile := r.Alarm.Sound.File != ""
		if hasDefault == hasFile {
			return fmt.Errorf("alarm.sound: must have exactly one of 'default' or 'file'")
		}
	}
	return nil
}

func validateReminderAlert(_ *Config, idx int, r *Reminder) error {
	if r.Alert == nil {
		return fmt.Errorf("reminder %d (%q): type 'alert' requires 'alert' section", idx+1, r.Name)
	}
	if r.Alert.Value == "" {
		return fmt.Errorf("reminder %d (%q): alert.value is required", idx+1, r.Name)
	}
	if r.Counter != nil {
		return fmt.Errorf("reminder %d (%q): type 'alert' must not have 'counter' section", idx+1, r.Name)
	}
	return nil
}

func validateReminderCounter(_ *Config, idx int, r *Reminder) error {
	if r.Counter == nil {
		return fmt.Errorf("reminder %d (%q): type 'counter' requires 'counter' section", idx+1, r.Name)
	}
	if r.Alert != nil {
		return fmt.Errorf("reminder %d (%q): type 'counter' must not have 'alert' section", idx+1, r.Name)
	}
	return validateCounter(r, idx)
}

func validateCounter(r *Reminder, idx int) error {
	c := r.Counter
	if c.Format != "duration" && c.Format != "number" {
		return fmt.Errorf("reminder %d (%q): counter.format must be \"duration\" or \"number\", got %q", idx+1, r.Name, c.Format)
	}
	if len(c.Rules) == 0 {
		return fmt.Errorf("reminder %d (%q): counter.rules is required and must not be empty", idx+1, r.Name)
	}
	valueOK := map[string]func(string) bool{"duration": validDurationValue, "number": validNumberValue}
	for i, rule := range c.Rules {
		if rule.Every == "" {
			return fmt.Errorf("reminder %d (%q): counter.rules[%d] must have 'every'", idx+1, r.Name, i)
		}
		hasInc := rule.Increase != ""
		hasSet := rule.SetValue != ""
		if hasInc == hasSet {
			return fmt.Errorf("reminder %d (%q): counter.rules[%d] must have exactly one of 'increase' or 'setValue'", idx+1, r.Name, i)
		}
		if hasInc && !valueOK[c.Format](rule.Increase) {
			return fmt.Errorf("reminder %d (%q): counter.rules[%d] increase %q is not a valid %s", idx+1, r.Name, i, rule.Increase, c.Format)
		}
		if hasSet && !valueOK[c.Format](rule.SetValue) {
			return fmt.Errorf("reminder %d (%q): counter.rules[%d] setValue %q is not a valid %s", idx+1, r.Name, i, rule.SetValue, c.Format)
		}
	}
	return nil
}
