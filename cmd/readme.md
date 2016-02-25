# cmd collection

To install the tools you need to get [Go](http://golang.org) Then to 

    go get github.com/etnz/logfmt/cmd/{lrep}

# lrep 

    go get github.com/etnz/logfmt/cmd/lrep


`lrep` is a simple command line that filters a logfmt stream using query

    $ cat server.log
    
    at=info method=GET path=/ host=mutelight.org fwd="124.133.52.161"
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"
    
   	$ cat server.log | lrep ".method ~ /POST/"
    
    at=info method=POST path=/ host=mutelight.org fwd="124.133.52.161"

