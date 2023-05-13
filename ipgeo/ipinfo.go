package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func IPInfo(ip string) (*IPGeoData, error) {
	url := "https://ipinfo.io/" + ip + "?token=" + IPInfoToken

	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("IPInfo 请求超时(2s)")
		return nil, err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	var (
		country string
		region  string
		city    string
	)

	country = res.Get("country").String()
	region = res.Get("region").String()
	city = res.Get("city").String()

	switch country {
	case "HK":
		country = "中国"
		region = "香港"
		city = ""
	case "TW":
		country = "中国"
		region = "台湾省"
		switch city {
		case "Taipei":
			city = "台北市"
		case "Chang-hua":
			city = "彰化县"
		case "Hsinchu":
			city = "新竹市"
		case "Taoyuan City":
			city = "桃园市"
		case "Fengyuan":
			city = "丰原市"
		case "Banqiao":
			city = "新北市"
		case "Tainan":
			city = "台南市"
		case "Taichung":
			city = "台中市"
		}
	case "SG":
		country = "新加坡"
		region = ""
		city = ""
	case "MO":
		country = "中国"
		region = "澳门"
		city = ""
	}

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("asn").Get("asn").String(),
		Country:  country,
		Prov:     region,
		City:     city,
		Isp:      res.Get("asn").Get("domain").String(),
	}, nil
}

func IPInfoTMP(ip string) (*IPGeoData, error) {
	url := "https://ipinfo.io/" + ip + "?token=68578dc0b77d25"

	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("IPInfo 请求超时(2s)")
		return nil, err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	var (
		country string
		region  string
		city    string
	)

	country = res.Get("country").String()
	region = res.Get("region").String()
	city = res.Get("city").String()

	switch country {
	case "HK":
		country = "中国"
		region = "香港"
		city = ""
	case "TW":
		country = "中国"
		region = "台湾省"
		switch city {
		case "Taipei":
			city = "台北市"
		case "Chang-hua":
			city = "彰化县"
		case "Hsinchu":
			city = "新竹市"
		case "Taoyuan City":
			city = "桃园市"
		case "Fengyuan":
			city = "丰原市"
		case "Banqiao":
			city = "新北市"
		case "Tainan":
			city = "台南市"
		case "Taichung":
			city = "台中市"
		}
	case "SG":
		country = "新加坡"
		region = ""
		city = ""
	case "MO":
		country = "中国"
		region = "澳门"
		city = ""
	}

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("asn").Get("asn").String(),
		Country:  country,
		Prov:     region,
		City:     city,
		Isp:      res.Get("asn").Get("domain").String(),
	}, nil
}
