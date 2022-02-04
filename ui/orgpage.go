package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mytree"
	"github.com/tarunKoyalwar/sandman/mywidgets"
)

// OrgCheckList will be useful for organization/OSINT checklist
type OrgCheckList struct {
	MyTree     *mytree.Mytree
	BtnName    string
	Identifier string
}

// Store ...
func (o *OrgCheckList) Store() {
	if o.MyTree != nil {
		o.MyTree.Store()
	}
}

// Update loading on startup
func (o *OrgCheckList) Update() {
	if o.Identifier == "" {
		o.Identifier = "orgdefault"
		// fmt.Println("Changed to orgdefault")
	}
	//fetch all global checklists
	db.MDB.GetCollection("checklists")

	arr := mytree.GetAllTrees()

	//check if there is a tree named "orgdefault"

	found := false

	for _, v := range arr {
		if v.TreeIdentifer == o.Identifier {
			found = true

		}
	}

	if found {
		fmt.Printf("Found Organization Checklist\n")
		o.MyTree = mytree.NewMytree(o.Identifier)
		for _, v := range arr {
			if v.TreeIdentifer == o.Identifier {
				o.MyTree.BackupInstance = &v
				o.MyTree.UpdateFromInstance()
			}
		}
	} else {
		fmt.Printf("Organization Checklist Not FOund\n")

	}
}

// Render ...
func (o *OrgCheckList) Render() *fyne.Container {
	if o.MyTree == nil {
		o.MyTree = mytree.NewMytree(o.Identifier)
	}

	btn := widget.NewButton("", func() {})
	heading := canvas.NewText("Organization CheckLists", theme.ForegroundColor())
	heading.TextSize = 16
	heading.Alignment = fyne.TextAlignLeading
	heading.TextStyle = fyne.TextStyle{Bold: true}

	boxed := container.NewBorder(container.NewHBox(btn, heading), nil, nil, nil, o.MyTree.Panel)

	splitbox := container.NewHSplit(o.MyTree.GetContainer(w), boxed)
	splitbox.SetOffset(0.25)

	return container.NewMax(splitbox)
}

// GetButton ...
func (o *OrgCheckList) GetButton() *mywidgets.CustomButton {
	btn := mywidgets.CustomButton{
		Text:      fmt.Sprintf("\t%v", o.BtnName),
		TextSize:  15,
		Alignment: fyne.TextAlignLeading,
		// FillColor: color.RGBA{R: 49, G: 48, B: 66},
		OnTapped: func() {
			E.Show(o.Render())
		},
	}

	// btn.ExtendBaseWidget(btn)

	return &btn
}

// NewOrgCheckList ...
func NewOrgCheckList() *OrgCheckList {
	t := mytree.NewMytree("orgdefault")
	t.TreeIdentifer = "orgdefault"
	//rgba(49,48,66,255)
	t.Color = color.RGBA{R: 49, G: 48, B: 66, A: 255}

	c := &OrgCheckList{
		BtnName:    "Org CheckList",
		MyTree:     t,
		Identifier: t.TreeIdentifer,
	}

	c.Update()
	return c

}
