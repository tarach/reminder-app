package ui

import (
	"reminder-app/pkg/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ConfigReady is called when config is loaded and the app can start with it. Writable is false when config path is read-only.
type ConfigReady func(cfg *config.Config, writable bool)

// ConfigQuit is called when the user must quit (e.g. config not readable).
type ConfigQuit func()

// StartConfigFlow runs config lookup, load, validation and shows the appropriate dialogs.
// It calls onReady when config is ready (with writable=false when path is not writable or user chose not to save).
// It calls onQuit when the user must exit (config not readable).
// sampleJSON is the embedded sample config (e.g. from main).
func StartConfigFlow(parent fyne.Window, sampleJSON []byte, onReady ConfigReady, onQuit ConfigQuit) {
	result := config.Lookup()
	runFlow(result, parent, sampleJSON, onReady, onQuit)
}

func runFlow(result config.LookupResult, parent fyne.Window, sampleJSON []byte, onReady ConfigReady, onQuit ConfigQuit) {
	showNotFound := func() {
		dialog.ShowConfirm("Configuration", "No configuration file found. Load sample list of reminders?", func(load bool) {
			handleNotFoundChoice(load, parent, sampleJSON, onReady, onQuit)
		}, parent)
	}

	handleFound := func(path string) {
		readable := config.IsReadable(path)
		if !readable {
			d := dialog.NewError(configErrNotReadable, parent)
			d.SetOnClosed(func() { fyne.Do(onQuit) })
			d.Show()
			return
		}
		writable := config.IsWritable(path)
		if !writable {
			dialog.ShowConfirm("Configuration", "Configuration file is not writable. All changes will be lost. Continue anyway?", func(continueAnyway bool) {
				if !continueAnyway {
					fyne.Do(onQuit)
					return
				}
				loadAndValidate(path, parent, sampleJSON, onReady, writable)
			}, parent)
			return
		}
		loadAndValidate(path, parent, sampleJSON, onReady, writable)
	}

	flowHandlers := map[bool]func(){
		true:  func() { handleFound(result.Path) },
		false: showNotFound,
	}
	flowHandlers[result.OK]()
}

var configErrNotReadable = &configError{msg: "Configuration file is not readable."}

type configError struct{ msg string }

func (e *configError) Error() string { return e.msg }

func handleNotFoundChoice(load bool, parent fyne.Window, sampleJSON []byte, onReady ConfigReady, onQuit ConfigQuit) {
	if !load {
		askTimeFormatThenStart(parent, nil, onReady, false)
		return
	}
	cfg, err := config.LoadFromBytes(sampleJSON)
	if err != nil {
		dialog.ShowError(err, parent)
		askTimeFormatThenStart(parent, nil, onReady, false)
		return
	}
	askTimeFormatThenSaveAndStart(parent, cfg, sampleJSON, onReady)
}

func askTimeFormatThenSaveAndStart(parent fyne.Window, cfg *config.Config, sampleJSON []byte, onReady ConfigReady) {
	showTimeFormatDialog(parent, func(format string) {
		cfg.Format = format
		homePath := config.HomeConfigPath()
		_ = config.Save(homePath, cfg)
		onReady(cfg, true)
	})
}

func askTimeFormatThenStart(parent fyne.Window, cfg *config.Config, onReady ConfigReady, writable bool) {
	showTimeFormatDialog(parent, func(format string) {
		if cfg == nil {
			cfg = &config.Config{Format: format, Reminders: []config.Reminder{}}
		} else {
			cfg.Format = format
		}
		onReady(cfg, writable)
	})
}

func showTimeFormatDialog(parent fyne.Window, onChoice func(format string)) {
	content := widget.NewLabel("Choose time format for reminder display:")
	confirm := "12 hour"
	dismiss := "24 hour"
	dialog.ShowCustomConfirm("Time format", confirm, dismiss, content, func(use12h bool) {
		format := "24h"
		if use12h {
			format = "12h"
		}
		onChoice(format)
	}, parent)
}

func loadAndValidate(path string, parent fyne.Window, sampleJSON []byte, onReady ConfigReady, writable bool) {
	cfg, err := config.Load(path)
	if err != nil {
		showInvalidAndOfferSample(parent, err.Error(), sampleJSON, path, onReady, writable)
		return
	}
	warnings, err := config.Validate(cfg)
	if err != nil {
		showInvalidAndOfferSample(parent, err.Error(), sampleJSON, path, onReady, writable)
		return
	}
	showRuntimeWarnings(parent, warnings)
	onReady(cfg, writable)
}

func showInvalidAndOfferSample(parent fyne.Window, errMsg string, sampleJSON []byte, path string, onReady ConfigReady, writable bool) {
	content := widget.NewLabel("Invalid configuration: " + errMsg + "\n\nLoad sample configuration instead? If you decline, you can start the app but changes will not be saved.")
	dialog.ShowCustomConfirm("Configuration error", "Load sample", "No, start without saving", content, func(load bool) {
		handleInvalidChoice(load, parent, errMsg, sampleJSON, path, onReady, writable)
	}, parent)
}

func handleInvalidChoice(load bool, parent fyne.Window, _ string, sampleJSON []byte, path string, onReady ConfigReady, writable bool) {
	if load {
		cfg, err := config.LoadFromBytes(sampleJSON)
		if err != nil {
			askTimeFormatThenStart(parent, nil, onReady, false)
			return
		}
		_, err = config.Validate(cfg)
		if err != nil {
			askTimeFormatThenStart(parent, nil, onReady, false)
			return
		}
		cfg.Format = "24h"
		showTimeFormatDialog(parent, func(format string) {
			cfg.Format = format
			if writable {
				_ = config.Save(path, cfg)
			}
			onReady(cfg, writable)
		})
		return
	}
	askTimeFormatThenStart(parent, nil, onReady, false)
}

func showRuntimeWarnings(parent fyne.Window, warnings []config.Warning) {
	for _, w := range warnings {
		msg := w.Message
		if w.DisableSound {
			msg += " Sound will be disabled for \"" + w.ReminderName + "\"."
		}
		ShowInformation("Configuration warning", msg, parent)
	}
}
