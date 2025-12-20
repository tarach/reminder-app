package layout

import (
	"fyne.io/fyne/v2"
)

type SimpleFlowLayout struct {
	Margin float32
	Width  float32
	Height float32
}

func (d *SimpleFlowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(d.Width, d.Height)
}

func (d *SimpleFlowLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	margin := d.Margin
	if margin == 0 {
		margin = 10 // Default margin
	}

	// Start position with a margin offset
	pos := fyne.NewPos(margin, containerSize.Height-d.MinSize(objects).Height+margin)
	availableWidth := containerSize.Width - (margin * 2) // Account for left and right margins

	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		newXPos := size.Width + margin
		pos = pos.Add(fyne.NewPos(newXPos, 0))
		if pos.X > availableWidth {
			pos = fyne.NewPos(margin, pos.Y+size.Height)
			o.Move(pos)
			pos = pos.Add(fyne.NewPos(newXPos, 0))
		}

	}
}
