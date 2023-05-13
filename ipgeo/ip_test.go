package ipgeo

import (
	"log"
	"testing"
)

func TestIPInfo(t *testing.T) {
	log.Println(GetIPGeoData("68.142.82.41"))
}

// func TestCTBackbone(t *testing.T) {
// 	ZXINC()
// 	// fmt.Println(AiwenTechv6("240e::1:31:51:5303"))
// }
