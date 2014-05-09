package main

import (
	"net/http"

	"github.com/nictuku/webpprof/ppserver"
)

func main() {
	http.HandleFunc("/profile", ppserver.HandlePostProfile)
	http.ListenAndServe(":8080", nil)
}
