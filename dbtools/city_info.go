package dbtools

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unicode"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
	"github.com/tidwall/gjson"
)

type CityInfo struct {
	CountryZH          string
	CountryEN          string
	RegionZH           string
	RegionEN           string
	CityZH             string
	CityEN             string
	CenterLocationLat  float64
	CenterLocationLng  float64
	BoundsNortheastLat float64
	BoundsNortheastLng float64
	BoundsSouthwestLat float64
	BoundsSouthwestLng float64
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}

func AddData(c *CityInfo) {
	db := GetDB()
	db.Model(c).Create(c)
}

func CountryCodeToCountryName(country_code string) string {
	filePtr, _ := ioutil.ReadFile("/root/develop/nexttrace-backend/iso3166-1.json")
	return gjson.Get(string(filePtr), country_code).String()
}

func SearchData(r *ipgeo.IPGeoData) CityInfo {
	country := CountryCodeToCountryName(r.Country)
	if country == "" {
		country = r.Country
	}
	c := CityInfo{}
	db := GetDB()
	tx := db.Model(&CityInfo{})
	tx = tx.Where("country_en = ? OR country_zh = ?", country, country)

	if r.Prov != "" {
		tx = tx.Where("region_en = ? OR region_zh = ?", r.Prov, r.Prov)
	}

	if r.City != "" {
		tx = tx.Where("city_en = ? OR city_zh = ?", r.City, r.City)
	}

	tx.Take(&c)
	return c
}

func GetInfoData(g *Geo) (*CityInfo, error) {
	country_en := CountryCodeToCountryName(g.Country)
	var url_final string
	if country_en != "" {
		url_final = "https://maps.googleapis.com/maps/api/geocode/json?key=YOUR_GOOGLEMAP_KEY&address=" + url.QueryEscape(g.City+","+g.Prov+","+country_en)
	} else {
		url_final = "https://maps.googleapis.com/maps/api/geocode/json?key=YOUR_GOOGLEMAP_KEY&address=" + url.QueryEscape(g.City+","+g.Prov+","+g.Country)
	}
	client := &http.Client{
		// 5 秒超时
		Timeout: 5 * time.Second,
	}
	resp, _ := http.NewRequest("GET", url_final, nil)
	if !IsChinese(g.Prov) {
		resp.Header.Add("Accept-Language", "zh-CN")
	}
	response, _ := client.Do(resp)
	defer response.Body.Close()
	b, _ := io.ReadAll(response.Body)
	//
	info := CityInfo{}
	value := gjson.Get(string(b), "results.0.geometry.location")
	info.CenterLocationLat, _ = strconv.ParseFloat(value.Get("lat").String(), 32)
	info.CenterLocationLng, _ = strconv.ParseFloat(value.Get("lng").String(), 32)

	// 中文
	value = gjson.Get(string(b), "results.0.address_components")
	arr := value.Array()
	arr_len := len(arr)
	if arr_len == 0 {
		return nil, errors.New("no data")
	}
	if !IsChinese(g.Prov) {
		if g.Prov == "Taiwan" || g.Country == "TW" {
			info.CountryZH = "中国"
			info.RegionZH = "台湾省"
			info.CityZH = arr[0].Get("long_name").String()
		} else if arr_len == 1 {
			info.CountryZH = arr[0].Get("long_name").String()
		} else if arr_len == 2 {
			info.CountryZH = arr[1].Get("long_name").String()
			info.RegionZH = arr[0].Get("long_name").String()
		} else {
			if arr[arr_len-1].Get("types.0").String() != "country" {
				info.CountryZH = arr[arr_len-2].Get("long_name").String()
				info.CityZH = arr[0].Get("long_name").String()
			} else {
				info.CountryZH = arr[arr_len-1].Get("long_name").String()
				info.RegionZH = arr[arr_len-2].Get("long_name").String()
				info.CityZH = arr[0].Get("long_name").String()
			}

			if info.CountryZH == "美国" {
				info.RegionZH += "州"
			}

			if info.CityZH == info.RegionZH {
				info.CityZH = ""
			}
		}
		if country_en == "" {
			country_en = g.Country
		}
		info.CountryEN = country_en
		info.RegionEN = g.Prov
		info.CityEN = g.City
	} else {
		// 特殊地区单独判断
		if g.Prov == "台湾省" {
			info.CountryEN = "China"
			info.RegionEN = "Taiwan"
			info.CityEN = arr[0].Get("long_name").String()
		} else if arr_len == 1 {
			info.CountryEN = arr[0].Get("long_name").String()
		} else if arr_len == 2 {
			info.CountryEN = arr[1].Get("long_name").String()
			info.RegionEN = arr[0].Get("long_name").String()
		} else {
			info.CountryEN = arr[arr_len-1].Get("long_name").String()
			info.RegionEN = arr[arr_len-2].Get("long_name").String()
			info.CityEN = arr[0].Get("long_name").String()

			if info.RegionEN == info.CityEN {
				info.CityEN = ""
			}
		}

		info.CountryZH = g.Country
		info.RegionZH = g.Prov
		info.CityZH = g.City
	}
	info.BoundsNortheastLat, _ = strconv.ParseFloat(gjson.Get(string(b), "results.0.geometry.bounds.northeast.lat").String(), 32)
	info.BoundsNortheastLng, _ = strconv.ParseFloat(gjson.Get(string(b), "results.0.geometry.bounds.northeast.lng").String(), 32)

	info.BoundsSouthwestLat, _ = strconv.ParseFloat(gjson.Get(string(b), "results.0.geometry.bounds.southwest.lat").String(), 32)
	info.BoundsSouthwestLng, _ = strconv.ParseFloat(gjson.Get(string(b), "results.0.geometry.bounds.southwest.lng").String(), 32)
	return &info, nil
}
