package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

// E EmptyContainer New Instance
var E *EmptyContainer

// EmptyContainer : Container Beside The Nav Panel
type EmptyContainer struct {
	Pages        []Page
	CheckList    *CheckLists
	OrgCheckL    *OrgCheckList
	Cont         fyne.Container
	PeriodicSave bool //commit to db when page is changed
}

// Show : Show This Object
func (e *EmptyContainer) Show(obj fyne.CanvasObject) {
	for _, v := range e.Cont.Objects {
		e.Cont.Remove(v)
	}
	e.Cont.Add(obj)
	e.Cont.Refresh()
	if e.PeriodicSave {
		e.SaveState()
	}
}

// SaveState : Commit Changes to DB
func (e *EmptyContainer) SaveState() {
	for _, v := range e.Pages {
		v.Store()
	}
	e.CheckList.Store()
	e.OrgCheckL.Store()
	//not saving org by default has to do manually
}

// GetPage : Get Page by Its Name
func (e *EmptyContainer) GetPage(pagename string) (*Page, error) {
	for _, v := range e.Pages {
		if v.BtnName == pagename {
			return &v, nil
		}
	}

	return &Page{}, fmt.Errorf("Page Not Found")
}

// NewEmptyContainer ...
func NewEmptyContainer() *EmptyContainer {
	c := container.NewMax()
	e := EmptyContainer{
		Cont: *c,
	}
	return &e
}
