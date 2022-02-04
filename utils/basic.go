package utils

import (
	"encoding/binary"
	"log"
	"net"
	"regexp"
)

// MultiCIDRtoHosts ...
func MultiCIDRtoHosts(cidrs []string) []string {
	res := []string{}
	for _, v := range cidrs {
		res = append(res, CIDRHosts(v)...)
	}
	return res
}

// CIDRHosts ...
func CIDRHosts(netw string) []string {
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(netw)
	if err != nil {
		log.Fatal(err)
	}
	// convert IPNet struct mask and address to uint32
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	// fmt.Printf("Mask is %v\n", mask)
	// find the start IP address
	start := binary.BigEndian.Uint32(ipv4Net.IP)
	// fmt.Printf("Start is %v\n", start)
	// find the final IP address
	finish := (start & mask) | (mask ^ 0xffffffff)
	// fmt.Printf("finish is %v\n", finish)
	// make a slice to return host addresses
	var hosts []string
	// loop through addresses as uint32.
	for i := start; i <= finish; i++ {
		// convert back to net.IPs
		// Create IP address of type net.IP. IPv4 is 4 bytes, IPv6 is 16 bytes.
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		hosts = append(hosts, ip.String())
	}
	// return a slice of strings containing IP addresses
	return hosts
}

// RegexGetCIDRs :  returns cidr and basic filtering
func RegexGetCIDRs(dat []string) ([]string, string) {
	regexcidr := `(((25[0-5])|(2[0-4][0-9])|(1[0-9]{2})|([1-9][0-9])|[0-9]{1})[.]{1}){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9][0-9]|[0-9])[/]([1-9][0-9])`

	cidrs := []string{}
	re := regexp.MustCompile(regexcidr)

	all := ""

	for _, v := range dat {
		res := re.FindAllString(v, -1)
		if res != nil {
			cidrs = append(cidrs, res...)
			rem := re.ReplaceAllString(v, "")
			all = all + rem
		} else {
			all = all + v
		}
	}

	return cidrs, all
}
