package config

import (
	"testing"
)

func TestParseJSON_validRootObject(t *testing.T) {
	data := []byte(`{"format":"24h","reminders":[]}`)
	cfg, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON: %v", err)
	}
	if cfg.Format != "24h" {
		t.Errorf("format: got %q", cfg.Format)
	}
	if cfg.Reminders == nil || len(cfg.Reminders) != 0 {
		t.Errorf("reminders: got %v", cfg.Reminders)
	}
}

func TestParseJSON_invalidJSON(t *testing.T) {
	data := []byte(`{invalid}`)
	_, err := ParseJSON(data)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseJSON_rootNotObject(t *testing.T) {
	data := []byte(`["format","24h"]`)
	_, err := ParseJSON(data)
	if err == nil {
		t.Fatal("expected error when root is not object")
	}
}

func TestLoadFromBytes_sampleStructure(t *testing.T) {
	sample := []byte(`{
		"format": "24h",
		"reminders": [
			{"name":"Test","enabled":true,"type":"alert","daysActive":21,"alarm":{"rules":[{"at":"7h"}],"sound":{"default":true}},"alert":{"value":"07:00"}}
		]
	}`)
	cfg, err := LoadFromBytes(sample)
	if err != nil {
		t.Fatalf("LoadFromBytes: %v", err)
	}
	if cfg.Format != "24h" {
		t.Errorf("format: got %q", cfg.Format)
	}
	if len(cfg.Reminders) != 1 {
		t.Fatalf("reminders: got len %d", len(cfg.Reminders))
	}
	r := cfg.Reminders[0]
	if r.Name != "Test" || r.Type != "alert" || r.DaysActive != 21 {
		t.Errorf("reminder: %+v", r)
	}
	if r.Alert == nil || r.Alert.Value != "07:00" {
		t.Errorf("alert: %+v", r.Alert)
	}
}
