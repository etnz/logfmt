package reader

import (
	"fmt"
	"strings"
)

func ExampleReader() {

	src := `at=1234578 debug path=/login user="foo@bar.com"
	at=1234589 debug path=/login user="bar@bar.com"
	at=1234599 debug path=/login user="baz@bar.com"
	`

	r := NewReader(strings.NewReader(src))
	for r.HasNext() {
		rec, _ := r.Next()
		fmt.Println(rec)
	}
	//Output:
	// at=1234578 path=/login user=foo@bar.com debug
	// at=1234589 path=/login user=bar@bar.com debug
	// at=1234599 path=/login user=baz@bar.com debug
}
