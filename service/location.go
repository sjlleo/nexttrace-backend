package service

import (
	"github.com/sjlleo/nexttrace-backend/dbtools"
	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

func LangApply(r *ipgeo.IPGeoData, c dbtools.CityInfo) {
	r.Country = c.CountryZH
	r.CountryEN = c.CountryEN
	if r.Prov != "" && r.Prov != r.Country {
		r.Prov = c.RegionZH
		r.RegionEN = c.RegionEN
	}
	if r.City != "" && r.Prov != r.City {
		r.City = c.CityZH
		r.CityEN = c.CityEN
	}
	r.Lat = c.CenterLocationLat
	r.Lng = c.CenterLocationLng
}

func GetLocation(r *ipgeo.IPGeoData) {
	if c := dbtools.SearchData(r); c.CountryEN != "" {
		LangApply(r, c)
	} else {
		g := dbtools.Geo{
			Country: r.Country,
			Prov:    r.Prov,
			City:    r.City,
		}
		if c, err := dbtools.GetInfoData(&g); err == nil {
			LangApply(r, *c)
			dbtools.AddData(c)
		}
	}
}
