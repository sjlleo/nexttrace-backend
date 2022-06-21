package wslistener

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sjlleo/nexttrace-backend/dbtools"
	"github.com/sjlleo/nexttrace-backend/ipgeo"
	"golang.org/x/sync/semaphore"
)

var upgrader = websocket.Upgrader{} // use default options

type WsConn struct {
	Conn *websocket.Conn
	Mux  sync.Mutex
}

func getASNData(res *ipgeo.IPGeoData) {
	asnData, err := dbtools.SearchASNData(res.Asnumber)
	// 如果 ASN 不存在
	if asnData.Asn == 0 || err != nil {
		d := ipgeo.GetIPASNDomain(res.IP)
		log.Println(d)
		res.Domain = d.Domain
		// 添加到数据库里面
		dbtools.AddASNData(d)
	} else {
		// 存在的话直接放入
		res.Domain = asnData.Domain
	}
}

func getIPData(mt int, wsconn *WsConn, message []byte, uid int) {
	// 定义返回的数据结构
	var reply *ipgeo.IPGeoData
	ip := string(message)
	// DB 处理
	cache, err := dbtools.SearchIP(ip)
	// DB 没有数据
	if err != nil {
		// 向 IP API 请求数据
		reply = ipgeo.GetIPGeoData(ip)
		log.Println(reply)
		// 向 DB 写入经过整合后的数据
		err := dbtools.AddIP(reply, uid)

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
			Owner:    cache.Owner,
			Isp:      cache.Isp,
		}
	}
	getASNData(reply)
	res, _ := json.Marshal(reply)
	// 返回流
	wsconn.Mux.Lock()
	// websocket.Conn 不支持多协程同时向一个链接写数据，需要互斥锁保护，否则在高并发的请求下容易触发 panic
	// 设置写超时，如果10秒客户端都没有收到数据，则超时丢弃
	wsconn.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	err = wsconn.Conn.WriteMessage(mt, res)
	wsconn.Mux.Unlock()
	if err != nil {
		log.Println("write:", err)
	}
	wsconn.Conn.SetWriteDeadline(time.Time{})
}

func getIPAdress(r *http.Request) string {
	return r.Header.Get("X-Real-IP")
}

func saveUserIP(ip string) (int, error) {
	uid, err := dbtools.AddUsers(ip)
	if err != nil {
		return -1, err
	}
	return uid, nil
}

func Response(w http.ResponseWriter, r *http.Request) {
	// 用户编号
	var uid int
	c, err := upgrader.Upgrade(w, r, nil)
	wsconn := WsConn{
		Conn: c,
	}

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	go func() {
		userConnectingIP := getIPAdress(r)
		uid, err = dbtools.SearchUsers(userConnectingIP)
		if err != nil {
			uid, err = saveUserIP(userConnectingIP)
			if err != nil {
				log.Println("数据库插入用户IP失败")
			}
		}
	}()

	for {
		// 接收部分
		mt, message, err := wsconn.Conn.ReadMessage()
		if err != nil {
			// log.Println("read:", err)
			break
		}
		// 处理部分

		// 信号量设定，控制每个连接的总子协程数不高于20个，防止过多请求涌入导致服务器资源耗尽
		semaphore.NewWeighted(20)
		// 开启子协程处理 IP 地理位置数据
		go getIPData(mt, &wsconn, message, uid)
	}
}
