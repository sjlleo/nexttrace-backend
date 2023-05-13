package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func ipwhois(ip string) (*IPGeoData, error) {
	url := "https://ipwho.is/" + ip

	client := &http.Client{
		// 2 秒超时
		Timeout: 2 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:101.0) Gecko/20100101 Firefox/101.0")
	req.Header.Set("Referer", "https://ipwhois.io/")

	content, err := client.Do(req)
	if err != nil {
		log.Println("IPInsight 请求超时(2s)")
		return &IPGeoData{}, err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	return &IPGeoData{
		IP:       ip,
		Asnumber: res.Get("connection").Get("asn").String(),
		Country:  res.Get("country").String(),
		City:     res.Get("city").String(),
		Prov:     res.Get("region").String(),
		Owner:    res.Get("connection").Get("org").String(),
		Isp:      res.Get("connection").Get("isp").String(),
		Domain:   res.Get("connection").Get("domain").String(),
	}, nil
}
