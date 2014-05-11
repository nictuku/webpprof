package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/nictuku/webpprof/ppclient"
)

func main() {

	if err := ppclient.Start(); err != nil {
		log.Fatalln(err)
	}
	f, _ := os.Open("/dev/null")
	for i := 0; ; i++ {
		h := md5.New()
		io.WriteString(h, "The fog is getting thicker!")
		io.WriteString(h, "And Leon's getting laaarger!")
		fmt.Fprintf(f, "%x", h.Sum(nil))
		runtime.Gosched()
	}
}
