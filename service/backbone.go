package service

import (
	"net"
	"strings"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

func CheckBadBackBone(r *ipgeo.IPGeoData) {
	// Google
	if r.Owner == "谷歌公司" {
		if r.City == "Mountain View" || r.City == "山景城" {
			r.Country = "Unknown"
			r.Prov = ""
			r.City = ""
		}
	}

	if r.Asnumber == "13335" {
		if r.City == "San Francisco" || r.City == "旧金山" {
			r.Country = "Anycast"
			r.Prov = ""
			r.City = ""
		}
	}

	if r.Whois == "CLOUDFLARENET" {
		r.Domain = "cloudflare.com"
	}
}

func CheckChinaBackBone(r *ipgeo.IPGeoData) bool {
	if r.Asnumber == "23764" {
		r.Owner = "中国电信"
		r.Domain = "chinatelecomglobal.com "
		r.Domain = r.Domain + " " + strings.Trim(r.Owner, "中国")
		return true
	}

	if r.Owner == "是方电讯股份有限公司" || r.Isp == "是方电讯股份有限公司" || r.Owner == "3Link" || r.Isp == "3Link" {
		r.Domain = "chief.net.tw "
		r.Domain = r.Domain + " " + "是方电讯"
		return true
	}

	if (strings.HasPrefix(r.IP, "210.78.") && r.Asnumber == "") || r.Asnumber == "9929" {
		r.Owner = "中国联通 CUII"
		r.Domain = "chinaunicom.cn "
		r.Domain = r.Domain + " " + strings.Trim(r.Owner, "中国")
		return true
	}

	if strings.HasPrefix(r.IP, "59.43.") {
		r.Owner = "中国电信"
		r.Domain = "chinatelecom.cn "
	}

	if r.Owner == "中国电信" || r.Owner == "中国联通" || r.Owner == "中国移动" || r.Owner == "中国教育网" || r.Owner == "中国科技网" {
		r.Domain = r.Domain + " " + strings.Trim(r.Owner, "中国")
	}

	if r.Owner == "中华电信股份有限公司" || r.Isp == "中华电信股份有限公司" {
		r.Domain = "hinet.net " + " " + "中华电信"
		return false
	}

	if r.Owner == "台湾固网" || r.Isp == "台湾固网" {
		r.Domain = "taiwanmobile.com " + " " + "台湾移动"
		return false
	}

	// switch r.Asnumber {
	// case "20485":
	// 	r.Domain = "ttk.ru  " + "俄罗斯铁通"
	// 	return false
	// case "12389":
	// 	r.Domain = "rt.ru  " + "俄罗斯电信"
	// 	return false
	// case "8359":
	// 	r.Domain = "mts.ru  " + "俄罗斯移动"
	// 	return false
	// case "3491":
	// 	r.Domain = "pccw.com  " + "电讯盈科"
	// 	return false
	// case "17676":
	// 	r.Domain = "bbtec.net  " + "软银"
	// 	return false
	// case "2914":
	// 	r.Domain = "ntt.com  " + "日本电信"
	// 	return false
	// case "4713":
	// 	r.Domain = "ntt.com  " + "日本电信"
	// 	return false
	// case "6453":
	// 	r.Domain = "tatacommunications.com   " + "塔塔通信"
	// 	return false
	// case "9304":
	// 	r.Domain = "hgc.com.hk   " + "和记环球电讯"
	// 	return false
	// }

	//"154.54", "129.250", "2001:218:0:", "2001:1900:", "101.4.", "62.115.", "2001:2034:1:", "2001:668:0:3:ffff:", "2600:80a:", "204.148.", "184.105.", "2001:470:0:", "64.125.", "2001:438:", "63.218.", "2400:8800:", "2001:5a0:", "193.251.", "2001:688:0:", "80.91."
	backboneList := []string{"202.97.0.0/16", "59.43.0.0/16", "219.158.0.0/16", "221.183.0.0/16", "111.5.0.0/16", "218.105.0.0/16", "210.78.0.0/16", "240e:0::/32", "240e:2::/32", "240e::f:0:0:0/80", "240e::1:0:0:0/80", "2409:8080::/32", "2408:8000:2::/48"}
	if r.Country != "中国" || r.Prov == "台湾省" || r.Prov == "香港" || r.Prov == "澳门特别行政区" || r.Prov == "澳门" {
		return false
	}

	for _, v := range backboneList {
		ip := net.ParseIP(r.IP)
		_, tmp, _ := net.ParseCIDR(v)
		if tmp.Contains(ip) {
			// r.Prov = ""
			r.City = ""

		}
	}
	r.District = ""
	return false
}
