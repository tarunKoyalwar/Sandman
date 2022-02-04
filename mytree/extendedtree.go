package mytree

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mylayouts"
	"github.com/tarunKoyalwar/sandman/mywidgets"
	"go.mongodb.org/mongo-driver/bson"
)

//Default USe this while creating CHecklists
var Default *CheckListOnly

// OrgCheck ...
var OrgCheck *CheckListOnly

// TreeBackup ...
type TreeBackup struct {
	Nodes         map[string][]string    `bson:"nodes"`
	BCheckBoxes   map[string]bool        `bson:"bcheckboxes"`
	ItemStore     map[string]*TreeObject `bson:"itemstore"`
	TreeIdentifer string                 `bson:"treeidentifer"`
	ProgressBar   float64                `bson:"progressbar"`
}

// CheckListOnly ...
type CheckListOnly struct {
	Name         string                 `bson:"checklistname"`
	TreeNodes    map[string][]string    `bson:"treenodes"`
	Descriptions map[string]*TreeObject `bson:"description"`
}

// Mytree ...
type Mytree struct {
	widget.Tree
	Nodes          map[string][]string
	Panel          *fyne.Container
	CheckBoxes     map[string]binding.Bool
	BCheckBoxes    map[string]bool
	Color          color.Color
	ItemStore      map[string]*TreeObject
	StatusBar      binding.Float
	Bar            *widget.ProgressBar
	TreeIdentifer  string
	Coll           string
	BackupInstance *TreeBackup
	HashMap        map[string]string //only use is for fast indexing
}

// ExportAs : function to export checklist under different name
func (m *Mytree) ExportAs() {
	defined := mylayouts.BorderBoxMin{
		TopIntend:   8,
		LeftIntend:  8,
		RightIntend: 8,
		SpacerSize:  2,
		Padding:     4,
	}
	dat1 := binding.NewString()
	e1 := widget.NewEntryWithData(dat1)
	e1.SetText("default")

	label := widget.NewLabel("Export As")

	content := container.New(&defined, layout.NewSpacer(), label, e1, layout.NewSpacer())

	callbacks := func(b bool) {
		if b {
			n, _ := dat1.Get()
			m.ExportCheckList(n)
		}
	}

	dx := dialog.NewCustomConfirm("Export CheckList", "Done", "Cancel", content, callbacks, Win)
	dx.Show()
}

// ExportCheckList ...
func (m *Mytree) ExportCheckList(name string) {
	z := CheckListOnly{}
	z.Name = name
	z.TreeNodes = m.Nodes
	descriptions := map[string]*TreeObject{}

	for k, v := range m.ItemStore {
		d1 := TreeObject{}
		d1.HeaderText = v.HeaderText
		d1.Notes = v.Notes
		d1.Payload = v.Payload
		descriptions[k] = &d1
	}

	z.Descriptions = descriptions
	PutCheckList(z)
}

//Done
func (m *Mytree) ImportCheckList(d CheckListOnly) {
	m.Nodes = d.TreeNodes
	m.ItemStore = d.Descriptions
	for _, v := range m.ItemStore {
		v.Update()
	}
	m.Refresh()
}

// SaveInstance ...
func (m *Mytree) SaveInstance() {
	m.BackupChecks()

	for _, v := range m.ItemStore {
		v.Store()
	}

	bk := TreeBackup{}
	bk.Nodes = m.Nodes
	bk.BCheckBoxes = m.BCheckBoxes
	bk.ItemStore = m.ItemStore
	bk.TreeIdentifer = m.TreeIdentifer
	bk.ProgressBar = m.Bar.Value

	m.BackupInstance = &bk
}

// UpdateFromInstance ...
func (m *Mytree) UpdateFromInstance() {
	ptr := m.BackupInstance
	m.BCheckBoxes = ptr.BCheckBoxes
	m.ItemStore = ptr.ItemStore
	for _, v := range m.ItemStore {
		v.Update()
	}
	m.Nodes = ptr.Nodes
	m.TreeIdentifer = ptr.TreeIdentifer

	m.UpdateChecks()
	m.Refresh()
	if m.Bar == nil {
		m.StatusBar.Set(ptr.ProgressBar)
	} else {
		m.Bar.SetValue(ptr.ProgressBar)
	}

	//check if checklist is default and update latest default checklist
	if m.TreeIdentifer == "default" {
		arr := GetAllCheckLists()
		var DefChecklist *CheckListOnly
		for _, v := range arr {
			if v.Name == "default" {
				DefChecklist = &v
			}
		}
		if DefChecklist != nil {
			//update default tree
			m.ImportCheckList(*DefChecklist)
		}

	}
}

