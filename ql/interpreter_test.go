package ql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/etnz/logfmt/logreader"
)

var (
	cases = []struct{ rec, query, result string }{
		{
			"user=johndoe",
			".user",
			`"johndoe"`,
		},
		{
			`user=John`,
			`.user`,
			`"John"`,
		},
		{
			`in=120 out=125`,
			`.in < .out`,
			`true`,
		},
		{
			`user=johndoe@mail.com`,
			`.user ~ /john.*/`,
			`true`,
		},
		{
			`user=johndoe@mail.com age=20`,
			`.user ~ /john.*/ and .age < 40`,
			`true`,
		},
		{
			`a=true b=false c=false`,
			`.a OR .b AND .c`,
			`true`,
		},
		{
			`a=true b=false c=false`,
			`.a OR ( .b AND .c )`,
			`true`,
		},
		{
			`a=true b=false c=false`,
			`(.a OR  .b ) AND .c `,
			`false`,
		},
		{
			`a="path/subpath" `,
			`.a ~ /path\/sub.*/`,
			`true`,
		},
	}
)

func TestCases(t *testing.T) {

	for _, c := range cases {

		rec, err := logreader.Parse(c.rec)
		if err != nil {
			t.Fatalf("Invalid record in test %v: %v", c, err)
			return
		}
		q, err := Parse(strings.NewReader(c.query))
		if err != nil {
			t.Fatalf("Invalid query in test %v: %v", c, err)
			return
		}

		result := AsLiteral(Eval(q, rec))
		if result != c.result {
			t.Errorf("Invalid result in test %v: %q instead of %q", c, result, c.result)
		}
	}

}

func ExampleInterpreter() {

	x, err := Parse(strings.NewReader(".a ~ /t.t./ AND .b < 100")) //.a ? and
	if err != nil {
		panic(err)
	}

	a := "toto"
	b := "98"

	y, err := Eval(x, map[string]*string{
		"a": &a,
		"b": &b,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(y)
	//Output: true
}

func ExampleInterpreter_Duration() {

	x, err := Parse(strings.NewReader(".a < 100s"))
	if err != nil {
		panic(err)
	}
	a := "100µs"
	y, err := Eval(x, map[string]*string{"a": &a})
	if err != nil {
		panic(err)
	}
	fmt.Println(y)
	//Output: true
}

// func ExampleInterpreter_Since() {

// 	x, err := Parse(strings.NewReader("since ( .a ) > 1s"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	a := "2016/02/12 23:26:59"
// 	y, err := Eval(x, map[string]*string{"a": &a})
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(y)
// 	//Output: true
// }

// func ExampleInterpreter_EvalTime() {

// 	x, err := Parse(strings.NewReader(".a"))
// 	if err != nil {
// 		panic(err)
// 	}
// 	a := "\"2016/02/14 18:16:00\""
// 	y, err := EvalDuration(x, map[string]*string{"a": &a})
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(y)
// 	//Output: true
// }
