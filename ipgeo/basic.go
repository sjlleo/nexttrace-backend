package ipgeo

type IPGeoData struct {
	IP          string              `json:"ip"`
	Asnumber    string              `json:"asnumber"`
	Country     string              `json:"country"`
	CountryEN   string              `json:"country_en"`
	CountryCode string              `json:"country_code"`
	Prov        string              `json:"prov"`
	RegionEN    string              `json:"prov_en"`
	City        string              `json:"city"`
	CityEN      string              `json:"city_en"`
	District    string              `json:"district"`
	Owner       string              `json:"owner"`
	Isp         string              `json:"isp"`
	Domain      string              `json:"domain"`
	Whois       string              `json:"whois"`
	Prefix      string              `json:"prefix"`
	Lat         float64             `json:"lat"`
	Lng         float64             `json:"lng"`
	Router      map[string][]string `json:"router"`
	Source      string              `json:"source"`
}

var (
	IPInfoToken    string = ""
	AiWenToken     string = ""
	AiWenV6Token   string = ""
	IPDataToken    string = ""
	IPInsightToken string = ""
)
