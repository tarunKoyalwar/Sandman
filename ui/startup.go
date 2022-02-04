package ui

import (
	"fmt"
	"net"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mylayouts"
)

var w fyne.Window

var a fyne.App

// AddPointer ...
func AddPointer(wx fyne.Window, ap fyne.App) {
	w = wx
	a = ap
}

// StartupWindow ...
func StartupWindow() fyne.CanvasObject {

	base := container.NewGridWithRows(3,
		layout.NewSpacer(),
		container.NewGridWithColumns(3,
			layout.NewSpacer(),
			container.NewGridWithColumns(1,
				widget.NewButtonWithIcon("New Project", theme.FolderNewIcon(), newproject),
				widget.NewButtonWithIcon("Open Project", theme.FolderOpenIcon(), openproject),
				widget.NewButtonWithIcon("Temp Project", theme.FileIcon(), tempproject),
				widget.NewButtonWithIcon("Configure MongoDB", theme.WarningIcon(), configuremongo),
				widget.NewButtonWithIcon("Drop Projects", theme.DeleteIcon(), deleteprojects),
				widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), settings),
			),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
	)

	return base
}

//must deal with creating a new database and initialize project
//and handle error
func newproject() {
	er := db.MDB.Connect()
	if er != nil {
		FyneError(er)
		return
	}

	err2 := db.MDB.PingTest()
	FyneError(err2)
	if err2 != nil {
		return
	}

	Heading := canvas.NewText("Project Name", theme.ForegroundColor())
	Heading.TextSize = 20
	Heading.TextStyle.Bold = true
	Heading.Alignment = fyne.TextAlignLeading

	mystr := binding.NewString()

	entry := widget.NewEntryWithData(mystr)
	entry.SetPlaceHolder("Enter Project Name")
	done := widget.NewButton("Done", func() {
		fmt.Printf("Creating Database and Initializing the project\n")
		dat, err := mystr.Get()
		if err != nil {
			dialog.NewError(err, w)
		}
		db.MDB.GetDatabase(dat)
		err = db.MDB.InitializeProject()
		FyneError(err)
		fmt.Printf("...Done\n")
		fmt.Printf("Will Redirect to HomeScreen \n")
		w.SetContent(HomeScreen())
	})

	definedlayout := &mylayouts.BorderBox{
		TopIntend:   7,
		LeftIntend:  5,
		RightIntend: 5,
		SpacerSize:  3,
	}

	headingbox := container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		w.SetContent(StartupWindow())
	}), Heading)

	final := container.New(definedlayout, headingbox, widget.NewSeparator(), layout.NewSpacer(), entry, layout.NewSpacer(), container.NewHBox(layout.NewSpacer(), done))

	w.SetContent(final)

}

//lists project and then move to home screen
//open projects will move to navigation panel and home screen
func openproject() {
	er := db.MDB.Connect()
	if er != nil {
		FyneError(er)
		return
	}

	err2 := db.MDB.PingTest()
	FyneError(err2)
	if err2 != nil {
		return
	}
	names, err := db.MDB.ListDatabases()
	FyneError(err)
	if err != nil {
		return
	}
	filtered := []string{}
	for _, v := range names {
		if v == "admin" || v == "local" || v == "config" || v == "checklists" {
			continue
		} else {
			filtered = append(filtered, v)
		}
	}
	if len(filtered) < 1 {
		dialog.NewInformation("No Projects", "Looks Like You have Not Started Any Projects. Redirecting to Start Page", w)
		w.SetContent(StartupWindow())
		return
	}

	fmt.Printf("Found Projects %v\n", filtered)

	Heading := canvas.NewText("Projects Found", theme.ForegroundColor())
	Heading.TextSize = 20
	Heading.TextStyle.Bold = true
	Heading.Alignment = fyne.TextAlignLeading
	spacer := layout.NewSpacer()

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
		fmt.Printf("Calling %v\n", filtered[id])
		db.MDB.GetDatabase(filtered[id])
		if !db.MDB.ValidateProject() {
			db.MDB.InitializeProject()
			fmt.Println("Initialized the project...")

		}

		w.SetContent(HomeScreen())
		//setup new window
		//and load database data
	}

	definedlayout := &mylayouts.BorderBox{
		TopIntend:   7,
		LeftIntend:  5,
		RightIntend: 5,
		SpacerSize:  3,
	}

	headingbox := container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		w.SetContent(StartupWindow())
	}), Heading)

	// withscroll := container.NewVScroll(container.NewGridWithColumns(1, objects...))

	// grid := container.NewGridWithColumns(1)

	final := container.New(definedlayout, headingbox, widget.NewSeparator(), spacer, list)

	w.SetContent(final)
}

