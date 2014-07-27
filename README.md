webpprof
--------
webpprof is a tool and service that collects and stores profiling data about
Go programs.

Programs must import a special library to use webpprof. The library enables
low-overhead heap and CPU profiling and pushes them as a best effort to a
central repository. A private repository could also be used.

The service stores the data for a few days or weeks, aggregates information about
different profiles and generates reports that can be viewed only by the profiles'
owner.

A subset of the data will be aggregated, anonymized and shared publicly - that
will probably include stats about GC overhead and the performance of core
libraries.

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

Usage
------
Eventual usage will involve installing a reporting library, then viewing reports on the web.

Functionality is limited for now. You can start a web pprof server that stores pprofs in a QL database. There is also a test program in ppclient/test that collects and uploads sample profiles. For now, profiles can only be inspected with the ql tool directly.

Install the QL tool:

```
$ go get github.com/cznic/ql/ql
# Confirm that the tool is installed.
$ ql -help
```

Test server:
```
$ go build
# Create the database.
$ ql 'CREATE TABLE profiles (user string, name string, content blob, t time);'
$ ./webpprof
```

Test client:
```
$ cd ppclient/test
$ go run test_client.go
```
