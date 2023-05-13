package dbtools

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type Iata_code struct {
	Code           string
	Country        string
	Province       string
	City           string
	Isp_city_en    string
	Isp_city_4_ntt string
	Country_code   string
	Province_code  string
	Province_en    string
}

func matchesPattern(prefix string, s string) bool {
	pattern := fmt.Sprintf(`^(.*[-.\d]|^)%s[-.\d].*$`, prefix)
	log.Println(pattern)
	r, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regular expression:", err)
		return false
	}

	return r.MatchString(s)
}

func GetIATAGeo(ptr string) (*ipgeo.IPGeoData, error) {
	// // 混淆字符，这些字符可能会影响到 IATA 的判断结果
	// var obfuscated_str = []string{
	// 	"china",
	// 	"telecom",
	// 	"mobile",
	// 	"atlas",
	// 	"netvigator",
	// 	"chimon",
	// 	"customer",
	// 	"crossing",
	// 	"sttlwawb",
	// 	"france",
	// 	"seabone",
	// 	"best",
	// 	"fiord",
	// 	"taobao",
	// 	"border",
	// 	"viettel",
	// 	"aper",
	// 	"gldn",
	// 	"super",
	// 	"park",
	// 	"londen",
	// 	"melbi",
	// }
	// 处理 ptr
	var iata Iata_code
	// 转小写
	ptr = strings.ToLower(ptr)
	db := GetDB()
	iatas := []Iata_code{}
	db.Where("locate(isp_city_en, ?) > 0 or locate(code, ?) > 0 or locate(isp_city_4_ntt, ?) > 0 and isp_city_4_ntt <> ''", ptr, ptr, ptr).Find(&iatas)
	log.Println(iatas)
	for _, v := range iatas {
		if matchesPattern(v.Code, ptr) || matchesPattern(v.Isp_city_en, ptr) {
			iata = v
			break
		}
	}
	if iata.Country == "" {
		return nil, errors.New("no data")
	}
	return &ipgeo.IPGeoData{
		Country: iata.Country,
		Prov:    iata.Province,
		City:    iata.City,
	}, nil
}
