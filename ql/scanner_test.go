package ql

import (
	"fmt"
	"strings"
)

func ExampleScanner() {

	src := " .identifier ~ /name.*/ or .load = +12354 and not .user ~ /eric.*/ OR .at > 52.12345  and .elapsed < 375.12ms"

	s := newScanner(strings.NewReader(src))

	for s.ttype != EOF {
		s.Next()
		fmt.Printf("%-10s %q\n", s.ttype.String(), s.token.String())
	}
	//Output:
	// <IDENT>    ".identifier"
	// ~          "~"
	// <REGEXP>   "/name.*/"
	// OR         "or"
	// <IDENT>    ".load"
	// =          "="
	// <NUMBER>   "+12354"
	// AND        "and"
	// NOT        "not"
	// <IDENT>    ".user"
	// ~          "~"
	// <REGEXP>   "/eric.*/"
	// OR         "OR"
	// <IDENT>    ".at"
	// >          ">"
	// <DECIMAL>  "52.12345"
	// AND        "and"
	// <IDENT>    ".elapsed"
	// <          "<"
	// <DURATION> "375.12ms"
	// <EOF>      "\x00"

}

func ExampleScanner_Function() {
	src := " since( .at ) < 1h"

	s := newScanner(strings.NewReader(src))

	for s.ttype != EOF {
		s.Next()
		fmt.Printf("%-10s %q\n", s.ttype.String(), s.token.String())
	}
	//Output:
	// <FUNCTION> "since"
	// (          "("
	// <IDENT>    ".at"
	// )          ")"
	// <          "<"
	// <DURATION> "1h"
	// <EOF>      "\x00"

}

func ExampleScanner_Paren() {
	src := "( .x = 1)"

	s := newScanner(strings.NewReader(src))

	for s.ttype != EOF {
		s.Next()
		fmt.Printf("%-10s %q\n", s.ttype.String(), s.token.String())
	}
	// //Output:
	// (          "("
	// <IDENT>    ".x"
	// =          "="
	// <NUMBER>   "1"
	// )          ")"
	// <EOF>      "\x00"

}
