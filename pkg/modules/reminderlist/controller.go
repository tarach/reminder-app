package reminderlist

import (
	"reminder-app/pkg/config"

	"fyne.io/fyne/v2"
)

// Callbacks are invoked by the reminder list screen (e.g. when Add is clicked).
type Callbacks struct {
	OnAdd func()
}

// Controller connects configuration with the reminder list view.
type Controller struct {
	view   *View
	onAdd  func()
}

// NewController creates a controller for the given view and wires callbacks.
func NewController(view *View, callbacks Callbacks) *Controller {
	if callbacks.OnAdd != nil {
		view.SetOnAdd(callbacks.OnAdd)
	}
	return &Controller{
		view:  view,
		onAdd: callbacks.OnAdd,
	}
}

// ApplyConfig converts config to the UI model and updates the view on the Fyne UI thread.
// If afterRender is non-nil, it is called after Render (on the same UI thread).
func (c *Controller) ApplyConfig(cfg *config.Config, afterRender func()) {
	model := ModelFromConfig(cfg)
	onAdd := c.onAdd
	fyne.Do(func() {
		c.view.Render(model, onAdd)
		if afterRender != nil {
			afterRender()
		}
	})
}