// UpdateUsingInterface ...
func (m *Mytree) UpdateUsingInterface(dat interface{}) {
	ptr := dat.(*TreeBackup)
	m.BCheckBoxes = ptr.BCheckBoxes
	m.ItemStore = ptr.ItemStore
	for _, v := range m.ItemStore {
		v.Update()
	}
	m.Nodes = ptr.Nodes
	m.TreeIdentifer = ptr.TreeIdentifer

	m.UpdateChecks()
	m.Refresh()
	if m.Bar == nil {
		m.StatusBar.Set(ptr.ProgressBar)
	} else {
		m.Bar.SetValue(ptr.ProgressBar)
	}
	fmt.Println("Updated using Interface")
}

// UpdateFromDB ...
func (m *Mytree) UpdateFromDB() {
	var z TreeBackup
	db.MDB.GetCollection(m.Coll)
	filter := bson.M{"treeidentifer": m.TreeIdentifer}
	dat, err := db.MDB.FindOne(filter, &z)
	if err != nil {
		panic(err)
	}
	ptr := dat.(*TreeBackup)
	m.BCheckBoxes = ptr.BCheckBoxes
	m.ItemStore = ptr.ItemStore
	for _, v := range m.ItemStore {
		v.Update()
	}
	m.Nodes = ptr.Nodes
	m.TreeIdentifer = ptr.TreeIdentifer

	m.UpdateChecks()
	m.Refresh()
	if m.Bar == nil {
		m.StatusBar.Set(ptr.ProgressBar)
	} else {
		m.Bar.SetValue(ptr.ProgressBar)
	}
}

// Store ...
func (m *Mytree) Store() {
	//backup boxes
	m.BackupChecks()

	for _, v := range m.ItemStore {
		v.Store()
	}

	bk := TreeBackup{}
	bk.Nodes = m.Nodes
	bk.BCheckBoxes = m.BCheckBoxes
	bk.ItemStore = m.ItemStore
	bk.TreeIdentifer = m.TreeIdentifer
	if m.Bar != nil {
		bk.ProgressBar = m.Bar.Value
	}

	//actual saving to db
	if m.Coll != "" {
		db.MDB.GetCollection(m.Coll)
		//create interface and store it
		filter := bson.M{"treeidentifer": m.TreeIdentifer}
		res, err := db.MDB.UpdateDocument(filter, bson.M{"$set": bk})
		if err != nil {
			// fmt.Println(res)
			// panic(err)
			e := dialog.NewError(fmt.Errorf("failed to save checklist %v", err.Error()), Win)
			e.Show()
			e.SetOnClosed(fyne.CurrentApp().Quit)

		}
		if res != nil {
			if res.MatchedCount == 0 {
				fmt.Printf("Document Not Found Creating New One\n")
			} else {
				fmt.Printf("Updated Existing Document\n")
			}
		}
	}
}

// UpdateChecks ...
func (m *Mytree) UpdateChecks() {
	for k, v := range m.BCheckBoxes {
		if m.CheckBoxes[k] == nil {
			m.CheckBoxes[k] = binding.NewBool()
		}
		m.CheckBoxes[k].Set(v)
	}
}

// BackupChecks ...
func (m *Mytree) BackupChecks() {
	if m.BCheckBoxes == nil {
		m.BCheckBoxes = map[string]bool{}
	}

	for k, v := range m.CheckBoxes {
		m.BCheckBoxes[k], _ = v.Get()
	}
}

// UpdateBar ...
func (m *Mytree) UpdateBar() {
	total := len(m.CheckBoxes)
	if len(m.CheckBoxes) == 0 {
		total = 1
	}

	positive := 0
	for _, v := range m.CheckBoxes {
		if stat, _ := v.Get(); stat {
			positive = positive + 1
		}
	}

	// fmt.Printf("Total is %v and positive is %v\n", total, positive)
	// fmt.Printf("Total positive %v Total Negative %v\n", positive, total)
	val := float64(positive) / float64(total)
	m.StatusBar.Set(val)
	if m.Bar != nil {
		m.Bar.SetValue(val)
		m.Bar.Refresh()
	}
	// fmt.Printf("Bar Updated new value is %v\n", val)
}

