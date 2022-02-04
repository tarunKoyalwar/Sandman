package mywidgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// TapLabel : label that can be tapped
type TapLabel struct {
	widget.Label
	OnTapped func()
}

// Tapped : Mouse Event Left Click will call the function
func (T *TapLabel) Tapped(_ *fyne.PointEvent) {
	if T.OnTapped != nil {
		T.OnTapped()
	}
}

// NewTapLabel ...
func NewTapLabel(text string, tapped func()) *TapLabel {
	t := TapLabel{
		OnTapped: tapped,
	}
	t.SetText(text)
	t.ExtendBaseWidget(&t)

	return &t
}
