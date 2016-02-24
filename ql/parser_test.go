package ql

import (
	"fmt"
	"strings"
)

func ExampleParser() {
	x, err := Parse(strings.NewReader(".name ~ /eric.*/  OR  .load < 35ms AND .debug ? OR .name ~ /.*/"))
	if err != nil {
		fmt.Printf("Error :%v", err)
		panic(err)
	}

	fmt.Println(Fmt(x))
	//Output: .name ~ /eric.*/   OR   .load < 35ms   AND   .debug ?   OR   .name ~ /.*/

}

func ExampleParser_Paren() {
	x, err := Parse(strings.NewReader("(.x = 12)"))
	if err != nil {
		fmt.Printf("Error :%v", err)
		panic(err)
	}

	fmt.Println(Fmt(x))
	//Output: (.x = 12)

}
func ExampleParser_Func() {
	x, err := Parse(strings.NewReader("since( .at ) > 1s"))
	if err != nil {
		fmt.Printf("Error :%v", err)
		panic(err)
	}

	fmt.Println(Fmt(x))
	//Output: since( .at ) > 1s

}
