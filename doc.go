// Package logfmt implements a logfmt Record reader and writer.
//
// Definition
//
// logfmt is defined in https://brandur.org/logfmt as a logging format optimal
// for easy development, consistency, and good legibility for humans and
// computers.
//
//    at=info method=GET path=/ host=mutelight.org fwd="124.133.52.161"
//    dyno=web.2 connect=4ms service=8ms status=200 bytes=1653
//
// A Record is a simple line in logfmt. It is a set of attribute that can have
// or have not a value.
//
// The value is either a valid identifier, or a quoted string.
//
// Reader
//
// Package logfmt can be used to parse logfmt compliant logs.
//
// Reader is a simple interface with two methods to read a stream as Records.
//
// It parses any io.Reader and return a Record each time.
//
//    r := NewReader(os.Stdin)
//    for r.HasNext() {
//        rec, _ := r.Next()
//        fmt.Println(rec)
//    }
//
//
// Logging API
//
// Package logfmt can be used as a simple tool to write logfmt logs.
//
//
// Record can be created using a fluent API based on the usual 'fmt' verbs
//
//    Q for %q : quoted string
//    S for %s : identifier
//    T for %t : boolean
//    D for %d : decimal base integer
//    G for %g : float
//    V for %v : interface{} uses go value
//    K        : just insert the key (no values)
//
// There are basically only one new "words" to learn: 'K' stands for Key only.
//
// Records are then created using a sequence of call to one of these functions.
//
//    Q("user", username).D("retry", retryCount).K("debug").Log()
//
// See Examples for more details.
//
// The Log method prints the Record in a single line, sorting attributes in
// 'significance' order. It means that general keys (short name) come first then
// specific attributes (long name).
//
// This Log method is fitted for 'defer'.
package logfmt
