package ql

const (

	// ILLEGAL token, the zero value for any Token
	ILLEGAL Token = iota

	// EOF is the end of file token
	EOF

	// IDENT ; anything above ' ' and not '"' or '='
	IDENT

	// FUNCTION ; anything usual alpha numeric
	FUNCTION

	// REGEXP anything between a pair of '/'. to write a '/' use '\/'
	REGEXP

	// NUMBER anything starting with a digit (or '+' or '-') but without '.'
	NUMBER

	// DECIMAL is like number but it contains a "." inside.
	DECIMAL

	// DURATION anything starting like a number, but ending with a time unit ( or of [smhdy] )
	DURATION

	// OR the 'or' or 'OR' token.
	OR

	// AND the 'and' or 'AND' token.
	AND

	// NOT the 'not' or 'NOT' token.
	NOT

	// MATCH the match token '~' (as in awk)
	MATCH

	// LT the "lower than" operator '<'
	LT

	// GT the "greater than" operator '>'
	GT

	// EQ the "equal" operator '='
	EQ

	// EXISTS the existance operator "?"
	EXISTS

	// LPAREN usual '('
	LPAREN

	// RPAREN usual ')'
	RPAREN
)

const (
	eof = '\x00'
)

// Token is a special int, to represent one of the const above
type Token int

// String gives a string representation of this token.
//
// Token with a constant value (like OR) are represent as such, others are represented between '<' '>'
func (t Token) String() string {
	switch t {

	case EOF:
		return "<EOF>"

	case IDENT:
		return "<IDENT>"

	case FUNCTION:
		return "<FUNCTION>"

	case REGEXP:
		return "<REGEXP>"

	case NUMBER:
		return "<NUMBER>"

	case DECIMAL:
		return "<DECIMAL>"

	case DURATION:
		return "<DURATION>"

	case OR:
		return "OR"

	case AND:
		return "AND"

	case NOT:
		return "NOT"

	case MATCH:
		return "~"

	case LT:
		return "<"

	case GT:
		return ">"

	case EQ:
		return "="

	case EXISTS:
		return "?"

	case LPAREN:
		return "("

	case RPAREN:
		return ")"
	default:
		return "<ILLEGAL>"
	}
}
