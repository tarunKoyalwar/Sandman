package mylayouts

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//This is similar to hbox + Border containers
//but with manual intend

//this supports dynamic scaling (horizontally)
//if object min size is greater than available space
//intends will be divided by 2

// BorderBox : Custom Layout
type BorderBox struct {
	//These are ratios not values
	//Ex left = 6 is like use 1/6th available width as left intend
	//these can be zero
	LeftIntend  float32 //left intend
	RightIntend float32
	TopIntend   float32
	SpacerSize  int //times of default font size
	// MenuObject  fyne.CanvasObject //present at top right
}

// MinSize : This will give maximum of (min width of an object)
func (n *BorderBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		// fmt.Printf("Min child size is %v\n", childSize)

		if w == 0 {
			w = childSize.Width
		} else if w < childSize.Width {
			w = childSize.Width
		}
		h += childSize.Height
	}
	// fmt.Printf("Min Size required to render is %v\n", fyne.NewSize(w, h))
	return fyne.NewSize(w, h)
}

//will divide recursively until intend is satisfied or intend reach zero
func divinedivide(left float32, right float32, objmin float32, total float32) (float32, float32) {
	if left > 0 || right > 0 {
		if objmin <= total-(left+right) {
			return left, right
		} else {
			//recursive
			left = left / 2
			left = float32(math.Floor(float64(left)))
			right = right / 2
			right = float32(math.Floor(float64(right)))
			return divinedivide(left, right, objmin, total)
		}
	} else {
		return 0, 0
	}
}

// Layout : Defines Arranging of Objects
func (n *BorderBox) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {

	var right float32
	var left float32
	//remove logic errors
	if n.LeftIntend <= 0 {
		left = 0
	} else {
		left = containerSize.Width / n.LeftIntend
	}

	//New Addition removing right intend
	if n.RightIntend <= 0 {
		right = 0
	} else {
		right = containerSize.Width / n.RightIntend
	}

	//first get minsizes of objects

	objmin := n.MinSize(objects).Width
	top := containerSize.Height / n.TopIntend

	if left > 0 || right > 0 {
		//If any intend can be reduced call divinedivide
		if objmin > containerSize.Width-(left+right) {
			//invoke divinedivide
			left, right = divinedivide(left, right, objmin, containerSize.Width)
		}
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

	//render the menu canvas object
	// if n.MenuObject != nil {
	// 	menumin := n.MenuObject.MinSize()
	// 	n.MenuObject.Move(fyne.NewPos(containerSize.Width-menumin.Width, menumin.Height))

	// }
}

// MenuBotton : Menu Btn Widget
type MenuButton struct{}

// MinSize : Get MinSize
func (M *MenuButton) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()
		// fmt.Printf("Min child size is %v\n", childSize)

		if w == 0 {
			w = childSize.Width
		} else if w < childSize.Width {
			w = childSize.Width
		}
		h += childSize.Height
	}
	// fmt.Printf("Min Size required to render is %v\n", fyne.NewSize(w, h))
	return fyne.NewSize(w, h)
}

// Layout : Defines Arranging of Objects
func (M *MenuButton) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	menumin := objects[0].MinSize()
	objects[0].Move(fyne.NewPos(containerSize.Width-menumin.Width, menumin.Height))
}
