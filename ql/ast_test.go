package ql

import "fmt"

func ExampleFmt() {

	ast := &BinaryExpr{
		Op: OR,
		X: &CompExpr{
			Op: MATCH,
			X: &Literal{
				Kind:  IDENT,
				Value: ".name",
			},
			Y: &Literal{
				Kind:  REGEXP,
				Value: "/eric.*/",
			},
		},
		Y: &CompExpr{
			Op: LT,
			X: &Literal{
				Kind:  IDENT,
				Value: ".load",
			},
			Y: &Literal{
				Kind:  DURATION,
				Value: "35ms",
			},
		},
	}

	fmt.Println(Fmt(ast))
	//Output: .name ~ /eric.*/   OR   .load < 35ms

}
