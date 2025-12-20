package controller

import (
	"fmt"
	"reminder-app/pkg/audio"
	"reminder-app/pkg/ui"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// PlayerController coordinates between UI and audio player
type PlayerController struct {
	player       *audio.Player
	mainWindow   *ui.PlayerWindow
	systemTray   *ui.SystemTray
	stopTicker   chan struct{}
	tickerMu     sync.Mutex
	tickerActive bool
	shutdownOnce sync.Once
}

// NewPlayerController creates a new application controller
func NewPlayerController(audioData []byte, mainWindow *ui.PlayerWindow, systemTray *ui.SystemTray) (*PlayerController, error) {
	player, err := audio.NewPlayer(audioData)
	if err != nil {
		return nil, err
	}

	c := &PlayerController{
		player:       player,
		mainWindow:   mainWindow,
		systemTray:   systemTray,
		stopTicker:   nil,
		tickerActive: false,
	}

	c.bindUICallbacks()

	return c, nil
}

func (pc *PlayerController) startUIUpdates() {
	pc.tickerMu.Lock()
	defer pc.tickerMu.Unlock()

	if pc.tickerActive {
		return // Already running
	}

	pc.stopTicker = make(chan struct{})
	pc.tickerActive = true

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pc.updateUI()
			case <-pc.stopTicker:
				return
			}
		}
	}()
}

func (pc *PlayerController) stopUIUpdates() {
	pc.tickerMu.Lock()
	defer pc.tickerMu.Unlock()

	if !pc.tickerActive {
		return // Already stopped
	}

	close(pc.stopTicker)
	pc.tickerActive = false
}

func (pc *PlayerController) bindUICallbacks() {
	pc.mainWindow.StartBtnFn = func() {
		pc.player.Play()
		pc.startUIUpdates()
	}
	pc.mainWindow.StopBtnFn = func() {
		pc.player.Stop()
		pc.stopUIUpdates()
	}
	pc.mainWindow.VolSlider.OnChanged = pc.player.SetVolume

	if pc.systemTray != nil {
		pc.systemTray.QuitMenuItemFn = pc.Shutdown
	}

	// Initialize UI
	pos, total := pc.player.GetPosition()
	totalText := pc.formatDuration(pc.player.GetDuration(total))
	posText := pc.formatDuration(pc.player.GetDuration(pos))
	fyne.Do(func() {
		pc.mainWindow.TimeLabel.SetText(fmt.Sprintf("%s / %s", posText, totalText))
	})
}

func (pc *PlayerController) updateUI() {
	pos, total := pc.player.GetPosition()

	frac := 0.0
	if total > 0 {
		frac = float64(pos) / float64(total)
	}

	posText := pc.formatDuration(pc.player.GetDuration(pos))
	totalText := pc.formatDuration(pc.player.GetDuration(total))

	fyne.Do(func() {
		pc.mainWindow.Progress.SetValue(frac)
		pc.mainWindow.TimeLabel.SetText(fmt.Sprintf("%s / %s", posText, totalText))
	})
}

func (pc *PlayerController) formatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(d.Minutes()), int(d.Seconds())%60)
}

// Shutdown gracefully shuts down the application
func (pc *PlayerController) Shutdown() {
	pc.shutdownOnce.Do(func() {
		pc.stopUIUpdates()
		pc.player.Close()
		pc.mainWindow.Application.Quit()
	})
}
