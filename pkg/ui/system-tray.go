package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// SystemTray manages the system tray icon and menu
type SystemTray struct {
	QuitMenuItemFn func()
	mainWindow     *PlayerWindow
	desk           desktop.App
	icon           *fyne.StaticResource
}

// NewSystemTray creates and initializes the system tray
func NewSystemTray(desk desktop.App, mainWindow *PlayerWindow, icon *fyne.StaticResource) *SystemTray {
	t := &SystemTray{
		mainWindow: mainWindow,
		desk:       desk,
		icon:       icon,
	}
	t.initialize()
	return t
}

func (t *SystemTray) initialize() {
	fyne.Do(func() {
		t.desk.SetSystemTrayIcon(t.icon)

		menu := fyne.NewMenu("Reminder",
			fyne.NewMenuItem("Show", t.showWindow),
			fyne.NewMenuItem("Quit", t.quit),
		)
		t.desk.SetSystemTrayMenu(menu)

		t.mainWindow.Window.SetCloseIntercept(func() {
			t.mainWindow.Window.Hide()
		})
	})
}

func (t *SystemTray) showWindow() {
	t.mainWindow.Window.Show()
	t.mainWindow.Window.RequestFocus()
}

func (t *SystemTray) quit() {
	if t.QuitMenuItemFn != nil {
		t.QuitMenuItemFn()
	}
}
