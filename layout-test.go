package main

import (
	"reminder-app/pkg/ui/layout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func main() {
	a := app.New()
	w := a.NewWindow("Diagonal")
	w.SetPadded(false)

	texts := []fyne.CanvasObject{
		layout.NewTextWithMargin(15, 15, "topleft", nil),
		canvas.NewText("Middle Label", nil),
		canvas.NewText("brawrr asdfgqwer", nil),
		layout.NewTextWithMargin(15, 15, "11111111111", nil),
		canvas.NewText("22222222222", nil),
		layout.NewTextWithMargin(15, 15, "33333333333", nil),
		canvas.NewText("44444444444", nil),
		canvas.NewText("55555555555", nil),
		canvas.NewText("66666666666", nil),
		canvas.NewText("77777777777", nil),
		canvas.NewText("88888888888", nil),
		canvas.NewText("aaaaaaaaaaa", nil),
		canvas.NewText("abbbbbbbbbb", nil),
		canvas.NewText("ccccccccccc", nil),
		canvas.NewText("ddddddddddd", nil),
		canvas.NewText("eeeeeeeeeee", nil),
		canvas.NewText("fffffffffff", nil),
		canvas.NewText("ggggggggggg", nil),
		canvas.NewText("hhhhhhhhhhh", nil),
		canvas.NewText("iiiiiiiiiii", nil),
		canvas.NewText("jjjjjjjjjjj", nil),
	}

	w.SetContent(container.New(&layout.MarginFlowLayout{Width: 300, Height: 500}, texts...))
	w.ShowAndRun()
}
