package dbtools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Coordinate struct {
	Lat string
	Lng string
}

type Ipv4_asn struct {
	Asn    int
	Org    string
	Isp    string
	Domain string
}

type IXData struct {
	Ipaddr4    string `json:"ipaddr4"`
	Ipaddr6    string `json:"ipaddr6"`
	Speed      int    `json:"speed"`
	IxName     string `json:"ix_name"`
	IxCountry  string `json:"ix_country"`
	IxCity     string `json:"ix_city"`
	NetAsn     int    `json:"net_asn"`
	NetName    string `json:"net_name"`
	NetWebsite string `json:"net_website"`
	NetAka     string `json:"net_aka"`
}

func TestMigrateCityInfo(t *testing.T) {
	db := GetDB()
	db.AutoMigrate(&CityInfo{})
}

func TestGeo(t *testing.T) {
	// g := Geo{
	// 	Country: "US",
	// 	Prov:    "California",
	// 	City:    "Los Angeles",
	// }
	g := Geo{
		Country: "",
		Prov:    "",
		City:    "",
	}
	log.Print(ExistInCityInfo(&g))
}

func TestNeiborHop(t *testing.T) {
	r, _ := FindNeiborHop("116.4.8.1")
	log.Println(r)
}

func TestAddIXHop(t *testing.T) {
	d := []IXData{}
	f, _ := ioutil.ReadFile("./ip.out.json")
	json.Unmarshal(f, &d)
	for _, v := range d {
		// 判断数据是否存在
		c := Cache{}
		c6 := Cache_v6{}
		db := GetDB()

		// 查看是不是已经有 district
		var both_exist_v4 bool
		var both_exist_v6 bool
		// IPv4
		if v.Ipaddr4 != "" {
			db.Model(Cache{}).Where("ip=?", v.Ipaddr4).Find(&c)
			if len(c.District) > 1 {
				both_exist_v4 = true
			}
		}
		// IPv6
		if v.Ipaddr6 != "" {
			db.Model(Cache_v6{}).Where("ip=?", v.Ipaddr6).Find(&c6)
			if len(c6.District) > 1 {
				both_exist_v6 = true
			}
		}

		if both_exist_v4 && both_exist_v6 {
			continue
		}

		// 获取地理位置
		//
		// 首先查询省份
		ci := CityInfo{}
		// 去除缩写
		tmp := strings.Split(v.IxCity, "/")
		if len(tmp) > 1 {
			v.IxCity = tmp[0]
		}
		// 特殊地区处理
		if v.IxCountry == "TW" {
			ci.CountryZH = "China"
			ci.RegionZH = "Taiwan"
			// 过滤敏感内容
			if v.IxCity == "Taipei/Taiwan" || strings.HasSuffix(v.IxCity, "Taiwan") || strings.HasSuffix(v.IxCity, "is a country") {
				ci.CityZH = "Taipei"
			} else {
				ci.CityZH = v.IxCity
			}
		} else if v.IxCountry == "HK" {
			ci.CountryZH = "China"
			ci.RegionZH = "Hong Kong"
		} else {
			db.Model(CityInfo{}).Where("country_en=?", CountryCodeToCountryName(v.IxCountry)).Where("city_en=?", v.IxCity).Take(&ci)
		}
		if ci.CountryZH != "" {
			c.Ip = v.Ipaddr4
			c.Asnumber = strconv.Itoa(v.NetAsn)
			c.Prefix = 32
			c.Country = ci.CountryZH
			c.Prov = ci.RegionZH
			c.City = ci.CityZH

			ix_note := ""
			if v.Speed == 0 {
				ix_note = v.IxName + " - " + v.NetName
			} else if v.Speed < 1000 {
				ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed) + "Mbps"
			} else {
				ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed/1000) + "Gbps"
			}
			c.District = ix_note

			c6.Ip = v.Ipaddr6
			c6.Asnumber = strconv.Itoa(v.NetAsn)
			c6.Prefix = 128
			c6.Country = ci.CountryZH
			c6.Prov = ci.RegionZH
			c6.City = ci.CityZH
			c6.District = ix_note
		} else {
			// 数据库无此数据，从 ipinfo 获取
			if v.Ipaddr4 != "" {
				r, _ := ipgeo.IPInfoTMP(v.Ipaddr4)
				c.Ip = v.Ipaddr4
				c.Asnumber = strconv.Itoa(v.NetAsn)
				c.Prefix = 32
				c.Country = r.Country
				c.Prov = r.Prov
				c.City = r.City
				var ix_note string
				if v.Speed == 0 {
					ix_note = v.IxName + " - " + v.NetName
				} else if v.Speed < 1000 {
					ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed) + "Mbps"
				} else {
					ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed/1000) + "Gbps"
				}
				c.District = ix_note
			}
			if v.Ipaddr6 != "" {
				r, _ := ipgeo.IPInfoTMP(v.Ipaddr6)
				c6.Ip = v.Ipaddr6
				c6.Asnumber = strconv.Itoa(v.NetAsn)
				c6.Prefix = 128
				c6.Country = r.Country
				c6.Prov = r.Prov
				c6.City = r.City
				var ix_note string
				if v.Speed == 0 {
					ix_note = v.IxName + " - " + v.NetName
				} else if v.Speed < 1000 {
					ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed) + "Mbps"
				} else {
					ix_note = v.IxName + " - " + v.NetName + " - " + strconv.Itoa(v.Speed/1000) + "Gbps"
				}
				c6.District = ix_note
			}
			log.Println(c)
		}

		// log.Println(c)
		// log.Println(c6)
		// 判断 IPv4
		if v.Ipaddr4 != "" {
			if db.Model(Cache{}).Where("ip=?", v.Ipaddr4).Updates(&c).RowsAffected == 0 {
				db.Model(Cache{}).Create(&c)
			}
		}
		// 判断 IPv6
		if v.Ipaddr6 != "" {
			if db.Model(Cache_v6{}).Where("ip=?", v.Ipaddr6).Updates(&c6).RowsAffected == 0 {
				db.Model(Cache_v6{}).Create(&c6)
			}
		}
	}
}

