package dbtools

import (
	"net"
	"strconv"

	"github.com/c-robinson/iplib"
)

type Router struct {
	Id           int
	Ip           string
	Ipend        string
	Prefix       string
	Asnumber     string
	Fromasnumber string
}

func AddRouter(ip string, prefix string, data map[string][]string) {
	p, _ := strconv.Atoi(prefix)
	n := iplib.NewNet(net.ParseIP(ip), p)

	db := GetDB()
	for k, v := range data {
		for _, j := range v {
			db.Model(&Router{}).Create(&Router{
				Ip:           n.FirstAddress().String(),
				Ipend:        n.LastAddress().String(),
				Prefix:       prefix,
				Asnumber:     j,
				Fromasnumber: k,
			})
		}
	}
}

func SearchRouter(ip string) []Router {
	var r []Router

	db := GetDB()

	ip_net := net.ParseIP(ip)

	if ip_net.To4() != nil {
		db.Raw("SELECT * FROM router WHERE INET_ATON(ip) <= INET_ATON(?) + 1 AND INET_ATON(?) <= INET_ATON(ipend)", ip, ip).Scan(&r)
	} else {
		db.Raw("SELECT * FROM router WHERE INET6_ATON(ip) <= INET6_ATON(?)  AND INET6_ATON(?) <= INET6_ATON(ipend)", ip, ip).Scan(&r)
	}
	return r
}
