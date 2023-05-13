package ipgeo

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

type IpWhois struct {
	IP           string
	Owner        string
	Inetnum      string
	Netname      string
	Country_code string
}

func AiwenIPWhois(ip string) (*IpWhois, error) {
	url := "https://api.ipplus360.com/ip/info/v1/ipWhois/?key=" + AiWenToken + "&ip=" + ip

	client := &http.Client{
		// 5 秒超时
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("埃文科技 请求超时(4s)")
		return &IpWhois{}, err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	IW := IpWhois{}
	// 查询IP
	IW.IP = ip
	// 获取 data 内的结构体信息
	res = res.Get("data")
	// 网段信息
	IW.Inetnum = res.Get("inetnum").String()
	// 网段名称
	IW.Netname = res.Get("netname").String()
	// 网段国家
	IW.Country_code = res.Get("country_code").String()
	// 网段归属
	IW.Owner = res.Get("owner").String()
	return &IW, nil
}
