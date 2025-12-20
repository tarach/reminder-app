package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ReminderType string

const (
	ReminderTypeCounter ReminderType = "counter"
	ReminderTypeClock   ReminderType = "clock"
)

type Reminder struct {
	Name       string
	Type       ReminderType
	TimeValue  string          // "00:00" for counter or clock display
	ActiveDays map[string]bool // "M", "T", "W", "Th", "F", "Sa", "Su"
	IsEnabled  bool
}

type ReminderListWindow struct {
	Window      fyne.Window
	Application fyne.App
	Icon        *fyne.StaticResource
	Reminders   []*Reminder
}

func NewReminderListWindow(app fyne.App, icon *fyne.StaticResource) *ReminderListWindow {
	w := &ReminderListWindow{
		Window:      app.NewWindow("Reminders"),
		Application: app,
		Icon:        icon,
		Reminders:   []*Reminder{},
	}

	// Add some sample reminders
	w.Reminders = append(w.Reminders,
		&Reminder{
			Name:       "Morning Workout",
			Type:       ReminderTypeClock,
			TimeValue:  "07:00",
			ActiveDays: map[string]bool{"M": true, "W": true, "F": true},
			IsEnabled:  true,
		},
		&Reminder{
			Name:       "Study Session",
			Type:       ReminderTypeCounter,
			TimeValue:  "25:00",
			ActiveDays: map[string]bool{"M": true, "T": true, "W": true, "Th": true, "F": true, "Sa": false, "Su": false},
			IsEnabled:  false,
		},
		&Reminder{
			Name:       "Water Break",
			Type:       ReminderTypeCounter,
			TimeValue:  "01:30",
			ActiveDays: map[string]bool{"M": true, "T": true, "W": true, "Th": true, "F": true, "Sa": true, "Su": true},
			IsEnabled:  true,
		},
	)
	w.Window.SetIcon(icon)
	w.setupLayout()
	return w
}

func (w *ReminderListWindow) setupLayout() {
	reminderList := container.NewVBox()

	// Create a reminder card for each reminder
	for _, reminder := range w.Reminders {
		reminderCard := w.createReminderCard(reminder)
		reminderList.Add(reminderCard)
	}

	// Add a button to create new reminders
	addButton := widget.NewButton("+ Add Reminder", func() {
		// TODO: Open dialog to create new reminder
	})

	content := container.NewBorder(
		nil,
		addButton,
		nil,
		nil,
		container.NewVScroll(reminderList),
	)

	w.Window.SetContent(content)
	w.Window.Resize(fyne.NewSize(500, 400))
}

func (w *ReminderListWindow) createReminderCard(reminder *Reminder) fyne.CanvasObject {
	// Name label
	nameLabel := widget.NewLabelWithStyle(reminder.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Time/Counter display in the center using canvas.Text for custom size
	timeText := canvas.NewText(reminder.TimeValue, theme.Color(theme.ColorNameForeground))
	timeText.TextSize = 32
	timeText.TextStyle = fyne.TextStyle{Bold: true}
	timeText.Alignment = fyne.TextAlignCenter

	// Type indicator using canvas.Text for custom size
	typeText := canvas.NewText(string(reminder.Type), theme.Color(theme.ColorNameDisabled))
	typeText.TextSize = 10
	typeText.Alignment = fyne.TextAlignCenter

	// Center section with time and type
	centerSection := container.NewVBox(
		timeText,
		typeText,
	)

	// Days of the week
	daysOrder := []string{"M", "T", "W", "Th", "F", "Sa", "Su"}
	dayButtons := container.NewHBox()

	for _, day := range daysOrder {
		isActive := reminder.ActiveDays[day]
		dayBtn := w.createDayButton(day, isActive)
		dayButtons.Add(dayBtn)
	}

	// On/Off switch
	toggle := widget.NewCheck("", func(checked bool) {
		reminder.IsEnabled = checked
	})
	toggle.Checked = reminder.IsEnabled

	// Layout: Name on top, center section, days, and toggle
	// Wrap both nameLabel and dayButtons in containers without extra nesting
	leftSection := container.NewBorder(
		nameLabel,
		nil,
		nil,
		nil,
		dayButtons,
	)

	cardContent := container.NewBorder(
		nil,
		nil,
		leftSection,
		toggle,
		container.NewCenter(centerSection),
	)

	// Create a rectangle background
	rect := canvas.NewRectangle(color.RGBA{R: 240, G: 240, B: 245, A: 255})
	rect.StrokeColor = color.RGBA{R: 200, G: 200, B: 210, A: 255}
	rect.StrokeWidth = 1

	// Use padding container with minimum size for the card
	cardWithBg := container.NewStack(
		rect,
		container.NewPadded(cardContent),
	)

	// Wrap in a container with minimum size
	sizedCard := container.NewPadded(cardWithBg)

	return sizedCard
}

func (w *ReminderListWindow) createDayButton(day string, isActive bool) fyne.CanvasObject {
	var bgColor = getDayBackgroundColor(isActive)
	circle := canvas.NewCircle(bgColor)

	var textColor color.Color
	textColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}

	dayText := canvas.NewText(day, textColor)
	dayText.TextSize = 10
	dayText.Alignment = fyne.TextAlignCenter
	dayText.TextStyle = fyne.TextStyle{Bold: true}

	// Calculate text size and add padding for the circle
	textSize := dayText.MinSize()
	padding := float32(4) // Adjust this value to control how much bigger the circle is
	circleSize := fyne.NewSize(
		textSize.Width+padding*2,
		textSize.Height+padding*2,
	)

	// Make it a perfect circle by using the larger dimension
	if circleSize.Width > circleSize.Height {
		circleSize.Height = circleSize.Width
	} else {
		circleSize.Width = circleSize.Height
	}

	// Create a spacer to define the minimum size based on text + padding
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(circleSize)

	// Stack: spacer (invisible, defines size), circle (fills the space), text (centered)
	return container.NewStack(
		spacer,
		circle,
		container.NewCenter(dayText),
	)
}

func getDayBackgroundColor(isActive bool) color.Color {
	if isActive {
		return theme.Color(theme.ColorNamePrimary)
	}
	return theme.Color(theme.ColorNameDisabled)
}

func (w *ReminderListWindow) Show() {
	w.Window.Show()
}

// fixedSizeContainer is a container that enforces a fixed size
type fixedSizeContainer struct {
	obj  fyne.CanvasObject
	size fyne.Size
}

func (f *fixedSizeContainer) MinSize() fyne.Size {
	return f.size
}

func (f *fixedSizeContainer) Resize(size fyne.Size) {
	f.obj.Resize(f.size)
}

func (f *fixedSizeContainer) Move(pos fyne.Position) {
	f.obj.Move(pos)
}

func (f *fixedSizeContainer) Position() fyne.Position {
	return f.obj.Position()
}

func (f *fixedSizeContainer) Size() fyne.Size {
	return f.size
}

func (f *fixedSizeContainer) Visible() bool {
	return f.obj.Visible()
}

func (f *fixedSizeContainer) Show() {
	f.obj.Show()
}

func (f *fixedSizeContainer) Hide() {
	f.obj.Hide()
}

func (f *fixedSizeContainer) Refresh() {
	f.obj.Refresh()
}
