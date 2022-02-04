package mywidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// ContextMenuButton : Drop Down Menu
type ContextMenuButton struct {
	widget.Button
	menu *fyne.Menu
}

//Tapped : Mouse Event
func (C *ContextMenuButton) Tapped(e *fyne.PointEvent) {
	// fmt.Println("Tapped")
	widget.ShowPopUpMenuAtPosition(C.menu, fyne.CurrentApp().Driver().CanvasForObject(C), e.AbsolutePosition)
}

//Tapped : Right Click Mouse Event
func (C *ContextMenuButton) TappedSecondary(*fyne.PointEvent) {}

// NewContextMenu ...
func NewContextMenu(icon fyne.Resource, m *fyne.Menu) *ContextMenuButton {
	c := &ContextMenuButton{
		menu: m,
	}
	c.Icon = icon
	c.ExtendBaseWidget(c)
	return c
}
