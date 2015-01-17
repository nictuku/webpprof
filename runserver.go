package main

import (
	"net/http"

	"github.com/nictuku/webpprof/ppstore"
)

func main() {
	http.HandleFunc("/write", ppstore.HandlePostProfile)
	http.HandleFunc("/profile", ppstore.HandleAuthReadProfile)
	http.ListenAndServe(":8080", nil)
}
