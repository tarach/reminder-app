package layout

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type TopLeftMargin struct {
	TopMargin  float32
	LeftMargin float32
}

type HasTopLeftMargin interface {
	GetTopLeftMargins() (float32, float32)
}

type TextWithMargin struct {
	widget.BaseWidget
	Text   *canvas.Text
	Margin TopLeftMargin
}

func (t *TextWithMargin) GetTopLeftMargins() (float32, float32) {
	return t.Margin.TopMargin, t.Margin.LeftMargin
}

func (t *TextWithMargin) CreateRenderer() fyne.WidgetRenderer {
	return &textWithMarginRenderer{text: t.Text, parent: t}
}

func (t *TextWithMargin) MinSize() fyne.Size {
	t.ExtendBaseWidget(t)
	return t.BaseWidget.MinSize()
}

type textWithMarginRenderer struct {
	text   *canvas.Text
	parent *TextWithMargin
}

func (r *textWithMarginRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *textWithMarginRenderer) MinSize() fyne.Size {
	return r.text.MinSize()
}

func (r *textWithMarginRenderer) Refresh() {
	r.text.Refresh()
}

func (r *textWithMarginRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *textWithMarginRenderer) Destroy() {}

func NewTextWithMargin(top float32, left float32, text string, color color.Color) *TextWithMargin {
	t := &TextWithMargin{
		Text: canvas.NewText(text, color),
		Margin: TopLeftMargin{
			TopMargin:  top,
			LeftMargin: left,
		},
	}
	t.ExtendBaseWidget(t)
	return t
}

type rowHeight struct {
	Height float32
}

func (rh *rowHeight) setHigher(height float32) {
	if rh.Height > height {
		return
	}
	rh.Height = height
}

func (rh *rowHeight) reset() {
	rh.Height = 0
}

type MarginFlowLayout struct {
	Width  float32
	Height float32
}

func (d *MarginFlowLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(d.Width, d.Height)
}

func (d *MarginFlowLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	// Start position with a margin offset
	pos := fyne.NewPos(0, containerSize.Height-d.Height)
	availableWidth := containerSize.Width
	rh := rowHeight{Height: 0}

	fmt.Println("containerSize: ", containerSize)

	for _, o := range objects {
		marginTop := float32(0)
		marginLeft := float32(0)
		owm, ok := o.(HasTopLeftMargin)
		if ok {
			marginTop, marginLeft = owm.GetTopLeftMargins()
		}

		fmt.Println("")
		output := "(not a text object)"
		if textObj, ok := o.(*canvas.Text); ok {
			output = textObj.Text
		}
		if textWithMargin, ok := o.(*TextWithMargin); ok {
			output = textWithMargin.Text.Text
		}
		fmt.Println("Text: ", output)
		fmt.Println("Pos x: ", pos.X, " Pos y: ", pos.Y)
		fmt.Println("Margin left: ", marginLeft, " Margin top: ", marginTop)

		pos = pos.Add(fyne.NewPos(marginLeft, marginTop))
		rh.setHigher(pos.Y + o.MinSize().Height)

		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width, marginTop*-1))
		if pos.X > availableWidth {
			pos = fyne.NewPos(0, pos.Y+size.Height)
			o.Move(pos)
			pos = pos.Add(fyne.NewPos(size.Width, 0))
			rh.reset()
		}
	}
}