// UpdateCheckStatus : Changing Status of Branch (Only works for a Branch Not nested ones)
func (m *Mytree) UpdateCheckStatus(uid string, status bool) {
	m.CheckBoxes[uid].Set(status)
	for _, v := range m.Nodes[uid] {
		if m.CheckBoxes[v] != nil {
			m.CheckBoxes[v].Set(status)
		}
	}
	m.Refresh()
}

// ViewItem this will remove any existing and render new tree object
func (m *Mytree) ViewItem(name string) {
	for _, v := range m.Panel.Objects {
		m.Panel.Remove(v)
	}
	//item is not present create a new one
	if m.ItemStore[name] == nil {
		m.ItemStore[name] = NewTreeObj(name)
		fmt.Println("Created NEw Object")
	}

	x := m.ItemStore[name]

	m.Panel.Add(x.Render())

	m.Panel.Refresh()
}

// AddItem : this will add new item if parent is found else returns error
func (m *Mytree) AddItem(parent string, child string, branch bool) error {
	if parent == "" {
		m.Nodes[parent] = append(m.Nodes[parent], child)
		if branch {
			m.Nodes[child] = []string{}
		}
	} else if parent != "" {
		count := 0
		for k := range m.Nodes {
			if k == parent {
				m.Nodes[k] = append(m.Nodes[k], child)
				if branch {
					m.Nodes[child] = []string{}
				}
			} else {
				count = count + 1
			}
		}
		if len(m.Nodes) == count {
			return fmt.Errorf("parent Node Was Not Found")
		}
	}

	m.Refresh()
	return nil
}

// Remove : Deletes and returns array if found
func Remove(arr []string, item string) []string {
	for k, v := range arr {
		if v == item {
			return append(arr[:k], arr[k+1:]...)
		}
	}
	return arr
}

// DeleteItem : this will find and delete item if parent is found else returns error
func (m *Mytree) DeleteItem(parent string, child string, branch bool) error {
	if parent == "" {
		m.Nodes[parent] = Remove(m.Nodes[parent], child)
		if branch {
			delete(m.Nodes, child)
		}
	} else if parent != "" {
		count := 0
		for k := range m.Nodes {
			if k == parent {
				m.Nodes[k] = Remove(m.Nodes[k], child)
				if branch {
					delete(m.Nodes, child)
				}
			} else {
				count = count + 1
			}
		}
		if len(m.Nodes) == count {
			return fmt.Errorf("parent Node Was Not Found")
		}
	}

	m.Refresh()
	return nil
}

// FindObject ...
func (m *Mytree) FindObject(name string) (*TreeObject, error) {
	//use maps for finding easy
	lowername := strings.ToLower(name)

	branchmatch := []string{}

	//branches with no children are skipped
	for k := range m.Nodes {
		klower := strings.ToLower(k)
		if possiblerelated(klower, lowername) {
			branchmatch = append(branchmatch, k)
		}
	}

	if len(branchmatch) > 0 {
		tip := "this API is intended to only work with check list items not branches\n try refining it to checklist item"
		newer := fmt.Errorf("matched with following branches %v\n %v", branchmatch, tip)
		return &TreeObject{}, newer
	} else {
		//check for checklist item match
		itemmatch := []string{}
		for k := range m.HashMap {
			klower := strings.ToLower(k)
			if possiblerelated(klower, lowername) {
				itemmatch = append(itemmatch, k)
			}
		}

		if len(itemmatch) > 0 {
			if len(itemmatch) == 1 {
				//return the item
				if m.ItemStore[name] == nil {
					m.ItemStore[name] = NewTreeObj(name)
					fmt.Println("Created NEw Object")
				}

				return m.ItemStore[itemmatch[0]], nil
			} else {
				newerror := fmt.Errorf("matched with multiple checklist items try refining it more %v", itemmatch)
				return &TreeObject{}, newerror
			}

		} else {
			//not matched with any item
			// fmt.Printf("you searched for %v\n", lowername)
			// fmt.Println(m.HashMap)
			case1 := "Branches with No children are not tracked if that's the case fix it"
			newerror := fmt.Errorf("not matched with any item use /checklist to list all items %v", case1)
			return &TreeObject{}, newerror

		}

	}

}

