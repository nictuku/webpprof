package main

import (
	"crypto/md5"
	"fmt"
	"github.com/nictuku/webpprof/ppclient"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	//_ "github.com/nictuku/webpprof/ppserver"
)

func main() {

	if false {
		go func() {
			log.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	}
	if err := ppclient.Start(); err != nil {
		log.Fatalln(err)
	}
	i := 0
	f, _ := os.Open("/dev/null")
	for {
		i++
		h := md5.New()
		io.WriteString(h, "The fog is getting thicker!")
		io.WriteString(h, "And Leon's getting laaarger!")
		fmt.Fprintf(f, "%x", h.Sum(nil))
		runtime.Gosched()
	}
}
