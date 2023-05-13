package dbtools

import (
	"log"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Cache_v6 struct {
	Id       uint
	Uid      uint
	Ip       string
	Prefix   int
	Asnumber string
	Country  string
	Prov     string
	City     string
	District string
	Owner    string
	Isp      string
}

func SearchIPv6(ip string) (*Cache_v6, error) {
	db := GetDB()
	c := Cache_v6{}
	res := db.Where("INET6_ATON(?) >> (128-prefix) = INET6_ATON(ip) >> (128-prefix)", ip).Order("prefix desc").Take(&c)
	return &c, res.Error
}

func UpdateIPv6(ip string, prefix int, asnumber string, country string, prov string, city string, district string) error {
	db := GetDB()
	c, _ := SearchIPv6(ip)
	log.Println(c)
	if c.Country == "" {
		c.Ip = ip
		c.Asnumber = asnumber
		c.Country = country
		c.Prov = prov
		c.City = city
		c.District = district
		c.Prefix = prefix
		err := db.Create(&c).Error
		return err
	}

	res := db.Model(&c).Updates(Cache_v6{
		Prefix:   prefix,
		Asnumber: asnumber,
		Country:  country,
		Prov:     prov,
		City:     city,
	})
	return res.Error
}

func AddIPv6(data *ipgeo.IPGeoData, uid int) error {
	prefix := 128

	if data.District != "" {
		prefix = 64
	}

	db := GetDB()

	c := Cache_v6{
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
