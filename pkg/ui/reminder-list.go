package ui

import (
	"image/color"
	"reminder-app/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Design palette - used for custom canvas elements
var (
	cardBgColor     = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	cardBorderColor = color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF}
	cardShadowColor = color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x18}
	dayInactiveBg   = color.NRGBA{R: 0xF1, G: 0xF5, B: 0xF9, A: 0xFF}
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
	w.Window.SetIcon(icon)
	w.setupLayout()
	return w
}

// SetFromConfig replaces reminders with those from cfg and refreshes the layout.
func (w *ReminderListWindow) SetFromConfig(cfg *config.Config) {
	w.Reminders = remindersFromConfig(cfg)
	w.setupLayout()
}

var typeToUI = map[string]ReminderType{
	"alert":   ReminderTypeClock,
	"counter": ReminderTypeCounter,
}

func remindersFromConfig(cfg *config.Config) []*Reminder {
	if cfg == nil || len(cfg.Reminders) == 0 {
		return nil
	}
	out := make([]*Reminder, 0, len(cfg.Reminders))
	for i := range cfg.Reminders {
		r := &cfg.Reminders[i]
		uiType := typeToUI[r.Type]
		if uiType == "" {
			uiType = ReminderTypeCounter
		}
		timeVal := ""
		if r.Alert != nil {
			timeVal = r.Alert.Value
		}
		if r.Counter != nil {
			timeVal = r.Counter.CurrentValue
		}
		out = append(out, &Reminder{
			Name:       r.Name,
			Type:       uiType,
			TimeValue:  timeVal,
			ActiveDays: daysActiveToMap(r.DaysActive),
			IsEnabled:  r.Enabled,
		})
	}
	return out
}

// daysActiveToMap converts bitmask (0-127) to map; bits 0..6 = M,T,W,Th,F,Sa,Su.
func daysActiveToMap(daysActive int) map[string]bool {
	daysOrder := []string{"M", "T", "W", "Th", "F", "Sa", "Su"}
	m := make(map[string]bool, 7)
	for i, day := range daysOrder {
		m[day] = (daysActive & (1 << i)) != 0
	}
	return m
}

func (w *ReminderListWindow) setupLayout() {
	reminderList := container.NewVBox()

	// Create a reminder card for each reminder
	for _, reminder := range w.Reminders {
		reminderCard := w.createReminderCard(reminder)
		reminderList.Add(reminderCard)
	}

	// Add a prominent button to create new reminders
	addButton := widget.NewButton("  +  Add Reminder  ", func() {
		// TODO: Open dialog to create new reminder
	})
	addButton.Importance = widget.HighImportance

	scrollContent := container.NewVBox(reminderList)
	scroll := container.NewVScroll(scrollContent)
	scroll.SetMinSize(fyne.NewSize(0, 0))

	content := container.NewBorder(
		nil,
		container.NewPadded(addButton),
		nil,
		nil,
		scroll,
	)

	w.Window.SetContent(content)
	w.Window.Resize(fyne.NewSize(500, 480))
}

func (w *ReminderListWindow) createReminderCard(reminder *Reminder) fyne.CanvasObject {
	// Name label - prominent, readable
	nameLabel := widget.NewLabelWithStyle(reminder.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	nameLabel.Wrapping = fyne.TextWrapOff

	// Time/Counter display - large, clear typography
	timeText := canvas.NewText(reminder.TimeValue, theme.Color(theme.ColorNameForeground))
	timeText.TextSize = 28
	timeText.TextStyle = fyne.TextStyle{Bold: true}
	timeText.Alignment = fyne.TextAlignCenter

	// Type indicator - subtle
	typeText := canvas.NewText(string(reminder.Type), theme.Color(theme.ColorNameDisabled))
	typeText.TextSize = 11
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

	// On/Off toggle switch - wrap in Center so it keeps its compact size (Border stretches right object to full height)
	toggle := NewToggleSwitch(reminder.IsEnabled, func(checked bool) {
		reminder.IsEnabled = checked
	})
	toggleContainer := container.NewCenter(toggle)

	// Layout: Name on top, days below
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
		toggleContainer,
		container.NewCenter(centerSection),
	)

	// Card background - clean white with subtle border and rounded corners
	rect := canvas.NewRectangle(cardBgColor)
	rect.StrokeColor = cardBorderColor
	rect.StrokeWidth = 1
	rect.CornerRadius = 8

	// Card with padding
	cardWithBg := container.NewStack(
		rect,
		container.NewPadded(container.NewPadded(cardContent)),
	)

	// Shadow layer - extends 6px beyond card for outer shadow effect
	shadowRect := canvas.NewRectangle(cardShadowColor)
	shadowRect.CornerRadius = 14 // 8 (card radius) + 6 (shadow padding)
	shadowPadding := float32(6)
	paddedCard := container.New(layout.NewCustomPaddedLayout(0, shadowPadding, 0, shadowPadding), cardWithBg)
	cardWithShadow := container.NewStack(shadowRect, paddedCard)

	// Wrap with margin and padding for spacing between cards
	return container.NewPadded(container.NewPadded(cardWithShadow))
}

func (w *ReminderListWindow) createDayButton(day string, isActive bool) fyne.CanvasObject {
	bgColor := getDayBackgroundColor(isActive)
	circle := canvas.NewCircle(bgColor)

	textColor := theme.Color(theme.ColorNameForeground)
	if isActive {
		textColor = theme.Color(theme.ColorNameForegroundOnPrimary)
	}

	dayText := canvas.NewText(day, textColor)
	dayText.TextSize = 11
	dayText.Alignment = fyne.TextAlignCenter
	dayText.TextStyle = fyne.TextStyle{Bold: true}

	// Calculate text size and add padding for the circle
	textSize := dayText.MinSize()
	padding := float32(6)
	circleSize := fyne.NewSize(
		textSize.Width+padding*2,
		textSize.Height+padding*2,
	)

	// Make it a perfect circle
	maxDim := circleSize.Width
	if circleSize.Height > maxDim {
		maxDim = circleSize.Height
	}
	circleSize = fyne.NewSize(maxDim, maxDim)

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(circleSize)

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
	return dayInactiveBg
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
