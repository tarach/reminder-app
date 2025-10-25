package controller

import (
	"fmt"
	"reminder-app/pkg/audio"
	"reminder-app/pkg/ui"
	"sync"
	"time"

	"fyne.io/fyne/v2"
)

// AppController coordinates between UI and audio player
type AppController struct {
	player       *audio.Player
	mainWindow   *ui.MainWindow
	systemTray   *ui.SystemTray
	stopTicker   chan struct{}
	tickerMu     sync.Mutex
	tickerActive bool
	shutdownOnce sync.Once
}

// NewAppController creates a new application controller
func NewAppController(audioData []byte, mainWindow *ui.MainWindow, systemTray *ui.SystemTray) (*AppController, error) {
	player, err := audio.NewPlayer(audioData)
	if err != nil {
		return nil, err
	}

	c := &AppController{
		player:       player,
		mainWindow:   mainWindow,
		systemTray:   systemTray,
		stopTicker:   nil,
		tickerActive: false,
	}

	c.bindUICallbacks()

	return c, nil
}

func (c *AppController) startUIUpdates() {
	c.tickerMu.Lock()
	defer c.tickerMu.Unlock()

	if c.tickerActive {
		return // Already running
	}

	c.stopTicker = make(chan struct{})
	c.tickerActive = true

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.updateUI()
			case <-c.stopTicker:
				return
			}
		}
	}()
}

func (c *AppController) stopUIUpdates() {
	c.tickerMu.Lock()
	defer c.tickerMu.Unlock()

	if !c.tickerActive {
		return // Already stopped
	}

	close(c.stopTicker)
	c.tickerActive = false
}

func (c *AppController) bindUICallbacks() {
	c.mainWindow.StartBtnFn = func() {
		c.player.Play()
		c.startUIUpdates()
	}
	c.mainWindow.StopBtnFn = func() {
		c.player.Stop()
		c.stopUIUpdates()
	}
	c.mainWindow.VolSlider.OnChanged = c.player.SetVolume

	if c.systemTray != nil {
		c.systemTray.QuitMenuItemFn = c.Shutdown
	}

	// Initialize UI
	pos, total := c.player.GetPosition()
	totalText := c.formatDuration(c.player.GetDuration(total))
	posText := c.formatDuration(c.player.GetDuration(pos))
	fyne.Do(func() {
		c.mainWindow.TimeLabel.SetText(fmt.Sprintf("%s / %s", posText, totalText))
	})
}

func (c *AppController) updateUI() {
	pos, total := c.player.GetPosition()

	frac := 0.0
	if total > 0 {
		frac = float64(pos) / float64(total)
	}

	posText := c.formatDuration(c.player.GetDuration(pos))
	totalText := c.formatDuration(c.player.GetDuration(total))

	fyne.Do(func() {
		c.mainWindow.Progress.SetValue(frac)
		c.mainWindow.TimeLabel.SetText(fmt.Sprintf("%s / %s", posText, totalText))
	})
}

func (c *AppController) formatDuration(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d", int(d.Minutes()), int(d.Seconds())%60)
}

// Shutdown gracefully shuts down the application
func (c *AppController) Shutdown() {
	c.shutdownOnce.Do(func() {
		c.stopUIUpdates()
		c.player.Close()
		c.mainWindow.Application.Quit()
	})
}
