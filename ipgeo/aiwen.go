package ipgeo

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func AiwenTech(ip string) (*IPGeoData, error) {
	url := "https://api.ipplus360.com/ip/geo/v1/district/?key=" + AiWenToken + "&ip=" + ip

	client := &http.Client{
		// 4 秒超时
		Timeout: 4 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("埃文科技 请求超时(4s)")
		return &IPGeoData{}, err
	}

	body, _ := ioutil.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	var country string
	var prov string
	var city string
	res = res.Get("data")
	country = res.Get("country").String()
	prov = res.Get("prov").String()
	city = res.Get("city").String()

	switch prov {
	case "中国香港":
		country = "中国"
		prov = "香港"
		city = ""
	}

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("asnumber").String(),
		Country:  country,
		Prov:     prov,
		City:     city,
		Owner:    res.Get("owner").String(),
		Isp:      res.Get("isp").String(),
	}, nil
}
