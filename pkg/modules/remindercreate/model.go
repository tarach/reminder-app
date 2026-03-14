package remindercreate

// Input is the form input for creating a reminder.
type Input struct {
	Name                 string
	Type                 string   // "alert" or "counter"
	Enabled              bool
	ActiveDays           [7]bool  // M, T, W, Th, F, Sa, Su
	AlertValue           string   // for type "alert", e.g. "07:00"
	CounterInitialValue  string   // for type "counter", e.g. "00:00"
	CounterDisplayFormat string   // e.g. "mm:ss" or "hh:mm"
}
