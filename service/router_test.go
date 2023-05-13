package service

import (
	"log"
	"testing"
)

func TestRouter(t *testing.T) {
	GetBGPToolsData("2408:8000:1100:21d::3")
}

func TestRouterSave(t *testing.T) {
	m := make(map[string][]string)
	SaveData("", "", m)
}

func TestRouterBasic(t *testing.T) {
	asNumber, prefix := GetBasic("2a06:a005:28f2::1")
	log.Println(asNumber, prefix)

}