//this will delete projects
func deleteprojects() {
	er := db.MDB.Connect()
	if er != nil {
		FyneError(er)
		return
	}

	err2 := db.MDB.PingTest()
	FyneError(err2)
	if err2 != nil {
		return
	}
	names, err := db.MDB.ListDatabases()
	FyneError(err)
	filtered := []string{}
	for _, v := range names {
		if v == "admin" || v == "local" || v == "config" || v == "checklists" {
			continue
		} else {
			filtered = append(filtered, v)
		}
	}
	if len(filtered) < 1 {
		dialog.NewInformation("No Projects", "Looks Like You have Not Started Any Projects. Redirecting to Start Page", w)
		w.SetContent(StartupWindow())
		return
	}

	fmt.Printf("Found Projects %v\n", filtered)

	label := widget.NewLabelWithStyle("Projects Found", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	spacer := layout.NewSpacer()

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
		fmt.Printf("Deleting %v\n", filtered[id])
		dialog.ShowConfirm(fmt.Sprintf("Deleting %v\n", filtered[id]), "Are You Sure ??", func(b bool) {
			if b {
				db.MDB.DropDatabase(filtered[id])
				w.SetContent(StartupWindow())
			}
		}, w)

		//setup new window
		//and load database data
	}

	definedlayout := &mylayouts.BorderBox{
		TopIntend:   6,
		LeftIntend:  3,
		RightIntend: 3,
		SpacerSize:  3,
	}

	// withscroll := container.NewVScroll(container.NewGridWithColumns(1, objects...))

	// grid := container.NewGridWithColumns(1)

	final := container.New(definedlayout, label, spacer, list)

	w.SetContent(final)
}

//temp database
//will save it though just drop if called again
func tempproject() {
	er := db.MDB.Connect()
	if er != nil {
		FyneError(er)
		return
	}

	err2 := db.MDB.PingTest()
	FyneError(err2)
	if err2 != nil {
		return
	}
	dbs, err := db.MDB.ListDatabases()
	if err != nil {
		fmt.Printf("List database errror : %v\n", err)
		FyneError(err)
		return
	}
	for _, v := range dbs {
		if v == "tempy" {
			err = db.MDB.DropDatabase("tempy")
			if err != nil {
				dialog.NewError(err, w)
				return
			}
			fmt.Printf("Successfully Deleted Project %v\n", v)
		}
	}
	db.MDB.GetDatabase("tempy")
	err = db.MDB.InitializeProject()
	if err != nil {
		dialog.NewError(err, w)
		return
	}
	fmt.Printf("...Done\n")
	fmt.Printf("Will Redirect to HomeScreen \n")
	fmt.Printf("Note Temp Will be Dropped everytime you use this app")
	// time.Sleep(time.Duration(3) * time.Second)
	w.SetContent(HomeScreen())
}

//

//config mongodatabase only here
//will do this later not now
func configuremongo() {
	Heading := canvas.NewText("MongoDB Settings", theme.ForegroundColor())
	Heading.TextSize = 20
	Heading.TextStyle.Bold = true
	Heading.Alignment = fyne.TextAlignLeading

	// spacer := layout.NewSpacer()
	respdata := widget.NewMultiLineEntry()
	respdata.SetPlaceHolder("Response From Server")
	respdata.Wrapping = fyne.TextWrapBreak

	label := canvas.NewText("MongoDB Connection String", theme.ForegroundColor())
	label.TextSize = 16
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignLeading

	entry := widget.NewMultiLineEntry()
	entry.SetText(GetPreferenceDBString())
	entry.SetPlaceHolder("mongodb://localhost:27017")
	entry.Wrapping = fyne.TextWrapBreak

	entry.OnChanged = func(s string) {
		db.MDB.URL = s
	}

	//Connect With Context
	connectbtn := widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
		err := db.MDB.Connect()
		fmt.Printf("Connection Error Log %v\n", err)
		if err != nil {
			respdata.SetText(err.Error())
		} else {
			respdata.SetText("Connected to Database ")
			SavePreferenceDBString(entry.Text)
			//save the string
		}
	})

	//Home Button
	homebtn := widget.NewButtonWithIcon("Home", theme.HomeIcon(), func() {
		w.SetContent(StartupWindow())
	})

	definedlayout := &mylayouts.BorderBoxMin{
		TopIntend:   7,
		LeftIntend:  6,
		RightIntend: 6,
		SpacerSize:  3,
	}

	resplabel := canvas.NewText("MongoDB Response", theme.ForegroundColor())
	resplabel.TextSize = 16
	resplabel.TextStyle.Bold = true
	resplabel.Alignment = fyne.TextAlignLeading

	final := container.New(definedlayout, Heading, widget.NewSeparator(), layout.NewSpacer(), label, entry, layout.NewSpacer(), container.NewHBox(homebtn, layout.NewSpacer(), connectbtn), layout.NewSpacer(), resplabel, respdata)

	w.SetContent(final)
}

