package service

import (
	"fmt"
	"testing"

	"github.com/sjlleo/nexttrace-backend/ipgeo"
)

func TestCheckBone(t *testing.T) {
	CheckChinaBackBone(&ipgeo.IPGeoData{
		IP: "202.97.1.1",
	})
}

func TestGetIPData(t *testing.T) {
	// res := GetIPGeoData("127.0.0.1", 1)
	// log.Println(string(res))
}

func TestPRT(t *testing.T) {
	res, _ := lookupPTRWithContext("1.1.1.1")
	fmt.Println(res)
}
