// Package ppclient enables low-overhead performance profiling and sends them
// to a central repository.
//
// The profiler may skip profiling cycles if it's not given any CPU cycles,
// which can only happen if the scheduler finds spare cycles or if one of the
// busy gouroutines yields control often enough. If your program is constantly
// using the threads it's allowed to (via GOMAXPROCS), the profiler may not
// have the opportunity to run. To mitigate that, try adding runtime.<yield>
// in your main loop(s).
//
// - Count the number of missed profiling collections and report them at
// least. Alternatively, let the CPU profiling be on all the time with low
// frequency and hope that profiling works despite the scheduler issues
// (assuming it relies on OS signals).
//
// - Find the cost of CPU profiling at the currently used 100hz.
package ppclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/nictuku/webpprof/ppcommon"
)

const PrrrURL = "http://localhost:8080/profile"

var (
	CPUProfilingInterval = 30 * time.Second
	CPUProfilingDuration = 5 * time.Second
)
var Profiles = []string{"heap", "block", "goroutine", "threadcreate"}

func createProfile(name string) *ppcommon.Profile {
	buf := new(bytes.Buffer)
	p := pprof.Lookup(name)
	if p == nil {
		return nil
	}
	if err := p.WriteTo(buf, 1); err != nil {
		return nil
	}
	fmt.Println("profile data", buf.String())
	return &ppcommon.Profile{"", name, buf.Bytes(), time.Now()}
}

func profileURL(name string) string {
	u := *prrrURL
	q := u.Query()
	q.Set("p", name)
	u.RawQuery = q.Encode()
	return u.String()
}

func sendProfile(p *ppcommon.Profile) error {
	// TODO: Authentication.
	// TODO: identify the program name and arguments.
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(p)
	if err != nil {
		log.Printf("sendPr enc error %v", err)
		return err
	}
	req, err := http.NewRequest("POST", profileURL(p.Name), buf)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("sendProfile %v error reading response body: %v", p.Name, err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("sendProfile %v error reading response body: %v", p.Name, err)
		return err
	}
	resp.Body.Close()
	log.Println("sendProfile got", string(body))
	return nil
}

// XXX rewrite stuff copied from net/http/pprof.
// better docs.

// Cmdline responds with the running program's
// command line, with arguments separated by NUL bytes.
func cmdline() string {
	return strings.Join(os.Args, "\x00")
}

func sendCPUProfile(buf *bytes.Buffer) {
	// Send the hpprof to the remote server.

	// TODO: Authentication.

	// TODO: identify the program name and arguments.
	fmt.Printf("%v prof result\n", string(buf.String()))
}

// cpuProfile enables CPU profiling for the provided duration and writes the
// resulting pprof-formatted report to w. As of Go 1.1, the CPU profiler
// collects counters of the stack positions at 100hz frequency. An error will
// be reported if, among other reasons, CPU profiling is already enabled when
// this function is called.
func cpuProfile(w io.Writer, duration time.Duration) error {
	if duration.Seconds() == 0 {
		duration = 30 * time.Second
	}

	if err := pprof.StartCPUProfile(w); err != nil {
		return err
	}
	time.Sleep(duration)
	pprof.StopCPUProfile() // Also finishes the write to w.
	return nil
}

func send() {
	for _, name := range Profiles {
		fmt.Println(name)
		if name == "cpu" {
			buf := new(bytes.Buffer)
			if err := cpuProfile(buf, CPUProfilingDuration); err != nil {
				log.Printf("ppclient cpuProfile error: %v", err)
			} else {
				sendCPUProfile(buf)
			}
		} else {
			fmt.Println("creating")
			p := createProfile(name)
			fmt.Println("created")
			if p == nil {
				log.Printf("createProfile %v returned an empty profile", name)
				continue
			}
			fmt.Println("sending")
			sendProfile(p)
		}
	}
}

// profiler runs in the background, waking up occasionally to collect
// performance profiles and send them to prrr.
func profiler() {
	log.Println("profiler started")
	tick := time.Tick(CPUProfilingInterval)
	for {
		send()
		<-tick

	}
}

var prrrURL *url.URL

func Start() (err error) {
	prrrURL, err = url.Parse(PrrrURL)
	if err != nil {
		return err
	}

	go profiler()
	return nil
}
