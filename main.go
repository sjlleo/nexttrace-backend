package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/sjlleo/nexttrace-backend/wslistener"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/v2/ipGeoWs", wslistener.Response)
	http.ListenAndServe("0.0.0.0:5000", nil)
}
