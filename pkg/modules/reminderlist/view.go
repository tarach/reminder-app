package reminderlist

import (
	"image/color"

	"reminder-app/pkg/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	cardBgColor     = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	cardBorderColor = color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF}
	cardShadowColor = color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x18}
	dayInactiveBg   = color.NRGBA{R: 0xF1, G: 0xF5, B: 0xF9, A: 0xFF}
)

var daysOrder = [7]string{"M", "T", "W", "Th", "F", "Sa", "Su"}

// View builds the reminder list UI. It does not own a window.
type View struct {
	content fyne.CanvasObject
	onAdd   func()
}

// NewView creates a new list view. Call Render to build content; use Content() to get the root object.
func NewView() *View {
	v := &View{}
	v.content = container.NewVBox() // default until first Render
	return v
}

// SetOnAdd sets the callback for the Add Reminder button.
func (v *View) SetOnAdd(f func()) {
	v.onAdd = f
}

// Content returns the root canvas object for the list screen. Safe to call after Render.
func (v *View) Content() fyne.CanvasObject {
	if v.content == nil {
		return container.NewVBox()
	}
	return v.content
}

// Render updates the list content to display the given model (full rerender).
// onAdd is the callback for the Add button; the button captures this so it must be non-nil when the button is clicked.
func (v *View) Render(model Model, onAdd func()) {
	if onAdd != nil {
		v.onAdd = onAdd
	}
	reminderList := container.NewVBox()
	for i := range model.Items {
		item := &model.Items[i]
		reminderList.Add(v.cardForItem(item))
	}

	addBtnCallback := onAdd
	if addBtnCallback == nil {
		addBtnCallback = v.onAdd
	}
	addButton := widget.NewButton("  +  Add Reminder  ", func() {
		if addBtnCallback != nil {
			addBtnCallback()
		}
	})
	addButton.Importance = widget.HighImportance

	scrollContent := container.NewVBox(reminderList)
	scroll := container.NewVScroll(scrollContent)
	scroll.SetMinSize(fyne.NewSize(0, 0))

	v.content = container.NewBorder(
		nil,
		container.NewPadded(addButton),
		nil,
		nil,
		scroll,
	)
}

func (v *View) cardForItem(item *Item) fyne.CanvasObject {
	nameLabel := widget.NewLabelWithStyle(item.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	nameLabel.Wrapping = fyne.TextWrapOff

	timeText := canvas.NewText(item.DisplayText, theme.Color(theme.ColorNameForeground))
	timeText.TextSize = 28
	timeText.TextStyle = fyne.TextStyle{Bold: true}
	timeText.Alignment = fyne.TextAlignCenter

	typeText := canvas.NewText(item.Kind, theme.Color(theme.ColorNameDisabled))
	typeText.TextSize = 11
	typeText.Alignment = fyne.TextAlignCenter

	centerSection := container.NewVBox(timeText, typeText)

	dayButtons := container.NewHBox()
	for i, day := range daysOrder {
		dayButtons.Add(v.dayButton(day, item.ActiveDays[i]))
	}

	toggle := ui.NewToggleSwitch(item.Enabled, func(bool) {})
	toggleContainer := container.NewCenter(toggle)

	leftSection := container.NewBorder(nameLabel, nil, nil, nil, dayButtons)

	cardContent := container.NewBorder(
		nil,
		nil,
		leftSection,
		toggleContainer,
		container.NewCenter(centerSection),
	)

	rect := canvas.NewRectangle(cardBgColor)
	rect.StrokeColor = cardBorderColor
	rect.StrokeWidth = 1
	rect.CornerRadius = 8

	cardWithBg := container.NewStack(
		rect,
		container.NewPadded(container.NewPadded(cardContent)),
	)

	shadowRect := canvas.NewRectangle(cardShadowColor)
	shadowRect.CornerRadius = 14
	shadowPadding := float32(6)
	paddedCard := container.New(layout.NewCustomPaddedLayout(0, shadowPadding, 0, shadowPadding), cardWithBg)
	cardWithShadow := container.NewStack(shadowRect, paddedCard)

	return container.NewPadded(container.NewPadded(cardWithShadow))
}

func (v *View) dayButton(day string, isActive bool) fyne.CanvasObject {
	bgColor := dayBackgroundColor(isActive)
	circle := canvas.NewCircle(bgColor)

	textColor := theme.Color(theme.ColorNameForeground)
	if isActive {
		textColor = theme.Color(theme.ColorNameForegroundOnPrimary)
	}

	dayText := canvas.NewText(day, textColor)
	dayText.TextSize = 11
	dayText.Alignment = fyne.TextAlignCenter
	dayText.TextStyle = fyne.TextStyle{Bold: true}

	textSize := dayText.MinSize()
	padding := float32(6)
	circleSize := fyne.NewSize(textSize.Width+padding*2, textSize.Height+padding*2)
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

func dayBackgroundColor(isActive bool) color.Color {
	if isActive {
		return theme.Color(theme.ColorNamePrimary)
	}
	return dayInactiveBg
}
