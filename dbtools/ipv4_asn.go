package dbtools

import (
	"strconv"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Ipv4_asn struct {
	Asn    int
	Org    string
	Isp    string
	Domain string
}

func SearchASNData(asnNumber string) (*Ipv4_asn, error) {
	// ASN 数据
	db := GetDB()
	asn := Ipv4_asn{}
	res := db.Where("asn = ?", asnNumber).Find(&asn)
	return &asn, res.Error
}

func AddASNData(data *ipgeo.IPGeoData) error {
	// ASN 数据
	db := GetDB()
	asnNum, _ := strconv.Atoi(data.Asnumber)
	asn := Ipv4_asn{
		Asn:    asnNum,
		Org:    data.Owner,
		Isp:    data.Isp,
		Domain: data.Domain,
	}
	res := db.Create(&asn)
	return res.Error
}
