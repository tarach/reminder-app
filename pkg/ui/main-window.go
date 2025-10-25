package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// MainWindow represents the main application window
type MainWindow struct {
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

// DefaultConfig returns default window configuration
func DefaultConfig() Config {
	return Config{
		Width:         420,
		Height:        200,
		DefaultVolume: 100,
	}
}

// NewMainWindow creates and initializes the main application window
func NewMainWindow(app fyne.App, icon *fyne.StaticResource, config Config) *MainWindow {
	w := &MainWindow{
		Window:      app.NewWindow("Reminder"),
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

func (w *MainWindow) configureComponents(config Config) {
	w.VolSlider.SetValue(config.DefaultVolume)
	w.VolSlider.Step = 1

	w.Progress.Min = 0
	w.Progress.Max = 1
	w.Progress.SetValue(0)
}

func (w *MainWindow) setupLayout() {
	startBtn := widget.NewButton("Start", func() {
		if w.StartBtnFn != nil {
			w.StartBtnFn()
		}
	})
	stopBtn := widget.NewButton("Stop", func() {
		if w.StopBtnFn != nil {
			w.StopBtnFn()
		}
	})

	content := container.NewVBox(
		widget.NewLabel("Reminder Alarm"),
		container.NewGridWithColumns(2, startBtn, stopBtn),
		w.VolLabel,
		w.VolSlider,
		w.Progress,
		w.TimeLabel,
	)

	w.Window.SetContent(content)
}
