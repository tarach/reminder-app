package reminderlist

import "reminder-app/pkg/config"

var typeToKind = map[string]string{
	"alert":   "clock",
	"counter": "counter",
}

// daysActiveToArray converts bitmask (0–127) to a fixed array; bits 0..6 = M, T, W, Th, F, Sa, Su.
func daysActiveToArray(daysActive int) [7]bool {
	var out [7]bool
	for i := 0; i < 7; i++ {
		out[i] = (daysActive & (1 << i)) != 0
	}
	return out
}

// ModelFromConfig converts config into the UI model for the reminder list.
func ModelFromConfig(cfg *config.Config) Model {
	if cfg == nil || len(cfg.Reminders) == 0 {
		return Model{Items: nil}
	}
	items := make([]Item, 0, len(cfg.Reminders))
	for i := range cfg.Reminders {
		r := &cfg.Reminders[i]
		kind := typeToKind[r.Type]
		if kind == "" {
			kind = "counter"
		}
		displayText := ""
		if r.Alert != nil {
			displayText = r.Alert.Value
		}
		if r.Counter != nil {
			displayText = r.Counter.CurrentValue
		}
		items = append(items, Item{
			Name:        r.Name,
			Kind:        kind,
			DisplayText: displayText,
			ActiveDays:  daysActiveToArray(r.DaysActive),
			Enabled:     r.Enabled,
		})
	}
	return Model{Items: items}
}
