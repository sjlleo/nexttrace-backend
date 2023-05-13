package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func IPInsight(ip string) (*IPGeoData, string, error) {
	url := "https://api.ipinsight.io/ip/" + ip + "?token=" + IPInsightToken

	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0")
	req.Header.Set("Referer", "https://ipinsight.io/")

	content, err := client.Do(req)
	if err != nil {
		log.Println("IPInsight 请求超时(2s)")
		return &IPGeoData{}, "", err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	var country string
	var region string

	country = res.Get("country_name").String()
	region = res.Get("region_name").String()

	switch region {
	case "Hong Kong":
		country = "中国"
		region = "香港"
	}

	return &IPGeoData{
		IP:      ip,
		Source:  res.Get("continent_code").String(),
		Country: country,
		Prov:    region,
		City:    res.Get("city_name").String(),
	}, res.Get("country_code").String(), nil
}
