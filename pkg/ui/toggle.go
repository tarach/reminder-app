package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Toggle switch colors - grey when off, primary when on
var (
	toggleOffBg  = color.NRGBA{R: 0xD1, G: 0xD5, B: 0xDB, A: 0xFF} // light grey
	toggleOnBg   = color.NRGBA{R: 0x0D, G: 0x94, B: 0x88, A: 0xFF} // teal primary
	toggleKnobBg = color.White
)

// ToggleSwitch is a custom vertical toggle switch widget (off = bottom, on = top).
type ToggleSwitch struct {
	widget.BaseWidget
	Checked   bool
	OnChanged func(checked bool)

	track *canvas.Rectangle
	knob  *canvas.Rectangle // Now a rounded rectangle knob
}

// NewToggleSwitch creates a new vertical toggle switch.
func NewToggleSwitch(checked bool, onChanged func(checked bool)) *ToggleSwitch {
	track := canvas.NewRectangle(toggleOffBg)
	if checked {
		track.FillColor = toggleOnBg
	}
	track.CornerRadius = 10 // rounded ends for vertical track

	knob := canvas.NewRectangle(toggleKnobBg)
	knob.CornerRadius = 6 // rounded corners for square-ish knob

	t := &ToggleSwitch{
		Checked:   checked,
		OnChanged: onChanged,
		track:     track,
		knob:      knob,
	}
	t.ExtendBaseWidget(t)
	return t
}

// SetChecked updates the checked state and refreshes the widget.
func (t *ToggleSwitch) SetChecked(checked bool) {
	if t.Checked != checked {
		t.Checked = checked
		t.Refresh()
	}
}

// Tapped implements fyne.Tappable - toggles on tap.
func (t *ToggleSwitch) Tapped(*fyne.PointEvent) {
	t.Checked = !t.Checked
	if t.OnChanged != nil {
		t.OnChanged(t.Checked)
	}
	t.Refresh()
}

// CreateRenderer implements fyne.Widget.
func (t *ToggleSwitch) CreateRenderer() fyne.WidgetRenderer {
	return &toggleRenderer{
		toggle: t,
		track:  t.track,
		knob:   t.knob,
	}
}

type toggleRenderer struct {
	toggle *ToggleSwitch
	track  *canvas.Rectangle
	knob   *canvas.Rectangle // Now a rounded rectangle
}

func (r *toggleRenderer) Layout(size fyne.Size) {
	r.track.Resize(size)
	r.track.Move(fyne.NewPos(0, 0))

	// Knob: square with rounded corners, fits inside vertical track
	padding := float32(4)
	availW := size.Width - padding*2
	availH := size.Height - padding*2
	side := availW
	if availH < side {
		side = availH
	}
	if side < 12 {
		side = 12
	}

	r.knob.Resize(fyne.NewSize(side, side))

	// Position knob: top (on) or bottom (off) — vertical layout
	var knobY float32
	if r.toggle.Checked {
		knobY = padding // on = knob at top
	} else {
		knobY = size.Height - side - padding // off = knob at bottom
	}
	knobX := padding + (availW-side)/2 // center horizontally
	r.knob.Move(fyne.NewPos(knobX, knobY))
}

func (r *toggleRenderer) MinSize() fyne.Size {
	// Vertical toggle: track is taller than wide (portrait)
	return fyne.NewSize(28, 56)
}

func (r *toggleRenderer) Refresh() {
	if r.toggle.Checked {
		r.track.FillColor = toggleOnBg
	} else {
		r.track.FillColor = toggleOffBg
	}
	r.track.Refresh()
	r.knob.Refresh()
	r.Layout(r.toggle.Size())
}

func (r *toggleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.track, r.knob}
}

func (r *toggleRenderer) Destroy() {}
