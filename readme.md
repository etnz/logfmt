[![Travis](https://travis-ci.org/etnz/logfmt.svg?branch=master)](https://travis-ci.org/etnz/logfmt.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/etnz/logfmt?status.svg)](https://godoc.org/github.com/etnz/logfmt)

golang package 'logfmt' implements a logfmt Record reader and writer.


To write logfmt





See [Examples](https://godoc.org/github.com/etnz/logfmt#pkg-examples) or directly the [godoc](https://godoc.org/github.com/etnz/logfmt) for more details.


# logfmt Writer

To use the package:

    go get github.com/etnz/logfmt


logfmt is defined in [Brandur's blog](https://brandur.org/logfmt) as a logging format optimal
for easy development, consistency, and good legibility for humans and
computers.

    at=info method=GET path=/ host=mutelight.org fwd="124.133.52.161"
    dyno=web.2 connect=4ms service=8ms status=200 bytes=1653

The example above can be generated like:

```go
logfmt.S("at", "info").
  S("method",r.Method).
  S("path", r.URL.Path).
  Q("fwd", "124.133.52.161").
  V("service", time.Since(begin) )
  D("status", code ).
  Log()
```

Where 'S', 'Q', 'D' are the usual 'fmt' verbs.

The keys are not generated at random, but in significance order: the shortest key is considered to be the most generic one, and then, come first.

See [Examples](https://godoc.org/github.com/etnz/logfmt#pkg-examples) or directly the [godoc](https://godoc.org/github.com/etnz/logfmt) for more details.

# logfmt Reader

the same package offers a parser to read streams in logfmt, one record at a time.

```go
r := NewReader(src)
for r.HasNext() {
	rec, _ := r.Next()
	fmt.Println(rec)
}
```

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

      '.a OR .b AND .c' is equivalent to '.a  OR  ( .b  AND  .c )'

Space delimiter: logfmt keys can be anything but ' ', therefore key names *must* be delimited by space.

      query    '(.a AND .b)'  is not valid '.b)' is a single key name.
      query    '(.a AND .b )' is valid.

Be careful.  


# lrep 

    go get github.com/etnz/logfmt/cmd/lrep

`lrep` is a simple command line that filter a logfmt stream using query

    $ cat server.log
    
    at=info method=GET path=/ host=mutelight.org fwd="124.133.52.161"
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"
    
   	$ cat server.log | lrep ".method ~ /POST/"
    
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"


# Still on the workbench

- rewriter logs (sort keys, filter out some matches), highlight lowlight some lines or keys
- trim logs (before first match and after last match)