package wslistener

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sjlleo/nexttrace-backend/dbtools"
	"github.com/sjlleo/nexttrace-backend/ipgeo"
	"github.com/sjlleo/nexttrace-backend/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

type WsConn struct {
	Conn *websocket.Conn
	Mux  sync.Mutex
}

// func getIPSense(res *ipgeo.IPGeoData) {
// 	var aiwenScene string
// 	if ipgeo.HasLocalIPAddr(res.IP) {
// 		return
// 	}
// 	r, err := dbtools.SearchIPScense(res.IP)
// 	if err != nil || r.Scense == "" {
// 		// aiwenScene, _ = ipgeo.AiwenTechScense(res.IP)
// 		// dbtools.AddIPScense(res, aiwenScene)
// 	} else {
// 		aiwenScene = r.Scense
// 	}
// 	res.Domain = aiwenScene
// }

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

func getIPData(mt int, wsconn *WsConn, message []byte, uid int) {
	// 定义返回的数据结构
	var res []byte
	var reply *ipgeo.IPGeoData
	var searchMode bool = false
	var err error
	msg := string(message)
	switch {
	case msg == "ping":
		// 保活应答
		res = []byte("pong")
	case strings.HasPrefix(msg, "FindNeiborHop"):
		msgSlice := strings.Split(msg, "|")
		if msgSlice[1] != "" {
			str, _ := dbtools.FindNeiborHop(msgSlice[1])
			res = []byte(str)
		} else {
			res = []byte(`{"error": {"message": "IPAddr cannot be empty"}}`)
		}

	default:
		searchMode = true
		reply = service.GetIPGeoData(msg, uid)
		res, _ = json.Marshal(reply)
	}
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
	if searchMode {
		service.CheckNewUpdate(reply)
	}
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

		go getIPData(mt, &wsconn, message, uid)
	}
}
