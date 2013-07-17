webpprof
--------
webpprof is a tool and service that collects and stores profiling data about
Go programs.

Programs must import a special library to use webpprof. The library enables
low-overhead heap and CPU profiling and pushes them as a best effort to a
central repository, called Prrr.

Prrr stores the data for a few days or weeks, aggregates information about
different profiles and generates reports.

A subset of the data will be aggregated, anonymized and shared publicly - that
will probably include stats about GC overhead and the performance of core
libraries.

= Roadmap =

ppclient
  * library for collecting and sending performance profiles

ppserver
  * control access to all uploaded profiles. Only the original uploader can access them.
  * web UI for showing individual performance profiles
  * aggregated/anonymized reports about GC performance 

ppstore  
  * server that receives and stores performance profiles
