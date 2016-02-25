package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/etnz/logfmt"
	"github.com/etnz/logfmt/logreader"
	"github.com/etnz/logfmt/ql"
)

var (
	debug = flag.Bool("v", false, "set to true to print out extra log (lrep self logs) all with the lrep attribute")
	help  = flag.Bool("h", false, "display some help")
)

func main() {
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if len(flag.Args()) != 1 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "Expecting a single query expression, got %v arguments instead.\n", flag.Args())
		os.Exit(-1)
	}

	q := flag.Arg(0)
	cmd := os.Args[0]

	// start the job by parsing the ql expression
	x, err := ql.Parse(strings.NewReader(q))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid query:\n    %q\n    %v\n", q, err)
		os.Exit(-2)
	}
	if *debug {
		logfmt.
			K(cmd).
			Q("query", ql.Fmt(x)).
			Log()
	}

	//and now read from the stdin for logfmt
	reader := logreader.New(os.Stdin)
	for reader.HasNext() {
		//read the next record
		rec, err := reader.Next()
		if err != nil {
			if *debug {

				logfmt.
					K(cmd).
					K("read-error").
					Q("error", err.Error()).
					Log()
			}
			continue
		}

		//print out the next if it matches
		match, err := ql.AsBool(ql.Eval(x, rec))
		switch {

		case err == nil && match:
			//simply log it
			logfmt.Default.Log(rec)

		case err != nil && *debug:
			// error in debug mode: simply print the "faulty" log record and the error

			logfmt.Default.Log(rec)
			logfmt.
				K(cmd).
				K("runtime-error").
				Q("error", err.Error()).
				Log()
		}
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, ql.WriteAQuery)

}
