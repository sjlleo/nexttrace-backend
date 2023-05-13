package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func IPData(ip string) (*IPGeoData, error) {

	url := "https://api.ipdata.co/" + ip + "?api-key=" + IPDataToken + "&fields=ip,is_eu,city,region,region_code,country_name,country_code,continent_name,continent_code,latitude,longitude,postal,calling_code,flag,emoji_flag,emoji_unicode,asn"

	client := &http.Client{
		// 5 秒超时
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("IPData 请求超时(4s)")
		return &IPGeoData{}, err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)
	var country string
	var prov string
	var city string
	country = res.Get("country_name").String()
	prov = res.Get("region").String()
	city = res.Get("city").String()
	switch country {
	case "Hong Kong":
		country = "中国"
		prov = "香港"
		city = ""
	case "Taiwan":
		country = "中国"
		prov = "台湾省"
	case "Macao":
		country = "中国"
		prov = "澳门"
		city = ""
	}

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("asn").Get("asnumber").String(),
		Country:  country,
		Prov:     prov,
		City:     city,
	}, nil
}
