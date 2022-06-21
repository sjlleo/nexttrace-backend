package ipgeo

import "strings"

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

func GetIPGeoData(ip string) *IPGeoData {
	aiwenData, _ := AiwenTech(ip)
	ipInsightData, regionCode, _ := IPInsight(ip)

	if ipInsightData.Country != "Reserved" && ipInsightData.Country != "Loopback Address" && regionCode != "CN" {
		// 不在国内

		if ipInsightData.Source == "NA" && aiwenData.Prov != "" && aiwenBlacklisted(aiwenData.Asnumber) {
			// 埃文有数据，IPInsight 没有
			return aiwenData
		} else {
			// 默认优先使用 IPInsight
			if !ipinfoWhitelisted(aiwenData.Asnumber) && ipInsightData.Source != "NA" && !strings.Contains(ipInsightData.Country, "Area") {
				ipInsightData.Asnumber = aiwenData.Asnumber
				ipInsightData.Owner = aiwenData.Owner
				ipInsightData.Isp = aiwenData.Isp
				return ipInsightData
			}
			// 有些 ASN 只有 IPInfo 的库稍微准一些，必须使用IPInfo，又或者当 IPInsight 以及埃文科技都没有数据时
			ipInfoData, _ := IPInfo(ip)
			// IPInfo 如果不是专业套餐，并不会返回ASN信息，这时候可以用埃文科技的来填补
			if ipInfoData.Asnumber == "" {
				// 不是专业版套餐
				ipInfoData.Asnumber = aiwenData.Asnumber
				ipInfoData.Owner = aiwenData.Owner
				ipInfoData.Isp = aiwenData.Isp
			}
			return ipInfoData
		}
	}

	return aiwenData
}

func GetIPASNDomain(ip string) *IPGeoData {
	ipwhoisData, _ := ipwhois(ip)
	return ipwhoisData
}
