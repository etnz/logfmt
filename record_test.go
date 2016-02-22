package logfmt

import (
	"time"

	"os"
)

func ExampleQ() {
	Default = New(os.Stdout)
	Q("key1", "value1").
		Q("key1.detail", "another value").
		Log()

	//Output: key1="value1" key1.detail="another value"
}

func ExampleS() {
	Default = New(os.Stdout)
	S("key1", "ident1").S("key2", "ident2").Log()
	//Output: key1=ident1 key2=ident2
}
func ExampleS_autoquote() {
	Default = New(os.Stdout)
	S("key1", "ident1 ident2").Log() //this would lead to an invalid value
	S("key1", "\"ide nt1\"").Log()   //correctly quoted value is not super quoted
	S("key1", "\"ide nt1").Log()     // invalid value is properly quoted
	//Output:
	// key1="ident1 ident2"
	// key1="ide nt1"
	// key1="\"ide nt1"
}

func ExampleS_Duration() {
	Default = New(os.Stdout)
	V("load", 25*time.Millisecond).D("size", 523).V("cplx", 1+2i).Log()
	//Output: cplx=(1+2i) load=25ms size=523
}

func ExampleT() {
	Default = New(os.Stdout)
	T("key", true).T("debug", false).Log()
	//Output: key=true debug=false
}
func ExampleD() {
	Default = New(os.Stdout)
	D("key", 1).D("ends", 125).Log()
	//Output: key=1 ends=125
}
func ExampleG() {
	Default = New(os.Stdout)
	G("load", 12.3).Log()
	//Output: load=12.3
}
func ExampleK() {
	Default = New(os.Stdout)
	K("debug").K("verbose").Log()
	//Output: debug verbose
}

func ExampleRec() {
	Default = New(os.Stdout)
	//create a record asap, and defer the actual log write.
	r := Rec()
	defer r.Log()

	r.K("debug")
	//somewhere else in the code
	r.D("load", 125)
	//Output: load=125 debug
}
