package dbtools

import (
	"encoding/json"
	"net"
	"time"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

type History_information struct {
	Id          uint
	Uid         uint
	Ip          string
	Asnumber    string
	Country     string
	Prov        string
	City        string
	District    string
	Owner       string
	Isp         string
	Domain      string
	RequestDate time.Time
}

const (
	Before_Hop_Type = 0
	Next_Hop_Type   = 1
)

type Result struct {
	BeforeHop map[string]*Hop
	NextHop   map[string]*Hop
}

type Hop struct {
	Address *ipgeo.IPGeoData `json:"address"`
	Count   uint             `json:"count"`
}

func AddHistory(data *ipgeo.IPGeoData, uid int) error {

	db := GetDB()

	c := History_information{
		Ip:          data.IP,
		Uid:         uint(uid),
		Asnumber:    data.Asnumber,
		Country:     data.Country,
		Prov:        data.Prov,
		City:        data.City,
		District:    data.District,
		Owner:       data.Owner,
		Isp:         data.Isp,
		Domain:      data.Domain,
		RequestDate: time.Now(),
	}

	res := db.Create(&c)
	return res.Error
}

func FindNeiborHop(ip string) (string, error) {
	if ipgeo.HasLocalIPAddr(ip) {
		byte, _ := json.Marshal(&Result{})
		return string(byte), nil
	}
	db := GetDB()
	h := []History_information{}
	db.Where("ip = ?", ip).Find(&h)

	var valid_id []uint
	var valid_id_user []uint

	for i := 0; i < len(h); i++ {
		// 如果是第一条记录就直接录入
		if i == 0 {
			valid_id = append(valid_id, h[i].Id)
			valid_id_user = append(valid_id_user, h[i].Uid)
		} else {
			// 如果上条记录 ID 临近，则视为无效记录
			if h[i].Id > h[i-1].Id+5 {
				valid_id = append(valid_id, h[i].Id)
				valid_id_user = append(valid_id_user, h[i].Uid)
			}
		}

	}
	// 查找相邻记录
	index := 0
	highendFlag := false
	res := Result{
		BeforeHop: make(map[string]*Hop),
		NextHop:   make(map[string]*Hop),
	}
	for _, v := range valid_id {
		// log.Println("ID = ", v)
		db.Where("id > ? - 5 && id < ? + 5", v, v).Where("uid = ?", valid_id_user[index]).Find(&h)
		// log.Println(h)

		j := 0
		for _, v2 := range h {
			// TODO: 对分查找
			// 查找上沿
			if v2.Ip == ip {
				// 如果 j = 0 单独判断
				if j == 0 {
					// 不存在上沿
					// log.Println("不存在上沿")
					highendFlag = true
					j++
					continue
				}

				if h[j-1].Ip == ip {
					// 重复结果
					j++
					continue
				}
				// 找到后退回
				if v2.RequestDate.Unix()-h[j-1].RequestDate.Unix() < 5 && !ipgeo.HasLocalIPAddr(h[j-1].Ip) {
					if _, created := res.BeforeHop[h[j-1].Ip]; created {
						res.BeforeHop[h[j-1].Ip].Count++
						res.BeforeHop[h[j-1].Ip].Address = searchIPGeo(h[j-1].Ip)
					} else {
						res.BeforeHop[h[j-1].Ip] = &Hop{
							Address: searchIPGeo(h[j-1].Ip),
							Count:   1,
						}
					}

					// log.Println("上一跳IP：" + h[j-1].Ip)
				}
				highendFlag = true
			}
			// 查找下沿
			if highendFlag && v2.Ip != ip {
				if j > 1 && v2.RequestDate.Unix()-h[j-1].RequestDate.Unix() < 5 && !ipgeo.HasLocalIPAddr(v2.Ip) {
					if _, created := res.NextHop[v2.Ip]; created {
						res.NextHop[v2.Ip].Count++
						res.NextHop[v2.Ip].Address = searchIPGeo(v2.Ip)
					} else {
						res.NextHop[v2.Ip] = &Hop{
							Address: searchIPGeo(v2.Ip),
							Count:   1,
						}
					}
					// log.Println("下一跳IP：" + v2.Ip)
				}
				break
			}
			j++
		}
		highendFlag = false
		index++
	}
	byte, _ := json.Marshal(res)
	return string(byte), nil
}

func searchIPGeo(ip string) *ipgeo.IPGeoData {
	if ipn := net.ParseIP(ip); ipn != nil && ipn.To4() == nil {
		cache, _ := SearchIPv6(ip)
		return &ipgeo.IPGeoData{
			IP:       ip,
			Asnumber: cache.Asnumber,
			Country:  cache.Country,
			Prov:     cache.Prov,
			City:     cache.City,
			District: cache.District,
			Owner:    cache.Owner,
			Isp:      cache.Isp,
		}
	} else {
		cache, _ := SearchIP(ip)
		return &ipgeo.IPGeoData{
			IP:       ip,
			Asnumber: cache.Asnumber,
			Country:  cache.Country,
			Prov:     cache.Prov,
			City:     cache.City,
			District: cache.District,
			Owner:    cache.Owner,
			Isp:      cache.Isp,
		}
	}
}
