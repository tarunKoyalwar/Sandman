package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mylayouts"
	"github.com/tarunKoyalwar/sandman/mywidgets"
	"go.mongodb.org/mongo-driver/bson"
)

// CredPage ...
type CredPage struct {
	HeaderText     string          `bson:"headername"`
	Heading        *canvas.Text    `bson:"-"`
	Rows           []*CredRow      `bson:"rows"`
	Holder         *fyne.Container `bson:"-"`
	CollectionName string          `bson:"-"`
	PanelName      string          `bson:"panelname"`
}

// Store : Save to DB
func (P *CredPage) Store() {
	db.MDB.GetCollection(P.CollectionName)

	for _, v := range P.Rows {
		v.Store()
	}
	// P.Entrydata = P.Entry.GetData()
	// db.MDB.InsertDocument()
	//filter used to match document
	filter := bson.M{"panelname": P.PanelName}
	res, err := db.MDB.UpdateDocument(filter, bson.M{"$set": P})
	// fmt.Println(P.Rows)
	// fmt.Printf("Saved Document %v, %v\n", res, err)
	if err != nil {
		FyneError(err)
	}
	if res.MatchedCount == 0 {
		fmt.Printf("Document Not Found Creating New One\n")
	} else {
		fmt.Printf("Updated Existing Document\n")
	}
}

// Update  : Load From Database
func (P *CredPage) Update() {
	var z CredPage
	db.MDB.GetCollection(P.CollectionName)
	filter := bson.M{"panelname": P.PanelName}
	dat, _ := db.MDB.FindOne(filter, &z)
	// if err != nil {
	// 	// FyneLogError(err)
	// 	fyne.LogError("No data found", err)
	// }
	ptr := dat.(*CredPage)
	if ptr.HeaderText != "" {
		for _, v := range ptr.Rows {
			if v.Username != "" && v.Password != "" && v.Description != "" {
				r1 := NewCredRow()
				r1.Username = v.Username
				r1.Password = v.Password
				r1.Description = v.Description
				r1.Update()
				P.Rows = append(P.Rows, &r1)
				P.Holder.Add(r1.Render())
			}
		}
	}

	P.Holder.Refresh()
}

// Setup ...
func (P *CredPage) Setup() {
	P.Heading = canvas.NewText(P.HeaderText, theme.ForegroundColor())
	P.Heading.TextSize = 20
	P.Heading.TextStyle.Bold = true
	P.Heading.Alignment = fyne.TextAlignLeading
}

// Render ...
func (P *CredPage) Render() *fyne.Container {
	// P.Update()
	P.Heading.Text = P.HeaderText

	defined := mylayouts.BorderBox{
		TopIntend:   7,
		LeftIntend:  8,
		RightIntend: 0,
		SpacerSize:  3,
	}

	tapped := func() {
		r1 := NewCredRow()
		P.Rows = append(P.Rows, &r1)
		P.Holder.Add(r1.Render())
		fmt.Println("Added New Row")
		P.Holder.Refresh()
	}

	labelrow := container.NewGridWithRows(1,
		container.NewGridWithColumns(3,
			widget.NewLabelWithStyle("Username", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("Description", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
	)

	btn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), tapped)

	cont := container.New(&defined, P.Heading, container.NewHBox(btn), labelrow, P.Holder)

	P.Store()
	// P.Store()

	return cont

}

// GetButton ...
func (P *CredPage) GetButton() *mywidgets.CustomButton {
	btn := mywidgets.CustomButton{
		Text:      fmt.Sprintf("\t%v", P.PanelName),
		TextSize:  15,
		Alignment: fyne.TextAlignLeading,
		// FillColor: color.RGBA{R: 49, G: 48, B: 66},
		OnTapped: func() {
			E.Show(P.Render())
		},
	}

	// btn.ExtendBaseWidget(btn)

	return &btn
}

// NewCredPage ...
func NewCredPage() *CredPage {
	c := CredPage{}
	c.HeaderText = "API Keys & Credentials"
	c.Holder = container.NewGridWithColumns(1)
	c.Rows = []*CredRow{}
	//first update if no rows are present then add one
	// r1 := CredRow{}
	// c.Rows = append(c.Rows, r1)
	// z := CreateNewRow(r1)
	// c.Holder.Add(z)
	c.PanelName = "Keys & Credentials"
	c.CollectionName = "notes"
	c.Setup()
	c.Update()

	if len(c.Holder.Objects) == 0 {
		r1 := NewCredRow()
		c.Rows = append(c.Rows, &r1)
		c.Holder.Add(r1.Render())
	}
	// fmt.Println("Added New Row")

	return &c

}

// CredRow ...
type CredRow struct {
	Username    string                 `bson:"username"`
	Password    string                 `bson:"password"`
	Description string                 `bson:"description"`
	username    binding.ExternalString `bson:"-"`
	password    binding.ExternalString `bson:"-"`
	description binding.ExternalString `bson:"-"`
}

// Store ...
func (c *CredRow) Store() {
	c.Username, _ = c.username.Get()
	c.Password, _ = c.password.Get()
	c.Description, _ = c.description.Get()

}

// NewCredRow ...
func NewCredRow() CredRow {
	c := CredRow{}
	c.username = binding.BindString(&c.Username)
	c.password = binding.BindString(&c.Password)
	c.description = binding.BindString(&c.Description)

	return c

}

// Update : Updates Data to Widgets
func (c *CredRow) Update() {
	c.username.Set(c.Username)
	c.password.Set(c.Password)
	c.description.Set(c.Description)
}

// Render ...
func (c *CredRow) Render() *fyne.Container {

	username := widget.NewEntryWithData(c.username)
	username.SetPlaceHolder("Grimm")
	username.Wrapping = fyne.TextWrapBreak
	username.MultiLine = true

	password := widget.NewPasswordEntry()
	password.Bind(c.password)
	password.SetPlaceHolder("1234567890")

	description := widget.NewEntryWithData(c.description)
	description.SetPlaceHolder("Admin Panel Password")
	description.Wrapping = fyne.TextWrapBreak
	description.MultiLine = true
	// description.ExtendBaseWidget()

	row := container.NewGridWithRows(1,
		container.NewGridWithColumns(3, username, password, description),
	)

	return row
}
