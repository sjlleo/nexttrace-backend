package ipgeo

import (
	"log"
	"net"
)

func aiwenBlacklisted(as string) bool {
	blackList := []string{"58453"}
	for i := 0; i < len(blackList); i++ {
		if blackList[i] == as {
			return true
		}
	}
	return false
}

func ipinfoWhitelisted(as string) bool {
	whitelist := []string{"1299", "9002", "58453"}
	for i := 0; i < len(whitelist); i++ {
		if whitelist[i] == as {
			return true
		}
	}
	return false
}

func HasLocalIPAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP 检测 IP 地址是否是内网地址
// 通过直接对比ip段范围效率更高
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 9) ||
		(ip4[0] == 100 && ip4[1] < 128 && ip4[1] >= 64) ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

func GetIPGeoData(ip string) *IPGeoData {
	if net.ParseIP(ip) == nil {
		// IP 输入非法
		return &IPGeoData{
			IP: ip,
		}
	}

	if HasLocalIPAddr(ip) {
		log.Println("局域网")
		return &IPGeoData{
			IP: ip,
		}
	}

	aiwenData, _ := AiwenTech(ip)
	ipInfoData, _ := IPInfo(ip)

	if (ipinfoWhitelisted(aiwenData.Asnumber) || ipInfoData.Country != "CN") && aiwenData.Asnumber != "45102" {
		// IPInfo 如果不是专业套餐，并不会返回ASN信息，这时候可以用埃文科技的来填补
		if ipInfoData.Asnumber == "" {
			// 不是专业版套餐
			ipInfoData.Asnumber = aiwenData.Asnumber
			ipInfoData.Owner = aiwenData.Owner
			ipInfoData.Isp = aiwenData.Isp
		}
		return ipInfoData
	}
	return aiwenData

	// if regionCode == "" {
	// 	// IPInsight 异常
	// 	ipInfoData, _ := IPInfo(ip)

	// }

	// return aiwenData
}

func GetIPASNDomain(ip string) *IPGeoData {
	ipwhoisData, _ := ipwhois(ip)
	return ipwhoisData
}