func TestMergeData(t *testing.T) {
	db := GetDB()
	// Get All result from asn
	var asn []Ipv4_asn
	db.Find(&asn)
	for _, v := range asn {
		var new_asn Asn
		db.Model(Asn{}).Where("asn =?", v.Asn).Find(&new_asn)
		if new_asn.Isp == "" {
			// 找不到记录，直接插入数据
			new_asn = Asn{
				Asn:     v.Asn,
				Isp:     v.Isp,
				Alias:   v.Org,
				Website: v.Domain,
			}
			db.Create(&new_asn)
			continue
		}

		if new_asn.Website == "" {
			new_asn.Website = v.Domain
			db.Model(Asn{}).Where("asn =?", v.Asn).Updates(&new_asn)
		}
	}
}

func ExistInCityInfo(g *Geo) bool {
	country := CountryCodeToCountryName(g.Country)
	if country == "" {
		country = g.Country
	}
	c := CityInfo{}
	db := GetDB()
	tx := db.Model(&CityInfo{})
	tx = tx.Where("country_en = ? OR country_zh = ?", country, country)

	if g.Prov != "" {
		tx = tx.Where("region_en = ? OR region_zh = ?", g.Prov, g.Prov)
	}

	if g.City != "" {
		tx = tx.Where("city_en = ? OR city_zh = ?", g.City, g.City)
	}

	tx.Take(&c)
	return c.CountryZH != ""
}

func TestShowAllCities(t *testing.T) {
	r := ShowAllCities()

	for _, v := range *r {
		g := Geo{
			Country: v.Country,
			Prov:    v.Prov,
			City:    v.City,
		}
		// 判断是不是一个国家

		// 判断是不是已经在数据库里面了
		if !ExistInCityInfo(&g) {
			if res, err := GetInfoData(&g); err == nil {
				log.Println(res)
				AddData(res)
			}
		}

		<-time.After(time.Millisecond * 10)
	}

}

func TestNewFunc(t *testing.T) {
	n := iplib.NewNet(net.ParseIP("1.1.1.1"), 24)
	fmt.Println(n.FirstAddress()) // 2001:db8::
	fmt.Println(n.LastAddress())  // 2001:db8:0:ff:ffff:ffff:ffff:ffff
}
