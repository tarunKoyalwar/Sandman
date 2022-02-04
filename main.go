package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/tarunKoyalwar/sandman/assets"
	"github.com/tarunKoyalwar/sandman/db"
	"github.com/tarunKoyalwar/sandman/mytree"
	"github.com/tarunKoyalwar/sandman/ui"
)

func main() {
	a := app.NewWithID("Sandman")

	//Set Icon
	a.SetIcon(assets.ResourceSandmanJpg)
	a.Settings().SetTheme(assets.MyTheme{})

	w := a.NewWindow("The Dreaming")
	w.SetPadded(false)
	w.Resize(fyne.NewSize(640, 480))

	//window pointer for tree
	mytree.Win = w

	// db setup
	m := &db.MDB

	ui.AddPointer(w, a)
	ui.UI_INIT()
	ui.AddKeyboardEvents(w)

	//Loads previously saved db string
	ui.GetPreferenceDBStringtoURL()

	w.SetContent(ui.StartupWindow())

	w.CenterOnScreen()
	w.ShowAndRun()

	if m.Isconnected() {
		fmt.Println("Saving State")
		ui.E.SaveState()
		fmt.Printf("Saved State Before Exiting")
	}

	defer func() {
		if m.Isconnected() {
			m.Disconnect()
		}

	}()
}

func HandleError(er error) {
	if er != nil {
		panic(er)
	}
}

//Identify different types of signals and fix this
func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Saving State")
		ui.E.SaveState()
		fmt.Printf("Saved State Before Exiting")
		// data.Task()
		os.Exit(1)
	}()
}
