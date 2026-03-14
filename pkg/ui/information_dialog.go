package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	infoPadWidth         = 32
	infoPadHeight        = 16
	infoMessageButtonGap = 24 // vertical space between message content and OK button
)

// InformationDialogMinWidth is the minimum width of the information dialog.
// Increase this to widen the dialog.
var InformationDialogMinWidth float32 = 420

// InformationDialogIconOffsetRight is the horizontal inset of the icon from the right edge.
var InformationDialogIconOffsetRight float32 = 12

// InformationDialogIconOffsetTop is the vertical inset of the icon from the top edge.
var InformationDialogIconOffsetTop float32 = 12

// informationDialogLayout lays out the information dialog content.
type informationDialogLayout struct{}

func (l *informationDialogLayout) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	btnMin := obj[3].MinSize()
	labelMin := obj[4].MinSize()

	iconHeight := infoPadHeight*2 + labelMin.Height*2 - theme.Padding()

	// icon: top right, inset by configurable offsets from right and top
	obj[0].Resize(fyne.NewSize(iconHeight, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width-iconHeight-InformationDialogIconOffsetRight, InformationDialogIconOffsetTop))

	// background
	obj[1].Move(fyne.NewPos(0, 0))
	obj[1].Resize(size)

	// title (top)
	obj[4].Move(fyne.NewPos(infoPadWidth/2, infoPadHeight))
	obj[4].Resize(fyne.NewSize(size.Width-infoPadWidth, labelMin.Height))

	// content (message)
	contentStart := infoPadHeight + labelMin.Height + infoPadHeight
	contentEnd := size.Height - infoPadHeight - btnMin.Height - infoMessageButtonGap
	obj[2].Move(fyne.NewPos(infoPadWidth/2, contentStart))
	obj[2].Resize(fyne.NewSize(size.Width-infoPadWidth, contentEnd-contentStart))

	// buttons
	obj[3].Resize(btnMin)
	obj[3].Move(fyne.NewPos(size.Width/2-(btnMin.Width/2), size.Height-infoPadHeight-btnMin.Height))
}

func (l *informationDialogLayout) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[2].MinSize()
	btnMin := obj[3].MinSize()
	labelMin := obj[4].MinSize()

	naturalWidth := fyne.Max(fyne.Max(contentMin.Width, btnMin.Width), labelMin.Width) + infoPadWidth
	width := fyne.Max(naturalWidth, InformationDialogMinWidth)
	// Three padding zones: top, between title and content, bottom; plus explicit gap above button
	height := contentMin.Height + labelMin.Height + btnMin.Height + infoMessageButtonGap + infoPadHeight*3

	return fyne.NewSize(width, height)
}

// ShowInformation shows a dialog with title and message and an info icon in the top right.
// Width and icon position can be adjusted via InformationDialogMinWidth,
// InformationDialogIconOffsetRight, and InformationDialogIconOffsetTop.
func ShowInformation(title, message string, parent fyne.Window) {
	iconRes := theme.InfoIcon()
	icon := &canvas.Image{Resource: iconRes}

	bg := canvas.NewRectangle(theme.Color(theme.ColorNameOverlayBackground))
	if n, ok := theme.Color(theme.ColorNameOverlayBackground).(color.NRGBA); ok {
		bg.FillColor = &color.NRGBA{R: n.R, G: n.G, B: n.B, A: 230}
	}

	messageLabel := &widget.Label{
		Text:      message,
		Alignment: fyne.TextAlignCenter,
		Wrapping:  fyne.TextWrapWord,
	}

	okBtn := widget.NewButton("OK", nil)

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.New(&informationDialogLayout{},
		icon,
		bg,
		messageLabel,
		okBtn,
		titleLabel,
	)

	pop := widget.NewModalPopUp(content, parent.Canvas())
	okBtn.OnTapped = pop.Hide
	pop.Show()
}
