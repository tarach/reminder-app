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

// ToggleSwitch is a custom toggle switch widget (off = left, on = right).
type ToggleSwitch struct {
	widget.BaseWidget
	Checked   bool
	OnChanged func(checked bool)

	track *canvas.Rectangle
	knob  *canvas.Circle
}

// NewToggleSwitch creates a new toggle switch.
func NewToggleSwitch(checked bool, onChanged func(checked bool)) *ToggleSwitch {
	track := canvas.NewRectangle(toggleOffBg)
	if checked {
		track.FillColor = toggleOnBg
	}
	track.CornerRadius = 5 // semicircular ends (half of height)
	t := &ToggleSwitch{
		Checked:   checked,
		OnChanged: onChanged,
		track:     track,
		knob:      canvas.NewCircle(toggleKnobBg),
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
	knob   *canvas.Circle
}

func (r *toggleRenderer) Layout(size fyne.Size) {
	r.track.Resize(size)
	r.track.Move(fyne.NewPos(0, 0))

	// Knob: ~50% of track height, centered
	knobSize := size.Height / 2
	if knobSize < 4 {
		knobSize = 4
	}
	r.knob.Resize(fyne.NewSize(knobSize, knobSize))

	// Position knob: left when off, right when on
	knobY := (size.Height - knobSize) / 2
	var knobX float32
	if r.toggle.Checked {
		knobX = size.Width - knobSize - 2
	} else {
		knobX = 2
	}
	r.knob.Move(fyne.NewPos(knobX, knobY))
}

func (r *toggleRenderer) MinSize() fyne.Size {
	// Compact toggle: wider track, slimmer height
	return fyne.NewSize(60, 20)
}

func (r *toggleRenderer) Refresh() {
	if r.toggle.Checked {
		r.track.FillColor = toggleOnBg
	} else {
		r.track.FillColor = toggleOffBg
	}
	r.track.Refresh()
	r.Layout(r.toggle.Size())
}

func (r *toggleRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.track, r.knob}
}

func (r *toggleRenderer) Destroy() {}
