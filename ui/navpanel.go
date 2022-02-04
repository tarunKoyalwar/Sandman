package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// PanelObject ...
type PanelObject struct {
	Header  string
	Objects []fyne.CanvasObject
}

// Render ...
func (p *PanelObject) Render() fyne.CanvasObject {
	c := container.NewVBox(p.Objects...)
	card := widget.NewCard("", "", c)
	return card
}

//Navigation Panel
type NavPanel struct {
	Color   color.Color
	Objects []fyne.CanvasObject
}

// Render ...
func (N *NavPanel) Render() fyne.CanvasObject {
	length := len(N.Objects)

	for i := length; i < 6; i = i + 1 {
		N.Objects = append(N.Objects, layout.NewSpacer())
	}

	final := container.NewVBox(N.Objects...)

	withscroll := container.NewVScroll(final)
	rectangles := canvas.NewRectangle(N.Color)

	// rectangles.FillColor = color.RGBA{R: 49, G: 48, B: 66}
	max := container.NewMax(
		withscroll,
		rectangles,
	)

	return max
}

// UI_INIT ...
func UI_INIT() {
	E = NewEmptyContainer()
	obj := container.NewCenter(widget.NewLabel("Default View"))
	E.Show(obj)
}
