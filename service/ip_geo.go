package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sjlleo/nexttrace-backend/dbtools"
	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

func getASNData(res *ipgeo.IPGeoData) {
	asnData, err := dbtools.SearchASNData(res.Asnumber)
	ip := net.ParseIP(res.IP)
	// 如果 ASN 不存在
	if asnData.Asn == 0 || err != nil {

	} else {
		// 存在的话直接放入
		if ip.To4() == nil {
			res.Domain = asnData.Website
			return
		}
		if asnData.Website != "" {
			res.Domain = asnData.Website + " " + res.Domain
		} else {
			res.Domain = res.Owner + " " + res.Domain
		}

	}
}

func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

func CheckNewUpdate(r *ipgeo.IPGeoData) error {

	if ptr, err := lookupPTRWithContext(r.IP); err == nil {
		log.Println(ptr)
		// 查询 IP 地理位置
		if iata_res, err := dbtools.GetIATAGeo(ptr); err == nil {
			r.Country = iata_res.Country
			r.Prov = iata_res.Prov
			r.City = iata_res.City
		}
	}
	log.Println(r)
	// IP - Prefix - ASN - Country - Province - City - District
	prefix_tmp, _ := strconv.Atoi(r.Prefix)
	_, ipType := ParseIP(r.IP)
	switch ipType {
	case 4:
		if err := dbtools.UpdateIP(r.IP, prefix_tmp, r.Asnumber, r.Country, r.Prov, r.City, r.District); err != nil {
			return err
		}
	case 6:
		if err := dbtools.UpdateIPv6(r.IP, prefix_tmp, r.Asnumber, r.Country, r.Prov, r.City, r.District); err != nil {
			return err
		}
	}
	return nil
}

func lookupPTRWithContext(ip string) (string, error) {
	// 手动构造 Resolver 以定制化 DNS 服务器 IP 等参数
	r := &net.Resolver{
		// 尽管编译器已经禁用 Cgo，这里以防万一，保证无论何种编译环境下都能优先使用 Pure-Go，构造详见 lookup.go 源码
		PreferGo: true,
	}

	r.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{
			Timeout: 1000 * time.Millisecond,
		}
		// 见文档 - Dial uses context.Background internally; to specify the context, use DialContext.
		return d.DialContext(ctx, "udp", "8.8.8.8:53")
	}

	ptrs, err := r.LookupAddr(context.Background(), ip)
	if err != nil {
		return "", err
	}

	if len(ptrs) == 0 {
		return "", fmt.Errorf("no PTR records found for %s", ip)
	}

	return ptrs[0], nil
}

func GetIPGeoData(msg string, uid int) *ipgeo.IPGeoData {
	var reply *ipgeo.IPGeoData

	// 判断 msg 输入的 IP 信息是否合法
	if ipValid := net.ParseIP(msg); ipValid != nil {
		// IP 合法
		ip := msg
		if ipgeo.HasLocalIPAddr(ip) {
			reply = &ipgeo.IPGeoData{
				IP:      ip,
				Country: "LAN Address",
			}
			return reply
		} else {
			// 先判断 rDNS 是否包含地理位置信息

			// IPv6 特殊处理
			if ipn := net.ParseIP(ip); ipn != nil && ipn.To4() == nil {
				cache, err := dbtools.SearchIPv6(ip)
				if err != nil {
					// 向 IP API 请求数据
					reply, _ = ipgeo.AiwenTechv6(ip)
					if ptr, err := lookupPTRWithContext(msg); err == nil {
						// 查询 IP 地理位置
						if iata_res, err := dbtools.GetIATAGeo(ptr); err == nil {
							reply.Country = iata_res.Country
							reply.Prov = iata_res.Prov
							reply.City = iata_res.City
						}
					}
					asNumber, _ := GetBasic(ip)
					if asNumber != "0" && asNumber != "" {
						reply.Asnumber = asNumber
					}
					// 向 DB 写入经过整合后的数据
					err := dbtools.AddIPv6(reply, uid)

					if err != nil {
						log.Println(err)
					}
				} else {
					// 数据库已有记录，将 cache 和 ipv4_asn 整合为 IPGeoData
					reply = &ipgeo.IPGeoData{
						IP:       ip,
						Asnumber: cache.Asnumber,
						Country:  cache.Country,
						Prov:     cache.Prov,
						City:     cache.City,
						District: cache.District,
						Owner:    cache.Isp,
						Isp:      cache.Owner,
					}
				}
			} else {
				// IPv4 部分
				// DB 处理
				cache, err := dbtools.SearchIP(ip)
				// DB 没有数据
				if err != nil {
					// 向 IP API 请求数据
					reply = ipgeo.GetIPGeoData(ip)
					log.Println(reply)
					// 获取 rDNS 信息

					if ptr, err := lookupPTRWithContext(msg); err == nil {
						// 查询 IP 地理位置
						if iata_res, err := dbtools.GetIATAGeo(ptr); err == nil {
							reply.Country = iata_res.Country
							reply.Prov = iata_res.Prov
							reply.City = iata_res.City
						}
					}
					asNumber, _ := GetBasic(ip)
					if asNumber != "0" && asNumber != "" {
						reply.Asnumber = asNumber
					}
					// 向 DB 写入经过整合后的数据
					err = dbtools.AddIP(reply, uid)

					if err != nil {
						log.Println(err)
					}
				} else {
					// 数据库已有记录，将 cache 和 ipv4_asn 整合为 IPGeoData
					reply = &ipgeo.IPGeoData{
						IP:       ip,
						Asnumber: cache.Asnumber,
						Country:  cache.Country,
						Prov:     cache.Prov,
						City:     cache.City,
						District: cache.District,
						Owner:    cache.Isp,
						Isp:      cache.Owner,
					}
				}

				// 查找 ipWhois 信息
				whois, err := dbtools.SearchIPWhois(ip)
				if err != nil {
					// 无数据
					// if res, err := ipgeo.AiwenIPWhois(ip); err == nil {
					// 	if res.Netname != "IANA-BLOCK" {
					// 		reply.Whois = res.Netname
					// 	}
					// 	dbtools.AddIPWhois(res)
					// }
				} else {
					// 有数据
					reply.Whois = whois.Netname
				}

				// 仅限中国内地开启
				// 	if reply.Country == "中国" && reply.Prov != "香港" {
				// 		getIPSense(reply)
				// 	}
			}
		}
		getASNData(reply)
		// log.Println(reply)
		reply.Source = "LeoMoeAPI"
	}
	if reply != nil {
		mapData, prefix := GetBGPToolsData(reply.IP)
		reply.Router = mapData
		reply.Prefix = prefix

		// 将用户的查询记录插入数据库
		dbtools.AddHistory(reply, uid)

		CheckChinaBackBone(reply)
		GetLocation(reply)
		CheckBadBackBone(reply)

	} else {
		reply = &ipgeo.IPGeoData{
			IP:      msg,
			Country: "IPAddr is illegal",
		}
	}
	return reply
}
