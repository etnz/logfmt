#logfmt [![Travis](https://travis-ci.org/etnz/logfmt.svg?branch=master)](https://travis-ci.org/etnz/logfmt.svg?branch=master) [![GoDoc](https://godoc.org/github.com/etnz/logfmt?status.svg)](https://godoc.org/github.com/etnz/logfmt)

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



# Reader

The same package offers a parser to read streams in logfmt, one record at a time.

```go
r := NewReader(src)
for r.HasNext() {
	rec, _ := r.Next()
	fmt.Println(rec)
}
```

# ql
# cmd



A 'rec' beeing a `map[string]*string`


# logfmt/ql query language

'ql' is a simple query langage and interpreter to evaluate a simple expression on top of a logfmt record.


## How to write a query ?

'Key' are evaluated to the key's value. Key name is prefixed by '.'

      record   user=John mail=john@doe.com
      query    .user
      result   John

Comparison is to compare two values

      record   in=120 out=125
      query    .in < .out
      result   true

Comparison available operators are: '<', '=', '>'

Matching: to match a value against a regular expression

      record   user=johndoe@mail.com
      query    .user ~ /john.*/
      result   true

Regular expression literal are delimited by '/' character. To write a '/'  inside the regular expression you need to escape it : '/path\/subpath/'

Logic arithmetic: Comparisons and matchings can be combined using usual boolean arithmetic

      record   user=johndoe@mail.com age=20
      query    .user ~ /john.*/ and .age < 40
      result   true

'AND' operator has priority over 'OR'

`.a OR .b AND .c` is equivalent to `.a  OR  ( .b  AND  .c )`

Space delimiter: logfmt keys can be anything but ' ', therefore key names *must* be delimited by space.

      query    '(.a AND .b)'  is not valid '.b)' is a single key name.
      query    '(.a AND .b )' is valid.

Be careful.  


# lrep 

    go get github.com/etnz/logfmt/cmd/lrep

`lrep` is a simple command line that filters a logfmt stream using query

    $ cat server.log
    
    at=info method=GET path=/ host=mutelight.org fwd="124.133.52.161"
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"
    
   	$ cat server.log | lrep ".method ~ /POST/"
    
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"


# Still on the workbench

- rewriter logs (sort keys, filter out some matches), highlight lowlight some lines or keys
- trim logs (before first match and after last match)