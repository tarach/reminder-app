package remindercreate

import "fyne.io/fyne/v2"

// Callbacks are invoked by the reminder create screen.
type Callbacks struct {
	OnSave   func(Input)
	OnCancel func()
}

// Controller binds the create view to callbacks.
type Controller struct {
	view *View
}

// NewController creates a controller and wires the view to callbacks.
func NewController(view *View, callbacks Callbacks) *Controller {
	view.SetCallbacks(callbacks.OnSave, callbacks.OnCancel)
	return &Controller{view: view}
}

// Content returns the root content for the create screen (delegates to view).
func (c *Controller) Content() fyne.CanvasObject {
	return c.view.Content()
}
