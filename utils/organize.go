package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// OrganizeData : Organizes data its best to leave tld size to 2
func OrganizeData(raw []string, tldsize int) map[string][]string {

	clustered := map[string][]string{}

	iponly := `$(((25[0-5])|(2[0-4][0-9])|(1[0-9]{2})|([1-9][0-9])|[0-9]{1})[.]{1}){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])^`

	re := regexp.MustCompile(iponly)

	ips := []string{}
	domains := []string{}

	for _, v := range raw {
		if re.MatchString(strings.TrimSpace(v)) {
			ips = append(ips, v)
		} else {

			domains = append(domains, v)
		}
	}

	apex := getrootnodes(domains, tldsize)

	//match all main domain to subdomains
	for _, v := range apex {
		//RegexSubsFromArray is far strict version of this
		//however it does not match the path for now
		//must come up with new regex for that
		//will do that later
		clustered[v] = *GetSubsFromArray(domains, v)
	}

	clustered["IP Addresses"] = ips

	return clustered

}

func getrootnodes(domains []string, tldsize int) []string {
	apex := map[string]bool{}

	//using tld size to find apex nodes

	for _, v := range domains {
		//also filter top domain paths
		//ex hackerone.com/abc static.hackerone.com/zab
		filtertldpath := strings.Split(v, "/")
		arr := strings.Split(filtertldpath[0], ".")
		if tldsize >= len(arr) {
			fmt.Println("No filtering can be done tldsize is incorrect ")
		}
		sub := len(arr) - tldsize
		req := arr[sub:]
		joined := ""
		for k, v := range req {
			if k != len(req)-1 {
				joined = joined + v + "."
			} else {
				joined = joined + v
			}
		}

		apex[joined] = true
	}

	results := []string{}

	for z := range apex {
		results = append(results, z)
	}

	return results

}

// RegexSubsFromArray ...
func RegexSubsFromArray(datax []string, query string) *[]string {
	data := strings.Join(datax, " ")
	results := []string{}
	// .* so that it can add extra paths ex : hackerone.com/appsec
	exp := `([a-zA-Z0-9-]+[.]+)+\Q` + query + `\E`
	re := regexp.MustCompile(exp)
	dat := re.FindAllStringSubmatch(data, -1)
	for _, z := range dat {
		if len(z) != 0 {
			results = append(results, string(z[0]))
		}
	}
	return &results
}

// GetSubsFromArray ...
func GetSubsFromArray(datax []string, query string) *[]string {
	results := []string{}

	for _, v := range datax {
		if strings.Contains(v, query) {
			results = append(results, v)
		}
	}

	return &results
}
