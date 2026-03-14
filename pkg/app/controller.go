package app

import (
	"reminder-app/pkg/config"
	"reminder-app/pkg/modules/remindercreate"
	"reminder-app/pkg/modules/reminderlist"
	"reminder-app/pkg/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

const listWindowWidth = 500
const listWindowHeight = 480

// Controller is the app-level coordinator: owns the main window and screen navigation.
type Controller struct {
	app          fyne.App
	window       fyne.Window
	icon         fyne.Resource
	sampleConfig []byte
	onQuit       func()

	listView       *reminderlist.View
	listController *reminderlist.Controller

	createView       *remindercreate.View
	createController *remindercreate.Controller

	currentConfig *config.Config
	writable      bool
}

// NewController creates the app controller. Call Start() to run.
func NewController(app fyne.App, icon fyne.Resource, sampleConfig []byte, onQuit func()) *Controller {
	c := &Controller{
		app:          app,
		icon:         icon,
		sampleConfig: sampleConfig,
		onQuit:       onQuit,
	}

	c.listView = reminderlist.NewView()
	c.listController = reminderlist.NewController(c.listView, reminderlist.Callbacks{
		OnAdd: func() { c.showCreate() },
	})

	c.createView = remindercreate.NewView()
	c.createController = remindercreate.NewController(c.createView, remindercreate.Callbacks{
		OnSave:   c.onSaveCreate,
		OnCancel: func() { c.showList() },
	})

	return c
}

// Start initializes the main window, runs the config flow, and shows the reminder list when config is ready.
func (c *Controller) Start() {
	c.window = c.app.NewWindow("Reminders")
	c.window.SetIcon(c.icon)
	c.window.SetContent(container.NewVBox()) // placeholder until config is loaded
	c.window.Show()

	ui.StartConfigFlow(c.window, c.sampleConfig, c.onConfigReady, c.onQuit)
}

// Run blocks on the main window event loop. Call after Start().
func (c *Controller) Run() {
	c.window.ShowAndRun()
}

func (c *Controller) onConfigReady(cfg *config.Config, writable bool) {
	c.currentConfig = cfg
	c.writable = writable
	c.listController.ApplyConfig(cfg, c.showList)
}

func (c *Controller) showList() {
	// Use fyne.Do so content updates are safe when called from button callbacks.
	fyne.Do(func() {
		c.window.SetContent(c.listView.Content())
		c.window.Resize(fyne.NewSize(listWindowWidth, listWindowHeight))
	})
}

func (c *Controller) showCreate() {
	// Defer content switch so we don't call SetContent from inside the button's
	// callback (can deadlock or prevent the UI from updating).
	fyne.Do(func() {
		c.window.SetContent(c.createController.Content())
		c.window.Resize(fyne.NewSize(500, 520))
	})
}

func (c *Controller) onSaveCreate(input remindercreate.Input) {
	rem := inputToReminder(input)
	if c.currentConfig == nil {
		c.currentConfig = &config.Config{Format: "24h", Reminders: []config.Reminder{}}
	}
	c.currentConfig.Reminders = append(c.currentConfig.Reminders, rem)

	if c.writable {
		result := config.Lookup()
		if result.OK {
			_ = config.Save(result.Path, c.currentConfig)
		}
	}

	c.listController.ApplyConfig(c.currentConfig, c.showList)
}

// inputToReminder converts create form input to config.Reminder.
func inputToReminder(in remindercreate.Input) config.Reminder {
	daysActive := 0
	for i := 0; i < 7; i++ {
		if in.ActiveDays[i] {
			daysActive |= 1 << i
		}
	}

	rem := config.Reminder{
		Name:       in.Name,
		Enabled:    in.Enabled,
		Type:       in.Type,
		DaysActive: daysActive,
	}

	if in.Type == "alert" {
		val := in.AlertValue
		if val == "" {
			val = "08:00"
		}
		rem.Alert = &config.Alert{Value: val}
		rem.Alarm = &config.Alarm{
			Rules: []config.AlarmRule{{At: "8h"}},
			Sound: &config.Sound{Default: boolPtr(true)},
		}
	} else {
		initVal := in.CounterInitialValue
		if initVal == "" {
			initVal = "00:00"
		}
		dispFmt := in.CounterDisplayFormat
		if dispFmt == "" {
			dispFmt = "mm:ss"
		}
		counterFormat := "duration"
		if dispFmt == "number" {
			if initVal == "" || initVal == "00:00" {
				initVal = "0"
			}
			counterFormat = "number"
		}
		rem.Counter = &config.Counter{
			Format:         counterFormat,
			CurrentValue:   initVal,
			InitialValue:   initVal,
			DisplayFormat:  dispFmt,
			Rules:         defaultCounterRules(counterFormat, dispFmt),
		}
		rem.Alarm = &config.Alarm{
			Rules: []config.AlarmRule{{Every: "1h"}},
			Sound: &config.Sound{Default: boolPtr(true)},
		}
	}

	return rem
}

// defaultCounterRules returns mandatory default rules for a new counter based on format and displayFormat.
func defaultCounterRules(format, displayFormat string) []config.CounterRule {
	switch format {
	case "number":
		return []config.CounterRule{
			{Every: "1s", Increase: "1"},
			{Every: "14s", SetValue: "0"},
		}
	case "duration":
		switch displayFormat {
		case "hh:mm":
			return []config.CounterRule{
				{Every: "1m", Increase: "1m"},
			}
		default: // "mm:ss" or any other
			return []config.CounterRule{
				{Every: "1s", Increase: "1s"},
			}
		}
	default:
		return []config.CounterRule{
			{Every: "1s", Increase: "1s"},
		}
	}
}

func boolPtr(b bool) *bool {
	return &b
}
