package mylayouts

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BorderBoxMin : Custom Layout
type BorderBoxMin struct {
	//This will also add intend to Minsize
	//These are ratios not values
	//Ex left = 6 is like use 1/6th available width as left intend
	//these can be zero
	LeftIntend  float32 //left intend
	RightIntend float32
	TopIntend   float32
	SpacerSize  int               //times of default font size
	MenuObject  fyne.CanvasObject //present at top right
	Padding     int               //This padding will be added to min size
}

// AddMenuObj : Adds Custom Menu Object
func (n *BorderBoxMin) AddMenuObj(z fyne.CanvasObject) {
	n.MenuObject = z
	// fmt.Printf("[debug] added obj %v\n", n.MenuObject)
}

// MinSize : This will give maximum of (min width of an object)
func (n *BorderBoxMin) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		// fmt.Printf("Min child size is %v\n", childSize)

		if _, ok := o.(*layout.Spacer); ok {
			h += theme.TextSize() * float32(n.SpacerSize)
			continue
		}

		if w == 0 {
			w = childSize.Width
		} else if w < childSize.Width {
			w = childSize.Width
		}
		h += childSize.Height
	}
	if n.Padding != 0 {
		h += float32(n.Padding) * theme.TextSize()
	}
	// fmt.Printf("Min Size required to render is %v\n", fyne.NewSize(w, h))
	return fyne.NewSize(w, h)
}

// Layout : Defines Arranging of Objects
func (n *BorderBoxMin) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {

	//remove logic errors
	if n.LeftIntend <= 1 {
		n.LeftIntend = 1
	}

	//first get minsizes of objects
	left := containerSize.Width / n.LeftIntend
	right := containerSize.Width / n.RightIntend
	objmin := n.MinSize(objects).Width
	top := containerSize.Height / n.TopIntend

	if left >= 0 && right >= 0 {
		if objmin > containerSize.Width-(left+right) {
			//invoke divinedivide
			left, right = divinedivide(left, right, objmin, containerSize.Width)
		}
	} else {
		fyne.LogError("Intend Cannot be negative", fmt.Errorf("intend cannot be negative"))
		return
	}

	pos := fyne.NewPos(0+left, 0+top)
	for index, o := range objects {
		osize := o.MinSize()
		//this layout is like grid columns width is max and height is min
		scaledobjheight := osize.Height
		scaledobjwidth := containerSize.Width - (left + right)

		if index == len(objects)-1 {
			//Last Item scale it to fit size
			switch o.(type) {
			case *widget.List:
				o.Resize(fyne.NewSize(scaledobjwidth, containerSize.Height-pos.Y))
			case *widget.TextGrid:
				o.Resize(fyne.NewSize(scaledobjwidth, containerSize.Height-pos.Y))
			case *widget.Tree:
				o.Resize(fyne.NewSize(scaledobjwidth, containerSize.Height-pos.Y))
			case *container.Scroll:
				o.Resize(fyne.NewSize(scaledobjwidth, containerSize.Height-pos.Y))
			case *widget.Entry:
				if osize.Height < (containerSize.Height-pos.Y)/2 {
					o.Resize(fyne.NewSize(scaledobjwidth, containerSize.Height-pos.Y))
				} else {
					o.Resize(fyne.NewSize(scaledobjwidth, osize.Height))
				}

			default:
				o.Resize(fyne.NewSize(scaledobjwidth, scaledobjheight))
			}

		} else {

			o.Resize(fyne.NewSize(scaledobjwidth, scaledobjheight))
		}

		if _, ok := o.(*layout.Spacer); ok {
			// fmt.Printf("resized spacer to %v\n", fyne.NewSize(scaledobjwidth, theme.TextSize()*float32(n.SpacerSize)))
			scaledobjheight = theme.TextSize() * float32(n.SpacerSize)
			o.Resize(fyne.NewSize(scaledobjwidth, theme.TextSize()*float32(n.SpacerSize)))
		}

		o.Move(pos)
		pos = pos.Add(fyne.NewSize(0, scaledobjheight))
	}

	// fmt.Printf("[Debug] Menu Object Not found \n")

	//render the menu canvas object
	//Heres a trick container avoids rendering duplicates
	//so if we add object twice in objects and as an menu object
	//final position of that object is changed to this
	if n.MenuObject != nil {
		menumin := n.MenuObject.MinSize()
		n.MenuObject.Move(fyne.NewPos(containerSize.Width-menumin.Width, menumin.Height))

	}
}
