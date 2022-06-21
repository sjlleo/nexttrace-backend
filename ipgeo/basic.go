package ipgeo

type IPGeoData struct {
	IP       string `json:"ip"`
	Asnumber string `json:"asnumber"`
	Country  string `json:"country"`
	Prov     string `json:"prov"`
	City     string `json:"city"`
	District string `json:"district"`
	Owner    string `json:"owner"`
	Isp      string `json:"isp"`
	Domain   string `json:"domain"`
	Source   string `json:"source"`
}

// 2 家服务需要付费，请自行购买后填入
var (
	IPInfoToken string = ""
	AiWenToken  string = ""
)
