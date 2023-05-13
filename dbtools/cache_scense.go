package dbtools

import "github.com/sjlleo/nexttrace-backend/ipgeo"

type Cache_scense struct {
	Id     uint
	Ip     string
	Scense string
}

func SearchIPScense(ip string) (*Cache_scense, error) {
	db := GetDB()
	c := Cache_scense{}
	res := db.Where("ip = ?", ip).Take(&c)
	return &c, res.Error
}

func AddIPScense(data *ipgeo.IPGeoData, scene string) error {

	db := GetDB()

	c := Cache_scense{
		Ip:  data.IP,
		Scense: scene,
	}

	res := db.Create(&c)
	return res.Error
}
