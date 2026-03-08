package config

import (
	"strings"
	"testing"
)

func TestValidate_rootMissingFormat(t *testing.T) {
	cfg := &Config{Reminders: []Reminder{}}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for missing format")
	}
	if !strings.Contains(err.Error(), "format") {
		t.Errorf("error should mention format: %v", err)
	}
}

func TestValidate_rootInvalidFormat(t *testing.T) {
	cfg := &Config{Format: "12H", Reminders: []Reminder{}}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestValidate_rootMissingReminders(t *testing.T) {
	cfg := &Config{Format: "24h"}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for nil reminders")
	}
}

func TestValidate_validMinimal(t *testing.T) {
	cfg := &Config{Format: "24h", Reminders: []Reminder{}}
	warnings, err := Validate(cfg)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
}

func TestValidate_reminderEmptyName(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "", Enabled: true, Type: "alert", DaysActive: 0, Alert: &Alert{Value: "07:00"}},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestValidate_reminderInvalidType(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "X", Enabled: true, Type: "invalid", DaysActive: 0},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid type")
	}
}

func TestValidate_reminderDaysActiveOutOfRange(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "X", Enabled: true, Type: "alert", DaysActive: 128, Alert: &Alert{Value: "07:00"}},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for daysActive > 127")
	}
}

func TestValidate_alertWithoutAlertSection(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "X", Enabled: true, Type: "alert", DaysActive: 0},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for alert type without alert section")
	}
}

func TestValidate_counterWithoutCounterSection(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "X", Enabled: true, Type: "counter", DaysActive: 0},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for counter type without counter section")
	}
}

func TestValidate_duplicateNames(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "A", Enabled: true, Type: "alert", DaysActive: 0, Alert: &Alert{Value: "07:00"}},
			{Name: "A", Enabled: false, Type: "alert", DaysActive: 0, Alert: &Alert{Value: "08:00"}},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for duplicate names")
	}
	if !strings.Contains(err.Error(), "unique") {
		t.Errorf("error should mention unique: %v", err)
	}
}

func TestValidate_counterInvalidCurrentValue(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{
				Name: "C", Enabled: true, Type: "counter", DaysActive: 0,
				Counter: &Counter{
					Format: "number", CurrentValue: "not-a-number", InitialValue: "0",
					Rules: []CounterRule{{Every: "1s", Increase: "1"}},
				},
			},
		},
	}
	_, err := Validate(cfg)
	if err == nil {
		t.Fatal("expected error for invalid counter currentValue")
	}
}

func TestValidate_validAlert(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{Name: "A", Enabled: true, Type: "alert", DaysActive: 21, Alert: &Alert{Value: "07:00"}},
		},
	}
	_, err := Validate(cfg)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidate_validCounterDuration(t *testing.T) {
	cfg := &Config{
		Format: "24h",
		Reminders: []Reminder{
			{
				Name: "C", Enabled: true, Type: "counter", DaysActive: 31,
				Counter: &Counter{
					Format: "duration", CurrentValue: "22:00", InitialValue: "00:00",
					Rules: []CounterRule{{Every: "1m", Increase: "1m"}},
				},
			},
		},
	}
	_, err := Validate(cfg)
	if err != nil {
		t.Fatalf("Validate: %v", err)
	}
}
