package main

import (
	_ "embed"
	"log"
	"reminder-app/pkg/controller"
	"reminder-app/pkg/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	_ "fyne.io/fyne/v2/driver/software"
)

//go:embed assets/icon.png
var IconPNG []byte

//go:embed assets/default-alarm.mp3
var DefaultAlarm []byte

func main() {
	a := app.New()
	icon := fyne.NewStaticResource("icon", IconPNG)

	windowConfig := ui.DefaultConfig()
	mainWindow := ui.NewMainWindow(a, icon, windowConfig)

	var systemTray *ui.SystemTray
	if desk, ok := a.(desktop.App); ok {
		systemTray = ui.NewSystemTray(desk, mainWindow, icon)
	}

	appController, err := controller.NewAppController(DefaultAlarm, mainWindow, systemTray)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer appController.Shutdown()

	// Hide window at startup
	fyne.Do(func() { mainWindow.Window.Hide() })

	mainWindow.Window.ShowAndRun()
}