// SelectChecklist : Creates a dialog box and select checklist and loads it
func (m *Mytree) SelectChecklist() {
	found := GetAllCheckLists()
	if len(found) == 0 {
		x := dialog.NewInformation("Error", "No Global CheckList Found \nTry Exporting one!!\n", Win)
		x.Show()
		return
	} else {
		defined := mylayouts.BorderBoxMin{
			TopIntend:   8,
			LeftIntend:  8,
			RightIntend: 8,
			SpacerSize:  2,
			Padding:     4,
		}

		label := widget.NewLabelWithStyle("Available checkLists", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		spacer := layout.NewSpacer()

		//create list of all results
		list := widget.NewList(
			//length function
			func() int {
				return len(found)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(found[i].Name)
			})

		list.OnSelected = func(id widget.ListItemID) {
			m.ImportCheckList(found[id])
		}

		content := container.New(&defined, label, spacer, container.NewVScroll(list))

		d := dialog.NewCustom("", "Done", content, Win)
		d.Show()
	}
}

// GetContainer ...
func (t *Mytree) GetContainer(w fyne.Window) *fyne.Container {
	defined := mylayouts.BorderBoxMin{
		TopIntend:   8,
		LeftIntend:  8,
		RightIntend: 8,
		SpacerSize:  2,
		Padding:     4,
	}
	dat1 := binding.NewString()
	dat2 := binding.NewString()
	MakeBranch := false
	l1 := widget.NewLabelWithStyle("Parent Branch Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	l2 := widget.NewLabelWithStyle("New Item Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	e1 := widget.NewEntryWithData(dat1)
	e2 := widget.NewEntryWithData(dat2)

	opts := []string{"Branch", "Item"}
	radio := widget.NewRadioGroup(opts, func(s string) {
		if s == "Branch" {
			MakeBranch = true
		} else if s == "Item" {
			MakeBranch = false
		}
	})
	radio.Horizontal = true
	l3 := widget.NewLabelWithStyle("Item type : ", fyne.TextAlignLeading, fyne.TextStyle{})
	radiobox := container.NewHBox(l3, radio)

	content := container.New(&defined, l1, e1, layout.NewSpacer(), l2, e2, layout.NewSpacer(), radiobox, layout.NewSpacer())

	addfunc := func() {
		callback := func(b bool) {
			if b {
				//find the parent
				parent, _ := dat1.Get()
				child, _ := dat2.Get()

				t.AddItem(parent, child, MakeBranch)
				fmt.Println("Done")
			}
		}
		// fmt.Println("Button tapped")

		z := dialog.NewCustomConfirm("Add New Item", "Yes", "No", content, callback, w)
		z.Show()
		// fmt.Println("Dialog called")

	}

	//Delete Item From tree

	dropfunc := func() {
		callback2 := func(b bool) {
			if b {
				//find the parent
				parent, _ := dat1.Get()
				child, _ := dat2.Get()

				t.DeleteItem(parent, child, MakeBranch)
				fmt.Println("Done")
			}
		}

		z2 := dialog.NewCustomConfirm("Delete Item", "Yes", "No", content, callback2, w)
		z2.Show()

	}
	save := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		t.Store()
	})

	// export := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {
	// 	t.ExportAs()
	// })

	// importc := widget.NewButtonWithIcon("", theme.ContentRedoIcon(), func() {
	// 	t.SelectChecklist()
	// })

	menitem1 := fyne.NewMenuItem("Import", func() {
		t.SelectChecklist()
	})

	menuitem2 := fyne.NewMenuItem("Export", func() {
		t.ExportAs()
	})

	men := fyne.NewMenu("", menitem1, menuitem2)

	ctxmen := mywidgets.NewContextMenu(theme.MenuDropDownIcon(), men)

	label := widget.NewLabelWithStyle(t.TreeIdentifer, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	top := container.NewHBox(label, layout.NewSpacer(), widget.NewButtonWithIcon("", theme.ContentAddIcon(), addfunc), widget.NewButtonWithIcon("", theme.ContentRemoveIcon(), dropfunc), save, ctxmen)

	t.Bar = widget.NewProgressBarWithData(t.StatusBar)
	t.Bar.Max = 1
	t.Bar.Min = 0

	return container.NewBorder(top, t.Bar, nil, nil, container.NewVScroll(t))
}

// NewMytree ...
func NewMytree(identifier string) *Mytree {

	t := Mytree{}
	t.HashMap = map[string]string{}
	t.Nodes = map[string][]string{}
	t.StatusBar = binding.NewFloat()
	t.CheckBoxes = make(map[string]binding.Bool)
	t.Panel = container.NewMax()
	t.TreeIdentifer = identifier
	t.ItemStore = map[string]*TreeObject{}
	t.BCheckBoxes = map[string]bool{}
	t.Coll = "checklists"

	baseview := container.NewCenter(widget.NewLabel("This is Default View"))
	t.Panel.Add(baseview)

	t.ChildUIDs = func(uid widget.TreeNodeID) (c []widget.TreeNodeID) {
		//returns an array of child nodes
		return t.Nodes[uid]
	}

	t.IsBranch = func(uid widget.TreeNodeID) (ok bool) {
		//check if it has a branch
		children, ok := t.Nodes[uid]

		// also if branch does not have any children
		return ok && len(children) >= 0
	}

	t.CreateNode = func(branch bool) (o fyne.CanvasObject) {
		cwidget := widget.NewCheck("", func(b bool) {})
		label := widget.NewLabelWithStyle("Template", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		return container.NewBorder(nil, nil, nil, cwidget, label)
	}

	t.UpdateNode = func(uid widget.TreeNodeID, branch bool, node fyne.CanvasObject) {
		//set item name to label
		objs := node.(*fyne.Container).Objects
		for _, v := range objs {
			if w, ok := v.(*widget.Label); ok {
				//change label value to uid
				w.SetText(uid)
			}
			if w, ok := v.(*widget.Check); ok {
				//change check status and tap function
				//chack if map has an object
				if _, ok := t.CheckBoxes[uid]; !ok {
					t.CheckBoxes[uid] = binding.NewBool()
				}
				//bind to node no need for fuzz
				w.Bind(t.CheckBoxes[uid])

				if branch {
					w.OnChanged = func(b bool) {
						t.UpdateCheckStatus(uid, b)
						t.UpdateBar()
					}
				} else {
					w.OnChanged = func(b bool) {
						t.CheckBoxes[uid].Set(b)
						t.UpdateBar()
					}
				}
			}
		}

		//use branches for reverse indexing childs
		if branch {
			childs := t.Nodes[uid]
			for _, v := range childs {
				t.HashMap[v] = uid
			}
		}
	}

	t.OnSelected = func(uid widget.TreeNodeID) {
		t.ViewItem(uid)
	}

	t.ExtendBaseWidget(&t)

	if identifier == "orgdefault" {
		if OrgCheck == nil {
			LoadOrgCData()
		}

		// fmt.Println(OrgDefault)
		t.ImportCheckList(*OrgCheck)
		return &t
	}

	//Get default Checklist if availabel and use it
	if Default == nil {
		LoadDefaultCheckList()
	}

	// fmt.Println("Not org checklist")
	// fmt.Println(Default)

	//Import Default Checklist
	t.ImportCheckList(*Default)

	return &t

}

// LoadOrgCData :  Only invoke if type of checklist is Organization
func LoadOrgCData() {
	cls := GetAllCheckLists()
	found := false

	for _, v := range cls {
		if v.Name == "orgdefault" {
			found = true
			OrgCheck = &v
		}
	}

	if !found {
		var tsemplate = map[string][]string{
			"":                      {"Subdomain Enumeration", "Github OSINT", "S3 Buckets"},
			"Subdomain Enumeration": {"findomain", "amass", "sublist3r"},
			"S3 Buckets":            {"GrayHatWarfare"},
		}

		OrgCheck = &CheckListOnly{
			Name:         "orgdefault",
			TreeNodes:    tsemplate,
			Descriptions: map[string]*TreeObject{},
		}

		fmt.Println("Something Went Wrong Could not load orgdefault checklist")
		fmt.Println("Loading a Simple Template")

	}
}

// LoadDefaultCheckList ...
func LoadDefaultCheckList() {
	cls := GetAllCheckLists()
	found := false

	for _, v := range cls {
		// fmt.Println(v.Name)
		if v.Name == "default" {
			// fmt.Printf("Condition v.name== default Met found %v\n", v.Name)
			found = true
			Default = &v
			fmt.Printf("Done Loading %v in Default \n", Default.Name)
			break
			// fmt.Println("Done loading default")
			// fmt.Println(v)
		}
	}

	if !found {
		var tssemplate = map[string][]string{
			"":                      {"Information Gathering", "Scanning", "Enumeration"},
			"Information Gathering": {"Nmap", "Whois", "Whatweb"},
			"Scanning":              {"Nikito", "JSFinder"},
		}

		Default = &CheckListOnly{
			Name:         "default",
			TreeNodes:    tssemplate,
			Descriptions: map[string]*TreeObject{},
		}

		fmt.Println("Something Went Wrong Could not load default checklist")
		fmt.Println("Loading a Simple Template")

	}
}

func possiblerelated(s1 string, s2 string) bool {
	if s1 == s2 || strings.Contains(s1, s2) {
		return true
	}

	return false
}
