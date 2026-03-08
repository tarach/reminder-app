package config

// Config is the root configuration structure.
type Config struct {
	Format    string     `json:"format"`
	Reminders []Reminder `json:"reminders"`
}

// Reminder represents a single reminder entry.
type Reminder struct {
	Name      string  `json:"name"`
	Enabled   bool    `json:"enabled"`
	Type      string  `json:"type"`
	DaysActive int    `json:"daysActive"`
	Alarm     *Alarm  `json:"alarm,omitempty"`
	Alert     *Alert  `json:"alert,omitempty"`
	Counter   *Counter `json:"counter,omitempty"`
}

// Alarm holds alarm rules and optional sound.
type Alarm struct {
	Rules []AlarmRule `json:"rules,omitempty"`
	Sound *Sound      `json:"sound,omitempty"`
}

// AlarmRule has exactly one of At or Every.
type AlarmRule struct {
	At    string `json:"at,omitempty"`
	Every string `json:"every,omitempty"`
}

// Sound has exactly one of Default or File.
type Sound struct {
	Default *bool  `json:"default,omitempty"`
	File    string `json:"file,omitempty"`
}

// Alert holds the alert time value for type "alert" reminders.
type Alert struct {
	Value string `json:"value"`
}

// Counter holds counter format, values, and rules for type "counter" reminders.
type Counter struct {
	Format       string         `json:"format"`
	CurrentValue string         `json:"currentValue"`
	InitialValue string         `json:"initialValue"`
	DisplayFormat string        `json:"displayFormat,omitempty"`
	Rules        []CounterRule  `json:"rules"`
}

// CounterRule has Every and exactly one of Increase or SetValue.
type CounterRule struct {
	Every     string `json:"every"`
	Increase  string `json:"increase,omitempty"`
	SetValue  string `json:"setValue,omitempty"`
}

// Warning holds a runtime validation warning (e.g. sound file missing); DisableSound tells the app to disable sound for that reminder.
type Warning struct {
	ReminderIndex int
	Message       string
	DisableSound  bool
}
