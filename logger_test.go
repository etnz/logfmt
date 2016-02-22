package logfmt

import (
	"bytes"
	"log"
	"testing"
	"time"
)

func BenchmarkLogger(b *testing.B) {
	// Benchmark a small log with three attribute
	var buf bytes.Buffer
	Default = New(&buf)
	for i := 0; i < b.N; i++ {
		n := time.Now().Unix()
		D("timestamp", int(n)).D("at", i).Q("username", "eric").K("debug").Log()
	}
}

func BenchmarkDefaultLogger(b *testing.B) {
	// Benchmark a small log with three attribute the usual way
	var buf bytes.Buffer
	log.SetOutput(&buf)
	for i := 0; i < b.N; i++ {
		//log.Printf("at=%d debug username=%q\n", int64(i), "eric")
		log.Printf("at=%d", i)
		log.Printf("debug")
		log.Printf("username=%q\n", "eric")
	}
}
