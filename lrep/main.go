package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/etnz/logfmt"
	"github.com/etnz/logfmt/ql"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid number of argument: waiting for one argument: the ql expression, got %v instead", len(os.Args)-1)
		os.Exit(-1)
	}

	x, err := ql.Parse(strings.NewReader(os.Args[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid query:\n%q\n%v\n", os.Args[1], err)
		os.Exit(-2)
	}
	fmt.Printf("filtering out %s\n", ql.Fmt(x))

	reader := logfmt.NewReader(os.Stdin)
	for reader.HasNext() {
		//read the next record
		rec, err := reader.Next()
		//print out the next if it matches
		if err == nil {
			continue
		}

		if match, err := ql.Eval(x, rec); match {
			logfmt.Default.Log(rec)
		} else if err != nil {
			log.Printf("err %v : %v", err, rec)
		}
	}
}
