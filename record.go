package logfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

//Record contains a single logfmt line.
type Record map[string]*string

// Log this Record now into  the 'Default' Logger
func (rec Record) Log() { Default.Log(rec) }

// Rec returns a newly created, empty Record
func Rec() *Record { rec := make(Record); return &rec }

// Q creates a new Record, call Q method, and return the Record
func Q(key, val string) *Record { return Rec().Q(key, val) }

// S creates a Record, call S method, and return the Record
func S(key, val string) *Record { return Rec().S(key, val) }

// D creates a Record, call D method, and return the Record
func D(key string, val int) *Record { return Rec().D(key, val) }

// T creates a Record, call T method, and return the Record
func T(key string, val bool) *Record { return Rec().T(key, val) }

// G creates a Record, call G method, and return the Record
func G(key string, val float64) *Record { return Rec().G(key, val) }

// K creates a Record, call K method, and return the Record
func K(key string) *Record { return Rec().K(key) }

// V creates a Record,call V method, and return the Record
func V(key string, value interface{}) *Record { return Rec().V(key, value) }

// set stores an attribute and return the record pointer
func (rec *Record) set(key string, val *string) *Record { (*rec)[key] = val; return rec }

// Q inserts a quoted string attribute like `key="val"`
func (rec *Record) Q(key, val string) *Record {
	x := strconv.Quote(val)
	return rec.set(key, &x)
}

// S insert an identifier attribute `key=val`.
//
// 'val' format is checked and quoted if needed.
func (rec *Record) S(key, val string) *Record {

	//check the validity of val
	var eos rune
	s := newScannerS(val)
	if r := s.Read(); r == '"' { //it will be a string
		s.Str()
		end := s.Read() //read the " //it must be a quote
		if end != '"' {
			return Q(key, val) //quote it
		}
		//ok it's valid str
		eos = s.Read() // read one more (must be eof)
	} else {
		s.Unread()     //unread the first char
		s.Identifier() //read the identifier
		eos = s.Read() // read the next, must be eof
	}
	if eos != eof {
		//there are extra stuff, this is not a valid value
		return Q(key, val)
	}
	return rec.set(key, &val)
}

// D insert an integer attribute `key=12`
func (rec *Record) D(key string, val int) *Record {
	x := strconv.FormatInt(int64(val), 10)
	return rec.set(key, &x)
}

// T insert a boolean attribute `key=false`
func (rec *Record) T(key string, val bool) *Record {
	x := strconv.FormatBool(val)
	return rec.set(key, &x)
}

// G insert a float attribute `key=12.3`
func (rec *Record) G(key string, val float64) *Record {
	x := strconv.FormatFloat(val, 'g', -1, 64)
	return rec.set(key, &x)
}

// K insert a key only attribute `debug verbose`
func (rec *Record) K(key string) *Record { (*rec)[key] = nil; return rec }

// V insert value using fmt.Printf verb "%v".
//
// The result is quoted if necessary
func (rec *Record) V(key string, value interface{}) *Record {
	//we use S to protect the value
	return rec.S(key, fmt.Sprint(value))
}

// String format the current record as a string
func (rec Record) String() string {

	// I can't reuse the fastLog, because this one need to be threadsafe
	var buffer bytes.Buffer
	buf := bufio.NewWriter(&buffer)

	L := len(rec)
	keyvals = &fastKeySorter{
		keys: make([]string, 0, L),
		vals: make([]*string, 0, L),
	}
	logTo(buf, rec, keyvals)
	buf.Flush()
	return buffer.String()
}

//logTo log a record into the buf, using the fastKeySorter struct to fill in
// the fastKeySorter can be reused (for speeding up the process)
func logTo(buf *bufio.Writer, rec Record, sorter *fastKeySorter) {
	L := len(rec)
	for k, v := range rec {
		keyvals.keys = append(keyvals.keys, k)
		keyvals.vals = append(keyvals.vals, v)
	}
	sort.Sort(keyvals)

	for i := 0; i < L; i++ {
		key, val := keyvals.keys[i], keyvals.vals[i]

		if i > 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(key)
		if val != nil {
			buf.WriteRune('=')
			buf.WriteString(*val)
		}
	}
}

type fastKeySorter struct {
	keys []string
	vals []*string
}

func (s *fastKeySorter) Len() int { return len(s.keys) }
func (s *fastKeySorter) Swap(i, j int) {
	s.keys[i], s.keys[j], s.vals[i], s.vals[j] = s.keys[j], s.keys[i], s.vals[j], s.vals[i]
}
func (s *fastKeySorter) Less(i, j int) bool {
	ki, kj := s.keys[i], s.keys[j]
	lki, lkj := len(ki), len(kj)
	if lki != lkj { //primary key (key length) is significant
		return lki < lkj //use it
	}
	//default use the secondary key: lexicographic
	return ki < kj
}
