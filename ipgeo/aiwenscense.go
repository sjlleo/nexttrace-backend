package ipgeo

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tidwall/gjson"
)

func AiwenTechScense(ip string) (string, error) {

	url := "https://api.ipplus360.com/ip/info/v1/scene/?key=" + AiWenToken + "&ip=" + ip

	client := &http.Client{
		// 5 秒超时
		Timeout: 5 * time.Second,
	}
	req, _ := http.NewRequest("GET", url, nil)

	content, err := client.Do(req)
	if err != nil {
		log.Println("埃文科技 请求超时(4s)")
		return "", err
	}

	body, _ := io.ReadAll(content.Body)
	res := gjson.ParseBytes(body)

	res = res.Get("data")
	scene := res.Get("scene").String()

	if scene == "基础设施" {
		scene = "骨干网"
	}

	fmt.Println(scene)
	return scene, nil
}
