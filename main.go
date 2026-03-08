package main

import (
	_ "embed"
	"reminder-app/pkg/config"
	"reminder-app/pkg/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/driver/software"
)

//go:embed assets/icon.png
var IconPNG []byte

//go:embed assets/default-alarm.mp3
var DefaultAlarm []byte

//go:embed config/sample.json
var SampleConfig []byte

func main() {
	a := app.New()
	a.Settings().SetTheme(&ui.AppTheme{})
	icon := fyne.NewStaticResource("icon", IconPNG)

	reminderListWindow := ui.NewReminderListWindow(a, icon)
	reminderListWindow.Window.Show()
	ui.StartConfigFlow(reminderListWindow.Window, SampleConfig, func(cfg *config.Config, _ bool) {
		reminderListWindow.SetFromConfig(cfg)
	}, a.Quit)
	reminderListWindow.Window.ShowAndRun()
	/*
		windowConfig := ui.DefaultConfig()
		mainWindow := ui.NewPlayerWindow(a, icon, windowConfig)

		var systemTray *ui.SystemTray
		if desk, ok := a.(desktop.App); ok {
			systemTray = ui.NewSystemTray(desk, mainWindow, icon)
		}

		appController, err := controller.NewPlayerController(DefaultAlarm, mainWindow, systemTray)
		if err != nil {
			log.Fatalf("Failed to initialize application: %v", err)
		}
		defer appController.Shutdown()

		// Hide window at startup
		fyne.Do(func() { mainWindow.Window.Hide() })

		mainWindow.Window.ShowAndRun()
	*/
}
