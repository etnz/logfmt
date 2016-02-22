package logfmt

import (
	"bufio"
	"io"
	"os"
	"sync"
)

// New creates a new Logger.
// Out is the destination to write logs to.
func New(out io.Writer) *Logger { return &Logger{out: bufio.NewWriter(out)} }

//Default is the default Logger implementation
var Default = New(os.Stderr)

//Logger is a basic object that logs records into an output io.Writer
//
type Logger struct {
	lock sync.Mutex // to force one log at a time
	out  *bufio.Writer
}

//Log a Record to the Logger Output
func (l *Logger) Log(rec Record) {
	// acquire the lock to make sure nobody writes at the same time
	l.lock.Lock()
	// write the log

	// not thread safe at all, it reuses a shared buffer (keyvals)
	// but it's ok, when doing Log() since Log is protected against
	keyvals.keys = keyvals.keys[0:0]
	keyvals.vals = keyvals.vals[0:0]
	logTo(l.out, rec, keyvals)
	l.out.WriteRune('\n')
	l.out.Flush()
	l.lock.Unlock() // we don't use defer, it takes a few extra seconds
}

// a local buffer of keyvals
var (
	bufSize int64 = 1024
	keyvals       = &fastKeySorter{
		keys: make([]string, 0, bufSize),
		vals: make([]*string, 0, bufSize),
	}
)
