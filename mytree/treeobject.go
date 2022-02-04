package mytree

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/hugetext"
	"github.com/tarunKoyalwar/sandman/mylayouts"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Win ...
var Win fyne.Window

// BackupTreeObj
type BackupTreeObj struct {
	Heading string   `bson:"heading"`
	Payload string   `bson:"payload"`
	Notes   string   `bson:"notes"`
	Toolout []string `bson:"toolout"`
}

//TreeObject
type TreeObject struct {
	HeaderText    string              `bson:"headertext"`
	Heading       *canvas.Text        `bson:"-"`
	Payload       string              `bson:"payload"`
	entry         *widget.Entry       `bson:"-"`
	Notes         string              `bson:"notes"`
	notesentry    *widget.Entry       `bson:"-"`
	ToolEntry     *hugetext.HugeEntry `bson:"-"`
	ToolEntryData []string            `bson:"toolentry"`
	// ToolPage   *Page
}

// Store : Copy required text and export
func (t *TreeObject) Store() {
	t.Payload = t.entry.Text
	t.Notes = t.notesentry.Text
	t.ToolEntryData = t.ToolEntry.GetData()
	// b := BackupTreeObj{
	// 	Heading: t.HeaderText,
	// 	Notes:   t.notesentry.Text,
	// 	Payload: t.Payload,
	// 	Toolout: t.ToolEntryData,
	// }

}

//Update Existing Values without changing original struct
func (t *TreeObject) Update() {
	t.Setup()
	t.ToolEntry.SetData(t.ToolEntryData)
	t.ToolEntry.Refresh()
}

// Setup ...
func (t *TreeObject) Setup() {
	t.Heading = canvas.NewText(t.HeaderText, theme.ForegroundColor())
	t.Heading.TextSize = 20
	t.Heading.TextStyle.Bold = true
	t.Heading.Alignment = fyne.TextAlignLeading

	t.entry = widget.NewMultiLineEntry()
	t.entry.Bind(binding.BindString(&t.Payload))
	t.entry.Wrapping = fyne.TextWrapWord
	t.notesentry = widget.NewMultiLineEntry()
	t.notesentry.Bind(binding.BindString(&t.Notes))
	t.notesentry.Wrapping = fyne.TextWrapBreak
	t.ToolEntry = hugetext.NewHugeEntry(Win)
}

// Render ...
func (t *TreeObject) Render() fyne.CanvasObject {

	t.Heading.Text = t.HeaderText

	defined := mylayouts.BorderBox{
		TopIntend:   8,
		LeftIntend:  8,
		RightIntend: 0,
		SpacerSize:  4,
	}

	label := canvas.NewText("Payload", theme.ForegroundColor())
	label.TextSize = 16
	label.TextStyle.Bold = true

	label2 := canvas.NewText("Notes", theme.ForegroundColor())
	label2.TextSize = 16
	label2.TextStyle.Bold = true

	cont := container.New(&defined, t.Heading, widget.NewSeparator(), layout.NewSpacer(), label, t.entry, layout.NewSpacer(), label2, t.notesentry)

	toolheading := canvas.NewText("Command Output", theme.ForegroundColor())
	toolheading.TextSize = 20
	toolheading.TextStyle.Bold = true
	toolheading.Alignment = fyne.TextAlignLeading

	headingbox := container.NewBorder(nil, nil, nil, container.NewHBox(t.ToolEntry.GetNavButtons()), toolheading)

	defined2 := mylayouts.BorderBox{
		TopIntend:   8,
		LeftIntend:  8,
		RightIntend: 0,
		SpacerSize:  4,
	}

	cont2 := container.New(&defined2, headingbox, widget.NewSeparator(), layout.NewSpacer(), container.NewVScroll(t.ToolEntry))

	tabs := container.NewAppTabs(container.NewTabItem(
		"Description", cont,
	),
		container.NewTabItem(
			"ToolOutput", cont2,
		),
	)

	// tabs.SelectIndex(0)
	tabs.SetTabLocation(container.TabLocationBottom)

	return tabs
}

// NewTreeObj ...
func NewTreeObj(name string) *TreeObject {
	t := TreeObject{}
	t.HeaderText = name
	t.Setup()
	return &t
}

// GetAllTrees ...
func GetAllTrees() []TreeBackup {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	/*
		cursor, err := episodesCollection.Find(ctx, bson.M{})
		if err != nil {
		    log.Fatal(err)
		}
		var episodes []bson.M
		if err = cursor.All(ctx, &episodes); err != nil {
		    log.Fatal(err)
		}

	*/
	coll := db.MDB.GetCollInstance()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var arr []TreeBackup
	if err = cursor.All(ctx, &arr); err != nil {
		panic(err)
	}

	// fmt.Println(arr)

	return arr
}

// GetAllCheckLists Gets All CheckLists From db
func GetAllCheckLists() []CheckListOnly {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//Get Temp Database instance
	//Retrieves a database with name checklist
	db := db.MDB.TempConnection("checklists")
	colllist, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	found := false
	for _, v := range colllist {
		if v == "webchecklists" {
			found = true
		}
	}
	if !found {
		return []CheckListOnly{}
	}
	coll := db.Collection("webchecklists")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	var arr []CheckListOnly
	if err = cursor.All(ctx, &arr); err != nil {
		panic(err)
	}

	// fmt.Println(arr)
	defer db.Client().Disconnect(ctx)
	// fmt.Println(arr)
	// fmt.Printf("While Loading CheckList")
	// for _, s := range arr {
	// 	fmt.Println(s.Name)
	// }
	return arr
}

// PutCheckList:  puts checklist to database
func PutCheckList(d CheckListOnly) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := db.MDB.TempConnection("checklists")
	colllist, err := db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		panic(err)
	}
	found := false
	for _, v := range colllist {
		if v == "webchecklists" {
			found = true
		}
	}
	if !found {
		fmt.Printf("Collection Not Found Creating New One")
		err := db.CreateCollection(ctx, "webchecklists")
		if err != nil {
			panic(err)
		}
	}

	coll := db.Collection("webchecklists")

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"checklistname": d.Name}
	res, err := coll.UpdateOne(ctx, filter, bson.M{"$set": d}, opts)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Matchedcount %v , Modified COunt %v, with upsert id %v\n ", result.MatchedCount, result.ModifiedCount, result.UpsertedID)
	fmt.Printf("Matched Document %v Updated Document %v\n", res.MatchedCount, res.UpsertedCount)

	defer db.Client().Disconnect(ctx)
}
