package logreader

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"github.com/etnz/logfmt"
)

const (
	//end of file
	eof = '\x00'
	//end of line
	eol = '\n'
)

// ErrUnterminatedQuote is an error returned when the end of file is reached before the end of quoted string.
var ErrUnterminatedQuote = errors.New("Error unterminated quoted string")

type scanner struct {
	*bufio.Reader
	err error
	eof bool
}

//Reader reads from any source successives records
type Reader interface {
	HasNext() bool
	Next() (rec logfmt.Record, err error)
}

// New instanciate a new Reader
func New(r io.Reader) Reader { return newScanner(r) }

// Parse a single record as string.
func Parse(src string) (rec logfmt.Record, err error) { return newScannerS(src).Next() }

func newScannerS(str string) *scanner { return newScanner(strings.NewReader(str)) }

func newScanner(r io.Reader) *scanner {
	return &scanner{
		Reader: bufio.NewReader(r),
	}
}

//the rules for scanning
func isIdentifier(r rune) bool { return r != eof && r > ' ' && r != '"' && r != '=' }
func isGarbage(r rune) bool    { return r != eof && r != eol && r <= ' ' }
func isString(r rune) bool     { return r != eof && r != '"' && r != '\\' }

// HasNext return true has long as the scanned has not found the end of file
func (s *scanner) HasNext() bool { return !s.eof && s.err == nil }

// Read read a single rune from the src
func (s *scanner) Read() (r rune) {
	r, _, s.err = s.ReadRune()
	if r == eof {
		s.eof = true
		s.err = nil
	}
	return
}

// Unread the previous rune from the src
func (s *scanner) Unread() {
	s.UnreadRune()
}

// Next read runes until it has found a full Record, returns it.
//
// If the source has errors it returns it
func (s *scanner) Next() (record logfmt.Record, err error) {
	rec := logfmt.Rec()
	for {

		if r := s.Read(); r == eol || r == eof {
			return *rec, s.err
		}
		s.Unread()

		key, val := s.Pair()
		(*rec)[key] = val
	}
}

// Pair scan for a couple key(=value)?
func (s *scanner) Pair() (key string, value *string) {

	s.Garbage() //consumer the possible garbage
	key = s.Identifier()

	s.Garbage() //consumer the possible garbage

	if r := s.Read(); r == '=' {
		// separator there might be a value
		s.Garbage() //consumer extra space after =
		//then attempt to read the value (might not exist)
		v := s.Value()
		value = &v
		s.Garbage() //consumer extra space
		return
	}
	s.Unread()
	return
}

//Value scan for a valid value
func (s *scanner) Value() (value string) {

	if r := s.Read(); r == '"' { //it will be a string
		value = s.Str()
		end := s.Read() //read the " or the eof one more time
		if end != '"' {
			s.err = ErrUnterminatedQuote
		}
		return
	}
	s.Unread()

	return s.Identifier()

}

// Garbage consumes as much as possible "garbage" char (separators)
func (s *scanner) Garbage() {
	for r := s.Read(); isGarbage(r); r = s.Read() {
	}
	s.Unread()
}

// Identifier scan for an identifier
func (s *scanner) Identifier() (identifier string) {

	var buf bytes.Buffer
	for r := s.Read(); isIdentifier(r); r = s.Read() {
		buf.WriteRune(r)
	}
	s.Unread()
	return buf.String()
}

// Str scan for a quoted string, and escape any character right after '\'
func (s *scanner) Str() (str string) {

	var buf bytes.Buffer

	for r := s.Read(); isString(r) || r == '\\'; r = s.Read() {
		if r == '\\' {
			r = s.Read()
			// write the next after \ This is more general than just the usual definition, as a side effect, '\' followed by a line break insert the line break
			// to implement a strict version we "could" test r right now, it must be '"' or '\\'
		}
		buf.WriteRune(r)
	}
	s.Unread()

	return buf.String()
}
