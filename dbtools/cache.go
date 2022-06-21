package dbtools

import (
	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Cache struct {
	Id         uint
	Uid        uint
	Ip         string
	Prefix     int
	Asnumber   string
	Country    string
	Prov       string
	City       string
	District   string
	Owner      string
	Isp        string
}

func SearchIP(ip string) (*Cache, error) {
	db := GetDB()
	c := Cache{}
	res := db.Where("ip = ?", ip).Take(&c)
	return &c, res.Error
}

func AddIP(data *ipgeo.IPGeoData, uid int) error {
	prefix := 32

	if data.District != "" {
		prefix = 24
	}

	db := GetDB()

	c := Cache{
		Ip:       data.IP,
		Uid:      uint(uid),
		Prefix:   prefix,
		Asnumber: data.Asnumber,
		Country:  data.Country,
		Prov:     data.Prov,
		City:     data.City,
		District: data.District,
		Owner:    data.Owner,
		Isp:      data.Isp,
	}

	res := db.Create(&c)
	return res.Error
}
