package ql

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type scanner struct {
	buf *bufio.Reader

	token bytes.Buffer
	pos   int
	start int
	ttype Token
	err   error
}

func newScanner(src io.Reader) *scanner {
	return &scanner{
		buf: bufio.NewReader(src),
	}
}
func newScannerS(src string) *scanner {
	return &scanner{
		buf: bufio.NewReader(strings.NewReader(src)),
	}
}

func isWhitespace(r rune) bool      { return r <= ' ' && r != eof }
func isIdentifier(r rune) bool      { return r != eof && r > ' ' && r != '"' && r != '=' }
func isRegexp(r rune) bool          { return r != eof && r > ' ' && r != '/' }
func isFunctionTrigger(r rune) bool { return unicode.IsLetter(r) || r == '_' }
func isFunction(r rune) bool        { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' }
func isNumberTrigger(r rune) bool   { return strings.ContainsRune("0123456789+-", r) }
func isNumber(r rune) bool          { return strings.ContainsRune("0123456789.smhnµu", r) }
func isOp(r rune) bool              { return strings.ContainsRune("andortANDORT", r) } //set of and or not

//read one rune ahead
func (s *scanner) read() (r rune) {
	r, _, s.err = s.buf.ReadRune()
	if s.err == nil {
		s.pos++
	}
	return
}
func (s *scanner) unread() {
	if s.err != nil {
		return
	}
	s.err = s.buf.UnreadRune()
	s.pos--
}
func (s *scanner) peek() rune { defer s.unread(); return s.read() }
func (s *scanner) begin() {
	s.token.Reset()
	s.start = s.pos
}

func (s *scanner) Next() {
	s.readAll(isWhitespace) // readall white space just in case

	s.begin() //remember this position as the token's begin
	r := s.read()
	s.token.WriteRune(r)
	switch {
	case r == eof:
		s.ttype = EOF
	case r == '.': //identifier marker
		s.readAll(isIdentifier)
		s.ttype = IDENT
	case r == '/': //identifier marker
		s.readAllEscape(isRegexp)
		s.ttype = REGEXP
		//there MUST be a closing '/'
		r = s.read()
		s.token.WriteRune(r)
		if r != '/' {
			s.err = fmt.Errorf("Invalid regexp token: must end with '/' got %q", r)
		}

	case isFunctionTrigger(r):
		// try to read it fully as an op
		s.readAll(isOp)
		content := s.token.String()
		content = strings.ToUpper(content)
		switch content {
		case "OR":
			s.ttype = OR
		case "AND":
			s.ttype = AND
		case "NOT":
			s.ttype = NOT
		default:
			// this is not an Op
			s.readAll(isFunction)
			s.ttype = FUNCTION
		}

	case isNumberTrigger(r): // number marker
		s.readAll(isNumber)
		content := s.token.String()
		s.ttype = NUMBER //Default

		//contains char reserved for decimals
		if strings.ContainsAny(content, ".") {
			s.ttype = DECIMAL
		}

		//contains char reserved for duration
		if strings.ContainsAny(content, "smhnµu") {
			s.ttype = DURATION
		}

	case r == '~':
		s.ttype = MATCH //~

	case r == '=':
		s.ttype = EQ

	case r == '<':
		s.ttype = LT

	case r == '>':
		s.ttype = GT

	case r == '?':
		s.ttype = EXISTS

	case r == '(':
		s.ttype = LPAREN

	case r == ')':
		s.ttype = RPAREN

	default:
		s.err = fmt.Errorf("Unknown symbol '%s'", r)
	}
}

// read all runes that match the criteria
func (s *scanner) readAll(criteria func(r rune) bool) {
	r := s.read()
	for criteria(r) {
		s.token.WriteRune(r)
		r = s.read()
	}
	s.unread() //unread the last one, it is not a space
}

// read all runes that match the criteria
func (s *scanner) readAllEscape(criteria func(r rune) bool) {
	r := s.read()
	for criteria(r) {
		if r == '\\' {
			s.token.WriteRune(r) //keep it, but skip criteria for the next rune
			r = s.read()
			if r == eof {
				return
			}
		}
		s.token.WriteRune(r)
		r = s.read()
	}
	s.unread() //unread the last one, it is not a space
}