//settings page
func settings() {
	Heading := canvas.NewText("Settings Page", theme.ForegroundColor())
	Heading.TextSize = 20
	Heading.TextStyle.Bold = true
	Heading.Alignment = fyne.TextAlignLeading

	savelabel := canvas.NewText("Save to DB when page is changed", theme.ForegroundColor())
	savelabel.TextSize = 16
	savelabel.TextStyle.Bold = true
	savelabel.Alignment = fyne.TextAlignLeading
	periodicsave := widget.NewRadioGroup([]string{"on", "off"}, func(s string) {
		if s == "on" && E != nil {
			E.PeriodicSave = true
		} else if s == "off" && E != nil {
			E.PeriodicSave = false
		}
	})

	label := canvas.NewText("REST API", theme.ForegroundColor())
	label.TextSize = 16
	label.TextStyle.Bold = true
	label.Alignment = fyne.TextAlignLeading

	addr := widget.NewEntryWithData(InterfaceIp)
	addr.SetText("")
	addr.SetPlaceHolder("127.0.0.1")

	port := widget.NewEntryWithData(binding.IntToString(Port))
	val := 8088
	for {
		stat, _ := Check(val)
		if stat {
			break
		} else {
			val += 1
		}
	}
	port.SetText(strconv.Itoa(val))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Address", Widget: addr, HintText: "Connections received to only this interface are allowed"},
			{Text: "Port", Widget: port, HintText: "Port to Use"},
		},
		OnCancel: func() {
			fmt.Println("Web Server Stopped")
			go StopServer()
		},
		OnSubmit: func() {
			fmt.Println("Web Server Started")
			go StartServer()
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Web Server Started",
				Content: addr.Text + ":" + port.Text,
			})
		},
	}

	headingbox := container.NewBorder(nil, nil, nil, widget.NewButtonWithIcon("", theme.HomeIcon(), func() {
		w.SetContent(StartupWindow())
	}), Heading)

	defined := mylayouts.BorderBox{
		TopIntend:   7,
		LeftIntend:  8,
		RightIntend: 0,
		SpacerSize:  3,
	}

	final := container.New(&defined, headingbox, widget.NewSeparator(), layout.NewSpacer(), savelabel, periodicsave, layout.NewSpacer(), label, form)

	w.SetContent(final)

}

//
func FyneError(err error) {
	if err != nil {

		label := canvas.NewText("MongoDB Error Description", theme.ForegroundColor())
		label.TextSize = 16
		label.TextStyle.Bold = true
		label.Alignment = fyne.TextAlignLeading

		defined := mylayouts.BorderBoxMin{
			TopIntend:   8,
			LeftIntend:  8,
			RightIntend: 8,
			SpacerSize:  2,
			Padding:     4,
		}

		dat := widget.NewMultiLineEntry()
		dat.SetText(err.Error())
		dat.Wrapping = fyne.TextWrapBreak
		// dat.TextStyle = fyne.TextStyle{Italic: true, Monospace: true}

		dialog.ShowCustom("Error Occurred", "Ok", container.New(&defined, label, layout.NewSpacer(), dat), w)

	}
}

// GetPreferenceDBString ...
func GetPreferenceDBString() string {
	res := a.Preferences().StringWithFallback("dbstring", "mongodb://localhost:27017")
	return res
}

// SavePreferenceDBString ...
func SavePreferenceDBString(db string) {
	a.Preferences().SetString("dbstring", db)
}

// GetPreferenceDBStringtoURL ...
func GetPreferenceDBStringtoURL() {
	db.MDB.URL = GetPreferenceDBString()
}

// Check ...
func Check(port int) (status bool, err error) {

	// Concatenate a colon and the port
	host := ":" + strconv.Itoa(port)

	// Try to create a server with the port
	server, err := net.Listen("tcp", host)

	// if it fails then the port is likely taken
	if err != nil {
		return false, err
	}

	// close the server
	server.Close()

	// we successfully used and closed the port
	// so it's now available to be used again
	return true, nil

}
