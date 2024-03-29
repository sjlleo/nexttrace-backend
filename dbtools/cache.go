package dbtools

import (
	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Cache struct {
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

func ShowAllCities() *[]Cache {
	db := GetDB()
	c := []Cache{}
	db.Model(&Cache{}).Distinct("country", "prov", "city").Find(&c)
	return &c
}

func SearchIP(ip string) (*Cache, error) {
	db := GetDB()
	c := Cache{}

	res := db.Where("INET_ATON(?) & ~(pow(2,32-prefix)-1) = INET_ATON(ip) & ~(pow(2,32-prefix)-1)", ip).Order("prefix desc").Take(&c)
	return &c, res.Error
}

func UpdateIP(ip string, prefix int, asnumber string, country string, prov string, city string, district string) error {
	db := GetDB()
	c, _ := SearchIP(ip)
	// 如果找不到
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
	res := db.Model(&c).Updates(Cache{
		Prefix:   prefix,
		Asnumber: asnumber,
		Country:  country,
		Prov:     prov,
		City:     city,
		District: district,
	})
	return res.Error
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
