package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// AddKeyboardEvents : Adds Shortcuts
func AddKeyboardEvents(w fyne.Window) {
	//ctrl +s
	printdat := desktop.CustomShortcut{
		KeyName:  fyne.KeyS,
		Modifier: desktop.ControlModifier,
	}
	w.Canvas().AddShortcut(&printdat, func(shortcut fyne.Shortcut) {
		fmt.Printf("Saving Data to DB Called By %v\n", shortcut.ShortcutName())
		E.SaveState()
	})

	//ctrl + p | to print resources

}
