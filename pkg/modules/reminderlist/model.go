package reminderlist

// Item is a single reminder row for the list UI.
type Item struct {
	Name        string
	Kind        string   // "counter" or "clock"
	DisplayText string   // time or counter value to show
	ActiveDays  [7]bool  // M, T, W, Th, F, Sa, Su
	Enabled     bool
}

// Model is the UI model for the reminder list screen.
type Model struct {
	Items []Item
}
