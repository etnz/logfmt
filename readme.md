#logfmt [![Travis](https://travis-ci.org/etnz/logfmt.svg?branch=master)](https://travis-ci.org/etnz/logfmt?branch=master) [![GoDoc](https://godoc.org/github.com/etnz/logfmt?status.svg)](https://godoc.org/github.com/etnz/logfmt)

- "logfmt" is a golang package for logging objects in [logfmt](https://brandur.org/logfmt) format: [details](#logfmt)
- "logfmt/reader" is a golang package to parse logfmt streams: [details](./reader/)
- "logfmt/ql" is a golang package that contains the definition of a query language on top of logfmt records. [details](./ql)
- "logfmt/cmd" is a collection of command line utilities for working with logfmt, it complements brandur's own [collection](https://github.com/brandur/hutils). [details](./cmd)






#logfmt


To use the package:

    go get github.com/etnz/logfmt


logfmt is defined in [Brandur's blog](https://brandur.org/logfmt) as a logging format optimal
for easy development, consistency, and good legibility for humans and
computers.

    at=info method=GET path=/ host=mutelight.org status=200

The example above can be generated like:

```go
logfmt.
  S("at", "info").
  S("method",r.Method).
  S("path", r.URL.Path).
  D("status", code ).
  Log()
```

Where 'S', 'Q', 'D' are method mapping the usual 'fmt' verbs.

Key/Value are stored in a map, so it's ok to call S("at", "info") several time, the last one is always right.

When writing the line in the output stream, keys are not generated at *random*. Logs would be hard to read, and not reproducible. 
Instead key/value pairs are sorted in *significance* order: the shortest keys first, then alphabetically

 
    at=info path=/ host=mutelight.org method=GET status=200

It enforce the tendency to keep generic keys short, and specific one longer, the first information you read is the most important.

Usually it is faster than the default "log" package

```go
    r := Rec()
    r.D("timestamp", int(n))
    r.D("at", i)
    r.Q("username", "eric")
    r.K("debug")
    r.Log()
```
At 1389 ns/op  : 720 logs per milliseconds

```go    
     log.Printf("at=%d", i)
     log.Printf("debug")
     log.Printf("username=%q\n", "eric")
```
At 2248 ns/op, that is 445 logs per milliseconds

Because using logfmt, you tend to group together information in less records, but richer.


Record  object is idiomatic:

```go
log:= logfmt.Rec()
defer log.Log()
//set an initial value
r.S("level", "debug")
...
if err != nil{
  r.S("level", "error") // escalate the record
  r.V("err", err)
}

```


See [Examples](https://godoc.org/github.com/etnz/logfmt#pkg-examples) or directly the [godoc](https://godoc.org/github.com/etnz/logfmt) for more details.


# Still on the workbench

- rewriter logs (sort keys, filter out some matches), highlight lowlight some lines or keys
- trim logs (before first match and after last match)