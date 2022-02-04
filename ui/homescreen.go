package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// HomeScreen ...
func HomeScreen() fyne.CanvasObject {

	//create pages and create panel objects with them

	inscopepage := NewPage("In Scope Domains", "global", "In Scope")
	outofscopep := NewPage("Out of Scope Domains", "global", "Out of Scope")
	invalidvulns := NewPage("Out of Scope Vulnerabilities", "global", "Not Accepted Vuls")

	summaryobj := PanelObject{
		Header: "Program Summary",
	}

	sumlabel := widget.NewLabelWithStyle("Program Summary", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	summaryobj.Objects = []fyne.CanvasObject{
		sumlabel,
		inscopepage.GetButton(),
		outofscopep.GetButton(),
		invalidvulns.GetButton(),
	}

	activedomains := NewPage("Active Subdomains", "global", "Active Subs")
	activedomains.Entry.Unique = true
	activedomains.ShowLength = true
	alldomains := NewPage("All Subdomains", "global", "All Subs")
	alldomains.Entry.Unique = true
	alldomains.ShowLength = true
	assets := NewPage("Assets (Ex S3 etc)", "global", "Other assets")

	domainsobj := PanelObject{
		Header: "Subdomain  Enum",
	}

	domainlabel := widget.NewLabelWithStyle("Subdomain Enum", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	domainsobj.Objects = []fyne.CanvasObject{
		domainlabel,
		activedomains.GetButton(),
		alldomains.GetButton(),
		assets.GetButton(),
	}

	activeurls := NewPage("All Found Urls", "global", "All URLs")
	activeurls.Entry.Unique = true
	activeurls.ShowLength = true
	allurls := NewPage("Active Urls", "global", "Active URLs")
	allurls.Entry.Unique = true
	allurls.ShowLength = true

	urlsobj := PanelObject{
		Header: "WayBack Urls",
	}

	urlslabel := widget.NewLabelWithStyle("WayBack Urls", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	urlsobj.Objects = []fyne.CanvasObject{
		urlslabel,
		activeurls.GetButton(),
		allurls.GetButton(),
	}

	notespage := NewPage("My Notes", "notes", "Notes")
	findingspage := NewPage("My Findings", "notes", "Findings")
	// credpage := NewPage("API and Credentials", "notes", "Creds")
	credpage := NewCredPage()
	NotConfirmed := NewPage("Not Confirmed", "notes", "Uncertain")

	notesobj := PanelObject{
		Header: "Notes",
	}

	noteslabel := widget.NewLabelWithStyle("Notes and Stuff", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	notesobj.Objects = []fyne.CanvasObject{
		noteslabel,
		notespage.GetButton(),
		findingspage.GetButton(),
		credpage.GetButton(),
		NotConfirmed.GetButton(),
	}

	checklistobj := PanelObject{
		Header: "Web CheckLists",
	}
	cpage := NewCheckLists()
	cpage.ActiveSubs = activedomains

	opage := NewOrgCheckList()

	checklistlabel := widget.NewLabelWithStyle("CheckLists", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	checklistobj.Objects = []fyne.CanvasObject{
		checklistlabel,
		cpage.GetButton(),
		opage.GetButton(),
	}

	// myl := []Page{
	// 	*notespage,
	// 	*findingspage,
	// 	*credpage,
	// }

	// myt := GetTree(myl, "Notes")
	// myt.OpenAllBranches()
	// card := widget.NewCard("", "List Check", GetList(myl))

	//Create A Panel
	nav := NavPanel{}
	nav.Objects = []fyne.CanvasObject{
		// myt,
		summaryobj.Render(),
		domainsobj.Render(),
		urlsobj.Render(),
		notesobj.Render(),
		checklistobj.Render(),
	}

	//add all pages to container for update
	PageArr := []Page{
		*inscopepage,
		*outofscopep,
		*invalidvulns,
		*alldomains,
		*activedomains,
		*allurls,
		*activeurls,
		*assets,
		*notespage,
		*findingspage,
		*NotConfirmed,
	}

	E.CheckList = cpage
	E.OrgCheckL = opage

	E.Pages = PageArr

	nav.Color = color.RGBA{R: 49, G: 48, B: 66}

	return container.NewBorder(nil, nil, nav.Render(), nil, &E.Cont)

}
