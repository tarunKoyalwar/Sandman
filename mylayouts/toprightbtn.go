package mylayouts

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

//adds top right button to container

//Adds object given to layout
//to top right in order

// MenuBtnlayout : Layout to add top right menu btn
type MenuBtnlayout struct {
}

// MinSize : Get MinSize
func (m *MenuBtnlayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childsize := o.MinSize()
		if h < childsize.Height {
			h = childsize.Height
		}
		w += childsize.Width
	}

	return fyne.NewSize(w, h)
}

// Layout : Defines Arranging of Objects
func (m *MenuBtnlayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	//n.MenuObject.Move(fyne.NewPos(containerSize.Width-menumin.Width, menumin.Height))

	//add all objects to vbox
	box := container.NewHBox(objects...)

	boxmin := box.MinSize()

	box.Resize(boxmin)

	box.Move(fyne.NewPos(containerSize.Width-boxmin.Width, boxmin.Height))
}
