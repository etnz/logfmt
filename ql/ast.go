package ql

// Expr All AST node implements this one
type Expr interface {
	Pos() int
	End() int
}

type (

	// Literal holds any literal in 'ql':
	//
	//     IDENT   : '.' followed by any char above ' '
	//     REGEXP  : a regular expression between '/' e.g. /com.*/
	//     NUMBER  : a digit or '+-' sign, followed by digits
	//     DECIMAL : a digit or '+-' sign, followed by digits and containing a '.'
	//     DURATION: a NUMBER followed by any of 'smhnµu' (the character used to for time unit 'ms' or 'µs' 'h' etc.)
	Literal struct {
		Kind   Token // IDENT, REGEXP, NUMBER, DURATION, or TIME
		Value  string
		LitPos int
	}

	// BinaryExpr is for boolean arithmetic binary operations, namely 'OR' or 'AND'
	BinaryExpr struct {
		X     Expr
		OpPos int
		Op    Token // AND, or OR
		Y     Expr
	}

	//UnaryExpr is for the boolean arithmetic 'NOT'
	UnaryExpr struct {
		OpPos int
		Op    Token // NOT is the only one
		X     Expr
	}

	// ParenExpr encapsulate any Expr, to override priority rules.
	ParenExpr struct {
		LParenPos int
		X         Expr
		RParenPos int
	}

	// FuncExpr is a function call stratement, nota bene, function call only support a single parameter
	FuncExpr struct {
		FuncPos   int    // function's name start position
		Func      string // function's name
		X         Expr
		RParenPos int
	}

	// CompExpr is to compare a logfmt keyvalue with something ( a regexp or a date or a number)
	CompExpr struct {
		X     Expr
		OpPos int
		Op    Token //  ~ < > =
		Y     Expr
	}

	// PostCompExpr is a postfix operator for keys, right now only '?' (exists) is supported
	PostCompExpr struct {
		X     *Literal
		Op    Token //  ?
		OpPos int
	}
)

// Pos returns the literal first character position
func (l Literal) Pos() int { return l.LitPos }

// Pos returns the binary first character position
func (l BinaryExpr) Pos() int { return l.X.Pos() }

// Pos returns the UnaryExpr first character position
func (l UnaryExpr) Pos() int { return l.OpPos }

// Pos returns the ParenExpr first character position
func (l ParenExpr) Pos() int { return l.LParenPos }

// Pos returns the FuncExpr first character position
func (l FuncExpr) Pos() int { return l.FuncPos }

// Pos returns the CompExpr first character position
func (l CompExpr) Pos() int { return l.X.Pos() }

// Pos returns the PostCompExpr first character position
func (l PostCompExpr) Pos() int { return l.X.Pos() }

// End returns the Literal last character position
func (l Literal) End() int { return l.LitPos + len(l.Value) }

// End returns the BinaryExpr last character position
func (l BinaryExpr) End() int { return l.Y.End() }

// End returns the UnaryExpr last character position
func (l UnaryExpr) End() int { return l.X.End() }

// End returns the ParenExpr last character position
func (l ParenExpr) End() int { return l.RParenPos + 1 }

// End returns the FuncExpr last character position
func (l FuncExpr) End() int { return l.RParenPos + 1 }

// End returns the CompExpr last character position
func (l CompExpr) End() int { return l.Y.End() }

// End returns the PostCompExpr last character position
func (l PostCompExpr) End() int { return l.OpPos + 1 }

// Fmt returns a literal representation of any Expr
//
// Suitable to parse an expression, and to format it in a fixed style.
func Fmt(exp Expr) string {
	switch x := exp.(type) {

	case *Literal:
		return x.Value

	case *BinaryExpr:
		return Fmt(x.X) + "   " + x.Op.String() + "   " + Fmt(x.Y)

	case *CompExpr:
		return Fmt(x.X) + " " + x.Op.String() + " " + Fmt(x.Y)

	case *UnaryExpr:
		return x.Op.String() + "  " + Fmt(x.X)

	case *PostCompExpr:
		return Fmt(x.X) + " " + x.Op.String()

	case *ParenExpr:
		return "(" + Fmt(x.X) + ")"

	case *FuncExpr:
		return x.Func + "( " + Fmt(x.X) + " )"

	default:
		return "<nil>"
	}
}
