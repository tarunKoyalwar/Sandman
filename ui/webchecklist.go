package ui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mylayouts"
	"github.com/tarunKoyalwar/sandman/mytree"
	"github.com/tarunKoyalwar/sandman/mywidgets"
	"github.com/tarunKoyalwar/sandman/utils"
)

// CheckLists : Web CheckLists Page
type CheckLists struct {
	MyTrees         map[string]*mytree.Mytree //Tree for each subdomain
	BtnName         string                    //Panel Button Name
	ActiveSubs      *Page                     //All subdomains page
	ActiveTreeIndex string                    //Current Active tree
	Organized       binding.Bool              //Organize data
}

// Store : save data
func (c *CheckLists) Store() {
	for _, v := range c.MyTrees {
		v.Store()
	}
}

// Home : View All subdomains
func (c *CheckLists) Home() *fyne.Container {

	if c.Organized == nil {
		c.Organized = binding.NewBool()
	}

	vieworganized, _ := c.Organized.Get()

	checkbox := widget.NewCheckWithData("Organized", c.Organized)
	checkbox.SetChecked(vieworganized)

	definedlayout := &mylayouts.BorderBoxMin{
		TopIntend:   7,
		LeftIntend:  3,
		RightIntend: 3,
		SpacerSize:  3,
	}

	definedlayout.AddMenuObj(checkbox)

	prefilter := []string{}
	//Get active subdomains from active subs page
	//filter cidr and convert to  hosts
	//sample cidr 22.22.22.222/29
	all := c.ActiveSubs.Entry.GetData()

	label := widget.NewLabelWithStyle("Available Subdomains", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	// spacer := layout.NewSpacer()

	//uses regex to parse CIDR's if exists
	cidrs, noncidrs := utils.RegexGetCIDRs(all)

	//converts multiple cidrs to ip hosts
	iplist := utils.MultiCIDRtoHosts(cidrs)

	for _, b := range strings.Split(noncidrs, "\n") {
		if strings.TrimSpace(b) != "" {
			prefilter = append(prefilter, b)
		}
	}

	GetView := func(organized bool) fyne.CanvasObject {

		if organized {

			tldsize := 2
			filtered := prefilter
			//will create a organized layout
			//tldsize is 2 for abc.hackerone.com if hackerone.com is to be considered as apex
			organized := utils.OrganizeData(filtered, tldsize)
			if len(cidrs) != 0 {
				for _, x := range cidrs {
					organized[x] = utils.CIDRHosts(x)
				}
			}

			// fmt.Printf("[Debug] Organized data %v\n", organized)

			if len(organized) == 0 {
				organized["domains"] = filtered
			}

			selectors := []string{}

			for k := range organized {
				selectors = append(selectors, k)
			}

			listholder := container.NewMax()

			newlist := func(apex string) *widget.List {
				listdata := organized[apex]

				//create list of all results
				list := widget.NewList(
					//length function
					func() int {
						return len(listdata)
					},
					func() fyne.CanvasObject {
						return widget.NewLabel("template")
					},
					func(i widget.ListItemID, o fyne.CanvasObject) {
						o.(*widget.Label).SetText(listdata[i])
					})

				list.OnSelected = func(id widget.ListItemID) {
					fmt.Printf("Displaying CheckList for %v\n", listdata[id])
					c.ActiveTreeIndex = listdata[id]
					E.Show(c.Render())
				}

				return list
			}

			selectorwid := widget.NewSelect(selectors, func(s string) {
				for _, v := range listholder.Objects {
					listholder.Remove(v)
				}
				listholder.Add(newlist(s))
			})

			bordered := container.NewBorder(selectorwid, nil, nil, nil, listholder)

			return container.NewVScroll(bordered)

		} else {
			filtered := prefilter
			//add all ipaddress to arr
			filtered = append(filtered, iplist...)

			//create list of all results
			list := widget.NewList(
				//length function
				func() int {
					return len(filtered)
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("template")
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					o.(*widget.Label).SetText(filtered[i])
				})

			list.OnSelected = func(id widget.ListItemID) {
				fmt.Printf("Displaying CheckList for %v\n", filtered[id])
				c.ActiveTreeIndex = filtered[id]
				E.Show(c.Render())
			}

			return container.NewVScroll(list)

		}
	}

	prerender := container.New(definedlayout, label, checkbox, GetView(false))

	checkbox.OnChanged = func(b bool) {

		//view organized
		for _, v := range prerender.Objects {
			switch v.(type) {
			case *container.Scroll:
				// fmt.Printf("[debug] FOund vscroll \n")
				prerender.Remove(v)
				prerender.Add(GetView(b))
				prerender.Refresh()
			}
		}

	}

	return prerender

}

// Update ...
func (c *CheckLists) Update() {
	//Load from checklists collection
	db.MDB.GetCollection("checklists")

	arr := mytree.GetAllTrees()

	for _, v := range arr {
		if v.TreeIdentifer == "orgdefault" {
			continue
		}
		c.MyTrees[v.TreeIdentifer] = mytree.NewMytree(v.TreeIdentifer)
		c.MyTrees[v.TreeIdentifer].BackupInstance = &v
		c.MyTrees[v.TreeIdentifer].UpdateFromInstance()
		fmt.Printf("Found Backup for %v\n", v.TreeIdentifer)
	}

}

// Render ...
func (c *CheckLists) Render() *fyne.Container {
	if c.MyTrees[c.ActiveTreeIndex] == nil {
		fmt.Println("No active tree present creating new entry")
		c.MyTrees[c.ActiveTreeIndex] = mytree.NewMytree(c.ActiveTreeIndex)
		fmt.Println(c.MyTrees[c.ActiveTreeIndex].TreeIdentifer)
	}

	heading := canvas.NewText("Web App CheckLists", theme.ForegroundColor())
	heading.TextSize = 16
	heading.Alignment = fyne.TextAlignLeading
	heading.TextStyle = fyne.TextStyle{Bold: true}
	btn := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		E.Show(c.Home())
	})

	boxed := container.NewBorder(container.NewHBox(btn, heading), nil, nil, nil, c.MyTrees[c.ActiveTreeIndex].Panel)

	splitbox := container.NewHSplit(c.MyTrees[c.ActiveTreeIndex].GetContainer(w), boxed)
	splitbox.SetOffset(0.25)

	return container.NewMax(splitbox)

}

// GetButton ...
func (c *CheckLists) GetButton() *mywidgets.CustomButton {
	btn := mywidgets.CustomButton{
		Text:      fmt.Sprintf("\t%v", c.BtnName),
		TextSize:  15,
		Alignment: fyne.TextAlignLeading,
		// FillColor: color.RGBA{R: 49, G: 48, B: 66},
		OnTapped: func() {
			E.Show(c.Render())
		},
	}

	return &btn
}

// GetDefaultCheckList ...
func (c *CheckLists) GetDefaultCheckList() map[string][]string {
	tree := mytree.NewMytree("default")
	return tree.Nodes
}

// NewCheckLists ...
func NewCheckLists() *CheckLists {

	t := mytree.NewMytree("default")
	t.TreeIdentifer = "Default"
	//rgba(49,48,66,255)
	t.Color = color.RGBA{R: 49, G: 48, B: 66, A: 255}

	c := &CheckLists{
		BtnName:         "Web CheckLists",
		MyTrees:         map[string]*mytree.Mytree{"default": t},
		ActiveTreeIndex: "default",
	}
	if c.Organized == nil {
		c.Organized = binding.NewBool()
	}
	c.Update()
	return c
}
