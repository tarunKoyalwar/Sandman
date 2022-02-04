package ui

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/tarunKoyalwar/sandman/utils"
)

var PathsforGenericPages map[string]string = map[string]string{
	"in-scope":           "In Scope",
	"out-of-scope":       "Out of Scope",
	"not-accepted-vulns": "Not Accepted Vuls",
	"active-subs":        "Active Subs",
	"all-subs":           "All Subs",
	"other-assets":       "Other Assets",
	"all-urls":           "All URLs",
	"active-urls":        "Active URLs",
	"notes":              "Notes",
	"findings":           "Findings",
	"uncertain":          "Uncertain",
}

//GenericPage Handler : Handle Function Would be /page/pagename
func GenericPage(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	keys := make([]string, 0, len(PathsforGenericPages))
	for k := range PathsforGenericPages {
		keys = append(keys, k)
	}

	ActualPageName := PathsforGenericPages[vars["pname"]]

	if ActualPageName == "" {
		log.Printf("Page Not Found")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Page Not Found"))
		rw.Write([]byte(fmt.Sprintf("%v\n", keys)))
		return
	}

	//get page
	p, err := E.GetPage(ActualPageName)
	if err != nil {
		log.Printf("Page Not Found")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Page Not Found"))
		rw.Write([]byte(fmt.Sprintf("%v\n", keys)))
		return
	}

	if r.Method == "GET" {
		//Send outofscope data
		dat := p.Entry.GetData()
		rw.WriteHeader(http.StatusOK)
		rw.Write(utils.StringArraytoByteArr(dat))
		return
	}

	if r.Method == "POST" {
		//Set data to outofscope
		respbody, er := ioutil.ReadAll(r.Body)
		if er != nil {
			log.Panic(er)
		}

		// fmt.Println(string(respbody))

		query := r.URL.Query()
		Overwrite := utils.StringArrtoString(query["overwrite"])
		if Overwrite == "true" {
			p.Entry.Set(string(respbody), false)
			p.Entry.Refresh()
		} else {
			p.Entry.Set(string(respbody), true)
			p.Entry.Refresh()
		}

		// fmt.Println(p.Entry.GetData())

	}

}

//WorkingSubdomain : Gets Current Working Subdomain Handle Function would be /web/working
func WorkingSubdomain(rw http.ResponseWriter, r *http.Request) {
	respdata := E.CheckList.ActiveTreeIndex

	if respdata != "" {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(respdata))
		return
	} else {
		rw.WriteHeader(http.StatusNotFound)
	}
}

// OrgCheckListRoute : List Organization CheckList Handle Function would be /org/checklist
func OrgCheckListRoute(rw http.ResponseWriter, r *http.Request) {
	respdata := E.OrgCheckL.MyTree.Nodes
	bin, err := json.Marshal(respdata)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(bin)
	// rw.WriteHeader(http.StatusOK)

}

// WebCheckListList WebCheckList Handle Function would be /web/checklist
func WebCheckList(rw http.ResponseWriter, r *http.Request) {
	respdata := E.CheckList.GetDefaultCheckList()
	bin, err := json.Marshal(respdata)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	rw.Write(bin)
	// rw.WriteHeader(http.StatusOK)
}

//Org GET PUT Tool Output
//Handle Function would be /org/checklist_item_name
func OrgToolOutput(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	decoded, _ := url.QueryUnescape(vars["orgitem"])

	found, err := E.OrgCheckL.MyTree.FindObject(decoded)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	if r.Method == "GET" {
		//Send outofscope data
		dat := found.ToolEntry.GetData()
		rw.WriteHeader(http.StatusOK)
		rw.Write(utils.StringArraytoByteArr(dat))
		return
	}

	if r.Method == "POST" {
		//Set data to outofscope
		respbody, _ := ioutil.ReadAll(r.Body)

		query := r.URL.Query()
		Overwrite := utils.StringArrtoString(query["overwrite"])
		if Overwrite == "true" {
			found.ToolEntry.Set(string(respbody), false)
		} else {
			found.ToolEntry.Set(string(respbody), true)
		}

	}
}

// WebToolOutput :  CheckList GET PUT Tool Output Handle Function would be /web/subdomain/checklist_item_name
func WebToolOutput(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	allsubs := []string{}
	for k := range E.CheckList.MyTrees {
		allsubs = append(allsubs, k)
	}

	if E.CheckList.MyTrees[vars["sub"]] == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		tip := "If sub is not in above list however is present in active-subs\n then open checklist in UI at least once\nThis was done to avoid creating of unnecessary objects\n"
		rw.Write([]byte(fmt.Sprintf("No checklist found for that subdomain %v\n %v\n %v", allsubs, tip, r.RequestURI)))
		return
	}

	found, err := E.CheckList.MyTrees[vars["sub"]].FindObject(vars["webitem"])
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(err.Error()))
		return
	}

	if r.Method == "GET" {
		//Send outofscope data
		dat := found.ToolEntry.GetData()
		rw.WriteHeader(http.StatusOK)
		rw.Write(utils.StringArraytoByteArr(dat))
		return
	}

	if r.Method == "POST" {
		//Set data to outofscope
		respbody, _ := ioutil.ReadAll(r.Body)

		query := r.URL.Query()
		Overwrite := utils.StringArrtoString(query["overwrite"])
		if Overwrite == "true" {
			found.ToolEntry.Set(string(respbody), false)
		} else {
			found.ToolEntry.Set(string(respbody), true)
		}

	}
}

// Commit : This is used to commit changes to mongodb ,Handle Function Would be /commit
func Commit(rw http.ResponseWriter, r *http.Request) {
	E.SaveState()
}
