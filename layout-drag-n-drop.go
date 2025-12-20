package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// TappableRectangle wraps a canvas.Rectangle and makes it tappable and draggable
type TappableRectangle struct {
	widget.BaseWidget
	rect       *canvas.Rectangle
	onTapped   func()
	dragOffset fyne.Position
}

func NewTappableRectangle(color color.Color, onTapped func()) *TappableRectangle {
	tr := &TappableRectangle{
		rect:     canvas.NewRectangle(color),
		onTapped: onTapped,
	}
	tr.ExtendBaseWidget(tr)
	return tr
}

func (tr *TappableRectangle) Tapped(*fyne.PointEvent) {
	if tr.onTapped != nil {
		tr.onTapped()
	}
}

func (tr *TappableRectangle) Dragged(event *fyne.DragEvent) {
	// Calculate new position based on drag
	newPos := tr.Position().Add(event.Dragged)
	tr.Move(newPos)
}

func (tr *TappableRectangle) DragEnd() {
	tr.Tapped(nil)
}

func (tr *TappableRectangle) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(tr.rect)
}

func (tr *TappableRectangle) GetRectangle() *canvas.Rectangle {
	return tr.rect
}

func main() {
	// Create a new application
	myApp := app.New()
	myWindow := myApp.NewWindow("Desktop App")

	// Set window size to 400x600
	myWindow.Resize(fyne.NewSize(400, 600))

	// Create a label to display position
	positionLabel := widget.NewLabel("Click the rectangle!")

	// Create a tappable rectangle at position 20x20
	var tappableRect *TappableRectangle
	tappableRect = NewTappableRectangle(color.RGBA{R: 100, G: 100, B: 200, A: 255}, func() {
		position := tappableRect.Position()
		size := tappableRect.Size()
		positionLabel.SetText(fmt.Sprintf("Rectangle Position: X=%.0f, Y=%.0f\nSize: W=%.0f, H=%.0f",
			position.X, position.Y, size.Width, size.Height))
	})

	// Position and size the rectangle
	tappableRect.Resize(fyne.NewSize(50, 50))
	tappableRect.Move(fyne.NewPos(20, 20))

	// Create a container with absolute positioning
	content := container.NewWithoutLayout(tappableRect)

	// Create main layout with rectangle and label
	mainContainer := container.NewBorder(nil, positionLabel, nil, nil, content)

	// Create 1px borders
	topBorder := canvas.NewRectangle(color.Black)
	topBorder.SetMinSize(fyne.NewSize(1, 1))

	bottomBorder := canvas.NewRectangle(color.Black)
	bottomBorder.SetMinSize(fyne.NewSize(1, 1))

	leftBorder := canvas.NewRectangle(color.Black)
	leftBorder.SetMinSize(fyne.NewSize(1, 1))

	rightBorder := canvas.NewRectangle(color.Black)
	rightBorder.SetMinSize(fyne.NewSize(1, 1))

	// Wrap content with 1px border
	borderedContent := container.NewBorder(topBorder, bottomBorder, leftBorder, rightBorder, mainContainer)

	myWindow.SetContent(borderedContent)
	myWindow.ShowAndRun()
}
