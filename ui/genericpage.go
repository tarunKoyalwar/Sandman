package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/hugetext"
	"github.com/tarunKoyalwar/sandman/mylayouts"
	"github.com/tarunKoyalwar/sandman/mywidgets"
	"go.mongodb.org/mongo-driver/bson"
)

// Page ...
type Page struct {
	HeaderText     string              `bson:"headername"`
	Heading        *canvas.Text        `bson:"-"`
	Entry          *hugetext.HugeEntry `bson:"-"`
	CollectionName string              `bson:"-"`
	Entrydata      []string            `bson:"data"`
	BtnName        string              `bson:"-"`
	PanelName      string              `bson:"panelname"`
	ShowLength     bool                `bson:"-"`
}

// Store ...
func (P *Page) Store() {
	db.MDB.GetCollection(P.CollectionName)
	P.Entrydata = P.Entry.GetData()
	// db.MDB.InsertDocument()
	//filter used to match document
	filter := bson.M{"panelname": P.PanelName}
	res, err := db.MDB.UpdateDocument(filter, bson.M{"$set": P})
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

// Update ...
func (P *Page) Update() {
	var z Page
	db.MDB.GetCollection(P.CollectionName)
	filter := bson.M{"panelname": P.PanelName}
	dat, _ := db.MDB.FindOne(filter, &z)
	// if err != nil {
	// 	fyne.LogError("No data to update", err)
	// }
	ptr := dat.(*Page)
	if ptr.HeaderText != "" {
		P.HeaderText = ptr.HeaderText
		P.Entrydata = ptr.Entrydata
		P.Entry.SetData(ptr.Entrydata)
	}
	P.Entry.Refresh()
}

// Setup ...
func (P *Page) Setup() {
	P.Heading = canvas.NewText(P.HeaderText, theme.ForegroundColor())
	P.Heading.TextSize = 20
	P.Heading.TextStyle.Bold = true
	P.Heading.Alignment = fyne.TextAlignLeading

	P.Entry = hugetext.NewHugeEntry(w)
}

// Render ...
func (P *Page) Render() *fyne.Container {
	P.Store()
	// P.Update()

	// Border render

	headingbox := container.NewBorder(nil, nil, nil, container.NewHBox(P.Entry.GetNavButtons()), P.Heading)

	P.Heading.Text = P.HeaderText

	defined := mylayouts.BorderBox{
		TopIntend:   7,
		LeftIntend:  8,
		RightIntend: 0,
		SpacerSize:  3,
	}

	if P.ShowLength {
		dat := P.Entry.GetData()
		size := 0
		for _, v := range dat {
			for _, b := range strings.Split(v, "\n") {
				if strings.TrimSpace(b) != "" {
					size += 1
				}
			}
		}
		if strings.Contains(P.Heading.Text, "(") {
			P.Heading.Text = regreplace(P.Heading.Text, "("+strconv.Itoa(size)+")")
		} else {
			P.Heading.Text += "(" + strconv.Itoa(size) + ")"
		}
	}

	cont := container.New(&defined, headingbox, widget.NewSeparator(), layout.NewSpacer(), container.NewVScroll(P.Entry))

	fmt.Println("Rendered and stored")

	return cont
}

// GetButton ...
func (P *Page) GetButton() *mywidgets.CustomButton {
	btn := mywidgets.CustomButton{
		Text:      fmt.Sprintf("\t%v", P.BtnName),
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

// NewPage If button name is empty header name is used
func NewPage(header string, coll string, btn string) *Page {
	if btn == "" {
		btn = header
	}
	p := Page{
		HeaderText:     header,
		CollectionName: coll,
	}
	p.BtnName = btn
	p.PanelName = p.BtnName

	p.Setup()
	p.Update()

	return &p
}

//replace (0-9) with something
func regreplace(old string, new string) string {
	re := regexp.MustCompile("[(][0-9]*[)]")

	return re.ReplaceAllString(old, new)
}
