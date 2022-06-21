package ipgeo

import (
	"io/ioutil"
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

	body, _ := ioutil.ReadAll(content.Body)
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
		region = "台湾"
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
