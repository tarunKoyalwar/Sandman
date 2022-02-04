package hugetext

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

//This aims to handle huge text by paging the data
//this will reduce cpu load and improve performance
//default paging size is 200

//TODO
//implement manual word wrapping by measuring text

// HugeEntry : Huge Entry Custom Widget
type HugeEntry struct {
	widget.Entry
	Pagesize      int
	DataGrid      []binding.String
	pos           int
	pasteshortcut desktop.CustomShortcut //alt + v
	leftshortcut  desktop.CustomShortcut //alt <-
	rightshortcut desktop.CustomShortcut //alt ->
	copyshortcut  desktop.CustomShortcut //alt+c
	w             fyne.Window
	Unique        bool //Make Data Unique
}

// GetData : Returns data
func (h *HugeEntry) GetData() []string {
	d := []string{}
	for _, v := range h.DataGrid {
		x, _ := v.Get()
		d = append(d, x)
	}
	return d
}

// SetData : Set Data
func (h *HugeEntry) SetData(x []string) {
	if h.Unique {
		x = GetUnique(x)
	}
	for k, v := range x {
		h.adddataatposition(k, v)
	}
	// fmt.Println("Data Pasted to Multiple pages")
	h.changeposandupdate(h.pos)
}

func (h *HugeEntry) changeposandupdate(pos int) {
	h.Unbind()
	h.pos = pos
	if h.pos < len(h.DataGrid) && h.pos >= 0 {
		h.Bind(h.DataGrid[h.pos])
	} else {
		fmt.Printf("Cusor is greater than length fixing it\n")
		h.pos = len(h.DataGrid) - 1
		h.changeposandupdate(h.pos)
	}
	h.SetPlaceHolder(fmt.Sprintf("On Page : %v", h.pos))
	h.Refresh()
}

func (h *HugeEntry) adddataatposition(pos int, dat string) {
	if pos < len(h.DataGrid) {
		h.DataGrid[pos].Set(dat)
		//create next entry in datagrid if this is last one
		if pos == len(h.DataGrid)-1 {
			h.DataGrid = append(h.DataGrid, binding.NewString())
		}
	} else {
		fmt.Printf("This was not supposed to happen paste position is greater than length")
	}
}

// Set : This is used to set data using string
func (h *HugeEntry) Set(dat string, append bool) {
	// fmt.Printf("[debug] received %v and %v\n", dat, append)
	pos := 0
	alldata := ""
	if append {
		arr := h.GetData()
		// fmt.Printf("[debug]all data in entry '%v'\n", arr)
		temp := strings.TrimSpace(StringArrtoString(arr))
		alldata = temp + "\n" + strings.TrimSpace(dat)
		// fmt.Printf("[debug] all data after adding strings %v\n", alldata)
	} else {
		alldata = strings.TrimSpace(dat)
		// fmt.Printf("[debug] called not append\n")
	}
	holder := ""
	splitdata := String2Arr(strings.TrimSpace(alldata))
	if h.Unique {
		splitdata = GetUnique(splitdata)
	}
	// fmt.Printf("[debug] sp")
	// fmt.Printf("[debug] huge set splitdata |%v| alldata |%v|\n", splitdata, alldata)
	for k, v := range splitdata {
		if k == 0 {
			holder = v + "\n"
			if k == len(splitdata)-1 {
				h.adddataatposition(pos, holder)
			}
			continue
		}

		holder = holder + v + "\n"

		if k == len(splitdata)-1 {
			//last item so save
			h.adddataatposition(pos, holder)
			holder = ""
		} else if k%h.Pagesize == 0 {
			//and k is not equal to 0
			h.adddataatposition(pos, holder)
			holder = ""
			pos += 1
		}
	}
	// h.adddataatposition(pos, holder)
	// fmt.Println("Data Pasted to Multiple pages")
	h.changeposandupdate(h.pos)
}

//THis will handle UI paste
func (h *HugeEntry) paste() {
	pos := h.pos
	dat := h.w.Clipboard().Content()
	current, _ := h.DataGrid[h.pos].Get()
	all := current + dat
	holder := ""
	splitdata := strings.Split(all, "\n")
	if h.Unique {
		splitdata = GetUnique(splitdata)
	}
	for k, v := range splitdata {
		if k == 0 {
			holder = v + "\n"
			if k == len(splitdata)-1 {
				h.adddataatposition(pos, holder)
			}
			continue
		}

		holder = holder + v + "\n"

		if k == len(splitdata)-1 {
			//last item so save
			h.adddataatposition(pos, holder)
			holder = ""
		} else if k%h.Pagesize == 0 {
			//and k is not equal to 0
			h.adddataatposition(pos, holder)
			holder = ""
			pos += 1
		}
	}
	// fmt.Println("Data Pasted to Multiple pages")
	h.changeposandupdate(h.pos)
}

