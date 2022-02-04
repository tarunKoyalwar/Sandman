package mywidgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CustomButton
type CustomButton struct {
	//this widget acts both like
	//button and label and is a lot more customizable

	widget.BaseWidget
	Text      string
	TextSize  float32
	Alignment fyne.TextAlign
	TextStyle fyne.TextStyle
	FillColor color.Color
	TextColor color.Color
	Padding   float32
	Icon      fyne.Resource

	OnTapped func()
}

// Tapped ...
func (c *CustomButton) Tapped(*fyne.PointEvent) {
	if c.OnTapped != nil {
		c.OnTapped()
	}
}

// CreateRenderer ...
func (c *CustomButton) CreateRenderer() fyne.WidgetRenderer {
	c.ExtendBaseWidget(c)
	ctext := canvas.Text{}
	ctext.Text = c.Text
	if c.TextColor == nil {
		ctext.Color = theme.ForegroundColor()
	} else {
		ctext.Color = c.TextColor
	}
	ctext.Alignment = c.Alignment
	ctext.TextSize = c.TextSize

	crect := canvas.NewRectangle(theme.BackgroundColor())
	if c.FillColor != nil {
		crect.FillColor = c.FillColor
	}

	var p float32 = theme.Padding()
	if c.Padding != 0.0 {
		p = c.Padding
	}

	r := CustomButtonRender{
		Padding: p,
		Text:    &ctext,
		// Rect:    crect,
		Icon: canvas.NewImageFromResource(c.Icon),
	}

	return &r
}

// CustomButtonRender ...
type CustomButtonRender struct {
	fyne.WidgetRenderer
	Padding float32
	Text    *canvas.Text
	// Rect    *canvas.Rectangle
	Icon    *canvas.Image
	CButton CustomButton
}

// MinSize calculates the minimum size of a check.
// This is based on the contained text, the check icon and a standard amount of padding added.
func (r *CustomButtonRender) MinSize() fyne.Size {
	tmin := r.Text.MinSize()
	var iconsize float32 = 0
	if r.Icon != nil {
		iconsize = theme.IconInlineSize()
	}
	tnewmin := fyne.NewSize(3*r.Padding+iconsize+tmin.Width, tmin.Height+2*r.Padding)

	return tnewmin

}

// Layout : Layout the components of the check widget
func (r *CustomButtonRender) Layout(size fyne.Size) {

	//resize all objects
	// r.Rect.Resize(size)
	r.Text.Resize(r.Text.MinSize())

	// pos := fyne.NewPos(size.Width, size.Height)

	// r.Rect.Move(pos)

	//resize icon and move it before the text if not null
	hasicon := false

	if r.Icon.Resource != nil {
		hasicon = true
	}
	textmin := r.Text.MinSize()
	// if r.Icon != nil {
	// 	iconisnil = false
	// }
	// fmt.Println(r.Icon)

	// r.Icon.

	//The problem is with keeping data exactly at center
	//padding required for x and y

	//for handling alignment
	var px float32
	var py float32
	if r.Text.Alignment == fyne.TextAlignCenter {
		px = (size.Width - textmin.Width) / 2
		py = (size.Height - textmin.Height) / 2
	} else if r.Text.Alignment == fyne.TextAlignLeading {
		px = 0
		py = (size.Height - textmin.Height) / 2
	} else {
		px = size.Width - textmin.Width
		py = (size.Height - textmin.Height) / 2
	}

	if hasicon {
		picony := (textmin.Height - theme.IconInlineSize()) / 2
		piconx := picony
		shiftfactor := (2*piconx + theme.IconInlineSize()) / 2
		iconloc := fyne.NewPos(px+piconx-shiftfactor, py+picony)
		r.Icon.Resize(fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize()))
		r.Icon.Move(iconloc)

		texpos := fyne.NewPos(2*piconx+theme.IconInlineSize()+px-shiftfactor, py)
		r.Text.Move(texpos)
	} else {
		textpos := fyne.NewPos(px, py)
		r.Text.Move(textpos)
	}
}

// Objects ...
func (r *CustomButtonRender) Objects() []fyne.CanvasObject {
	if r.Icon != nil {
		return []fyne.CanvasObject{r.Icon, r.Text}
	} else {
		return []fyne.CanvasObject{r.Text}
	}
}

// Refresh ...
func (r *CustomButtonRender) Refresh() {
	r.Text.Text = r.CButton.Text
	r.Text.TextSize = r.CButton.TextSize
	r.Text.TextStyle = r.CButton.TextStyle
	if r.CButton.TextColor != nil {
		r.Text.Color = r.CButton.TextColor
	} else {
		r.Text.Color = theme.ForegroundColor()
	}
	r.Text.Refresh()
	// r.Rect.Refresh()

}

// Destroy ...
func (r *CustomButtonRender) Destroy() {}
