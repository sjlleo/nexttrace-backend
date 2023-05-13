package dbtools

import (
	"fmt"
	"math/big"
	"net"
	"strings"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Ip_whois struct {
	Id           uint
	Begin        int64
	End          int64
	Netname      string
	Owner        string
	Country_code string
}

func InetNtoA(ip uint64) string {
	return fmt.Sprintf("%d.%d.%d.%d",
		byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func InetAtoN(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

func SearchIPWhois(ip string) (*Ip_whois, error) {
	db := GetDB()
	w := Ip_whois{}
	ipInt := InetAtoN(ip)
	res := db.Where("begin <= ? and end >= ?", ipInt, ipInt).Last(&w)
	return &w, res.Error
}

func ParseIPRange(inet string) (int64, int64) {
	arr := strings.Split(inet, " - ")
	ccd := arr[0]
	fmt.Println(ccd)
	return InetAtoN(arr[0]), InetAtoN(arr[1])
}

func AddIPWhois(data *ipgeo.IpWhois) error {
	db := GetDB()
	if data.Netname != "IANA-BLOCK" {
		b, e := ParseIPRange(data.Inetnum)
		w := Ip_whois{
			Begin:        b,
			End:          e,
			Netname:      data.Netname,
			Owner:        data.Owner,
			Country_code: data.Country_code,
		}
		res := db.Create(&w)
		return res.Error
	} else {
		w := Ip_whois{
			Begin: InetAtoN(data.IP),
			End:   InetAtoN(data.IP),
		}
		res := db.Create(&w)
		return res.Error
	}
}