// GetNavButtons : Left and Right Buttons
func (h *HugeEntry) GetNavButtons() (*widget.Button, *widget.Button) {
	left := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		if h.pos == 0 {
			fmt.Printf("Already at least position\n")
			return
		} else {
			h.changeposandupdate(h.pos - 1)
		}
	})

	right := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		//change position ot right
		if h.pos == len(h.DataGrid)-1 {
			fmt.Printf("Created New Page %v\n", h.pos+1)
			h.DataGrid = append(h.DataGrid, binding.NewString())
		}

		if h.pos < len(h.DataGrid) {
			h.changeposandupdate(h.pos + 1)
		}
	})

	return left, right
}

// TypedShortcut : Custom Shortcut to paste into multiple pages
func (h *HugeEntry) TypedShortcut(s fyne.Shortcut) {

	if h.pasteshortcut.ShortcutName() == s.ShortcutName() {
		fmt.Println("GOt shift + v")
		h.paste()
	} else if h.leftshortcut.ShortcutName() == s.ShortcutName() {
		if h.pos == 0 {
			fmt.Printf("Already at least position\n")
			return
		} else {
			h.changeposandupdate(h.pos - 1)
		}

	} else if h.rightshortcut.ShortcutName() == s.ShortcutName() {
		//change position ot right
		if h.pos == len(h.DataGrid)-1 {
			fmt.Printf("Created New Page %v\n", h.pos+1)
			h.DataGrid = append(h.DataGrid, binding.NewString())
		}

		if h.pos < len(h.DataGrid) {
			h.changeposandupdate(h.pos + 1)
		}
	} else if h.copyshortcut.ShortcutName() == s.ShortcutName() {
		//copy all data to clipboard
		dat := h.GetData()
		joined := ""
		for _, v := range dat {
			joined = joined + v
		}
		h.w.Clipboard().SetContent(joined)
		fmt.Printf("Copied all pages to clipboard")

	} else {
		//send others to global
		h.Entry.TypedShortcut(s)
	}

	// log.Println("Shortcut typed:", s.ShortcutName())
}

//NewHugeEntry : Create New HugeEntry Widget
func NewHugeEntry(win fyne.Window) *HugeEntry {
	h := HugeEntry{}
	h.MultiLine = true

	h.Pagesize = 200
	h.pos = 0
	h.w = win
	h.Wrapping = fyne.TextWrapBreak
	// h.Clipboard = clipboard
	h.SetPlaceHolder(fmt.Sprintf("On Page : %v", h.pos))
	h.DataGrid = []binding.String{
		binding.NewString(), binding.NewString(),
	}

	h.pasteshortcut = desktop.CustomShortcut{}
	h.pasteshortcut.KeyName = fyne.KeyV
	h.pasteshortcut.Modifier = desktop.AltModifier

	h.leftshortcut = desktop.CustomShortcut{}
	h.leftshortcut.KeyName = fyne.KeyLeft
	h.leftshortcut.Modifier = desktop.AltModifier

	h.rightshortcut = desktop.CustomShortcut{}
	h.rightshortcut.KeyName = fyne.KeyRight
	h.rightshortcut.Modifier = desktop.AltModifier

	h.copyshortcut = desktop.CustomShortcut{}
	h.copyshortcut.KeyName = fyne.KeyC
	h.copyshortcut.Modifier = desktop.AltModifier

	// fmt.Println(h.DataGrid[0])
	h.Bind(h.DataGrid[h.pos])
	h.ExtendBaseWidget(&h)
	// h.TypedShortcut()

	return &h
}

// StringArrtoString : Convert arr to string by joining using space
func StringArrtoString(dat []string) string {
	all := ""
	for _, v := range dat {
		if strings.TrimSpace(v) != "" {
			all += v
		}

	}
	return all
}

// String2Arr : Split string to arr using `\n`
func String2Arr(str string) []string {
	final := []string{}

	for _, v := range strings.Split(str, "\n") {
		if strings.TrimSpace(v) != "" {
			final = append(final, strings.TrimSpace(v))
		}
	}

	// fmt.Printf("Got %v Sent %v\n", str, final)

	return final
}

// GetUnique : Return Unique Elements from Array
func GetUnique(arr []string) []string {
	dat := map[string]bool{}

	for _, v := range arr {
		dat[strings.TrimSpace(v)] = true
	}

	keys := make([]string, 0, len(dat))
	for k := range dat {
		keys = append(keys, k)
	}

	return keys
}
