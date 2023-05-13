package service

import (
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/c-robinson/iplib"
	"github.com/sjlleo/nexttrace-backend/dbtools"
)

func SaveData(ip string, prefix string, m map[string][]string) {
	if len(m) == 0 {
		log.Println("异常数据")
		return
	}
	res := strings.Split(prefix, "/")
	dbtools.AddRouter(res[0], res[1], m)
}

func SearchData(ip string) (string, map[string][]string) {
	m := make(map[string][]string)
	r := dbtools.SearchRouter(ip)

	if len(r) == 0 {
		return "", nil
	}
	for _, v := range r {
		m[v.Fromasnumber] = append(m[v.Fromasnumber], v.Asnumber)
	}
	prefix, _ := strconv.Atoi(r[0].Prefix)
	n := iplib.NewNet(net.ParseIP(ip), prefix)
	return n.FirstAddress().String() + "/" + r[0].Prefix, m
}

func GetBGPToolsData(ip string) (map[string][]string, string) {
	if prefix, data := SearchData(ip); prefix != "" {
		return data, prefix
	} else {
		asnumber, prefix := GetBasic(ip)
		log.Println(asnumber)

		if prefix == "<nil>" {
			return nil, ""
		} else {
			return nil, prefix
		}

		// m := GetPath(asNumber, prefix)

		// prefix = strings.Replace(prefix, "_", "/", -1)
		// SaveData(ip, prefix, m)
		// return m, prefix
	}
}

func GetBasic(ip string) (string, string) {
	conn, _ := net.Dial("tcp", "bgp.tools:43")

	conn.Write([]byte(ip + "\r\n"))
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return "", ""
	}
	str := string(buf[:n])

	str = strings.Split(str, "\n")[1]
	str = strings.Replace(str, " ", "", -1)
	asNumber := strings.Split(str, "|")[0]
	prefix := strings.Split(str, "|")[2]
	prefix = strings.Replace(prefix, "/", "_", -1)
	return asNumber, prefix
}

func RegexpUtil(originAsn string, rep string, content string) map[string][]string {
	routeMap := make(map[string][]string)
	compile := regexp.MustCompile(rep)
	submatch := compile.FindAllStringSubmatch(content, -1)
	for _, value := range submatch {
		routeMap[value[1]] = append(routeMap[value[1]], value[2])
	}

	return routeMap
}

func BFS(r *map[string][]string, node string, i *int, end_sign bool) {
	routeMap := *r
	if routeMap[node] == nil {
		if *i == 0 {
			fmt.Print(" → ", node)
			if end_sign {
				fmt.Print(" ]")
			}
		} else {
			fmt.Print(" / ", node)
		}
		return
	}
	if *i == 0 {
		fmt.Print(" → ", node)
	} else {
		fmt.Print(" / [ ", node)
		end_sign = true
	}
	for j := 0; j < len(routeMap[node]); j++ {
		BFS(r, routeMap[node][j], &j, end_sign)
	}

	if len(routeMap[node]) != 1 {
		fmt.Println()
	}

}
