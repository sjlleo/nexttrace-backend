package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func AiwenTech(ip string) (*IPGeoData, error) {

	url := "https://api.ipplus360.com/ip/geo/v1/district/?key=" + AiWenToken + "&ip=" + ip

	client := &http.Client{
		// 5 秒超时
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("埃文科技 请求超时(4s)")
		return &IPGeoData{}, err
	}

	body, _ := io.ReadAll(content.Body)
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
	case "上海市":
		city = ""
	case "北京市":
		city = ""
	case "重庆市":
		city = ""
	case "天津市":
		city = ""
	}

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("asnumber").String(),
		Country:  country,
		Prov:     prov,
		City:     city,
		District: res.Get("district").String(),
		Owner:    res.Get("owner").String(),
		Isp:      res.Get("isp").String(),
	}, nil
}
