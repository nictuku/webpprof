webpprof
--------
webpprof is a tool and service that collects and stores profiling data about
Go programs. It's a work in progress and not ready for use yet.

Programs must import a special library to use webpprof. The library enables
low-overhead heap and CPU profiling and pushes them as a best effort to a
central repository. A private repository could also be used.

The service stores the data for a few days or weeks, aggregates information about
different profiles and generates reports ([example heap allocation report](http://1.a.magnets.im/static/heap12.svg)) that can be viewed only by the profiles' owner.

A subset of the data will be aggregated, anonymized and shared publicly - that
will probably include stats about GC overhead and the performance of core
libraries.

Development
-----------

See the project's [Trello board](https://trello.com/b/djCGWcRD/webpprof).


Components
--------

ppclient
  * library for collecting and sending performance profiles

ppserver
  * control access to all uploaded profiles. Only the original uploader can access them.
  * web UI for showing individual performance profiles
  * aggregated/anonymized reports about GC performance 

ppstore  
  * server that receives and stores performance profiles

Webpprof Usage
------
Import the ppclient package and start collecting profiles by adding this snippet somewhere in your program.

```
package main

import (
 "log"
 "github.com/nictuku/webpprof/ppclient"
)

func main() {
 // By default, the ppclient collects and transmits profiles every 1 minute.
 // Use the following to change that to, say, 10 minutes:
 ppclient.CPUProfilingInterval = 10 * time.Minute

 // Start collecting and transmiting profiles, in the background.
 if err := ppclient.Start(); err != nil {
		log.Println("ppclient startup failure:", err)
	}

	// ... Your program
}
```

The data will be transmitted to a webpprof server. You can use the public one (not yet available) or run your own.

webpprof server
-------------------
A central public repository will eventually be available. For now you can run a webpprof server for yourself. The web pprof server stores pprofs in a QL database.

```
$ go build
# Create the database.
$ ql 'CREATE TABLE profiles (user string, name string, content blob, t time);'
$ ./webpprof
```

Test client
------------
There is a test program in ppclient/test that collects and uploads sample profiles. It intentionally leaks memory to make things more interesting.

```
$ cd ppclient/test
$ go run test_client.go
```

Browsing the data
------------------
For now, profiles can only be inspected with the ql tool directly. Install the QL tool:

```
$ go get github.com/cznic/ql/ql
# Confirm that the tool is installed.
$ ql -help

# With the webserver stopped (sorry), inspect the latest heap profile:
$ ql 'SELECT string(content), t FROM profiles WHERE name == "heap" ORDER BY t DESC LIMIT 1;'

```
