package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// PlayerWindow represents the main application window
type PlayerWindow struct {
	Window      fyne.Window
	Application fyne.App
	Icon        *fyne.StaticResource

	// UI Components
	VolLabel  *widget.Label
	VolSlider *widget.Slider
	Progress  *widget.ProgressBar
	TimeLabel *widget.Label

	// Callbacks
	StartBtnFn func()
	StopBtnFn  func()
}

// Config holds window configuration
type Config struct {
	Width         float32
	Height        float32
	DefaultVolume float64
}

// DefaultConfig returns the default window configuration
func DefaultConfig() Config {
	return Config{
		Width:         420,
		Height:        200,
		DefaultVolume: 100,
	}
}

// NewPlayerWindow creates and initializes the main application window
func NewPlayerWindow(app fyne.App, icon *fyne.StaticResource, config Config) *PlayerWindow {
	w := &PlayerWindow{
		Window:      app.NewWindow("Player"),
		Application: app,
		Icon:        icon,
		VolLabel:    widget.NewLabel("Volume"),
		VolSlider:   widget.NewSlider(0, 100),
		Progress:    widget.NewProgressBar(),
		TimeLabel:   widget.NewLabel("00:00 / 00:00"),
	}

	w.configureComponents(config)
	w.setupLayout()
	w.Window.SetIcon(icon)
	w.Window.Resize(fyne.NewSize(config.Width, config.Height))

	return w
}

func (pw *PlayerWindow) configureComponents(config Config) {
	pw.VolSlider.SetValue(config.DefaultVolume)
	pw.VolSlider.Step = 1

	pw.Progress.Min = 0
	pw.Progress.Max = 1
	pw.Progress.SetValue(0)
}

func (pw *PlayerWindow) setupLayout() {
	startBtn := widget.NewButton("Start", func() {
		if pw.StartBtnFn != nil {
			pw.StartBtnFn()
		}
	})
	stopBtn := widget.NewButton("Stop", func() {
		if pw.StopBtnFn != nil {
			pw.StopBtnFn()
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Reminder Alarm"),
		container.NewGridWithColumns(2, startBtn, stopBtn),
		pw.VolLabel,
		pw.VolSlider,
		pw.Progress,
		pw.TimeLabel,
	)

	pw.Window.SetContent(content)
}
