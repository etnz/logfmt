package ql

import (
	"fmt"
	"io"
)

// parse an expression into it's AST counter part

type parser struct {
	src *scanner
	err error
}

// Parse convert any 'src' into an Expr (or error if it is not possible)
func Parse(src io.Reader) (x Expr, err error) {
	//Start by building the scanner and cosuming the first token
	parser := &parser{src: newScanner(src)}
	parser.Next()
	return parser.OrExpr(), parser.err
}
func (p *parser) Next() { p.src.Next() }

func (p *parser) OrExpr() Expr {

	x := p.AndExpr()

	if p.src.ttype == OR {
		pos := p.src.start
		p.Next()
		y := p.OrExpr()
		return &BinaryExpr{
			X:     x,
			Op:    OR,
			OpPos: pos,
			Y:     y,
		}
	}
	return x

}

func (p *parser) AndExpr() Expr {
	x := p.UnaryExpr()

	if p.src.ttype == AND {
		pos := p.src.start
		p.Next()
		return &BinaryExpr{
			X:     x,
			Op:    AND,
			OpPos: pos,
			Y:     p.UnaryExpr(),
		}

	}

	return x
}

func (p *parser) UnaryExpr() Expr {

	// unary are paren  or any comparator
	switch p.src.ttype {
	case LPAREN:

		lpos := p.src.start
		p.Next()

		x := p.OrExpr()
		if p.src.ttype != RPAREN {
			p.err = fmt.Errorf("%v parenthesis mismatch, expected ')' found %v instead", p.src.start, p.src.ttype)
			return nil
		}
		rpos := p.src.start
		p.Next()

		return &ParenExpr{
			LParenPos: lpos,
			X:         x,
			RParenPos: rpos,
		}

	case NOT:
		pos := p.src.start
		p.Next()

		return &UnaryExpr{
			Op:    NOT,
			OpPos: pos,
			X:     p.UnaryExpr(),
		}

	default:
		return p.LiteralOpExpr()

	}
}

func (p *parser) LiteralOpExpr() Expr {

	// lhs can be either a literal OR a function
	var lhs Expr

	switch p.src.ttype {

	case FUNCTION:
		pos := p.src.start
		fname := p.src.token.String()
		p.Next()

		// it has to be a '('
		if p.src.ttype != LPAREN {
			p.err = fmt.Errorf("%v Invalid function call, need to start with a '('. Found %v instead", p.src.start, p.src.ttype)
			return nil
		}
		p.Next()

		//the  X expr
		x := p.OrExpr()

		// it has to be a ')'
		if p.src.ttype != RPAREN {
			p.err = fmt.Errorf("%v Invalid function call, need to end with a ')'. Found %q:%v instead", p.src.start, p.src.token.String(), p.src.ttype)
			return nil
		}
		rpos := p.src.start
		p.Next()
		// Ok great
		lhs = &FuncExpr{
			Func:      fname,
			FuncPos:   pos,
			RParenPos: rpos,
			X:         x,
		}
	default:
		lhs = p.LiteralExpr()

	}

	op := p.src.ttype
	switch op {
	case EXISTS:
		p.Next() //consume it
		return &PostCompExpr{
			X:     lhs.(*Literal),
			Op:    op,
			OpPos: p.src.start,
		}
	case LT, GT, EQ, MATCH:
		p.Next() //consume it
		pos := p.src.start
		y := p.LiteralExpr()
		return &CompExpr{
			X:     lhs,
			Op:    op,
			OpPos: pos,
			Y:     y,
		}
	default:
		return lhs
	}

}

func (p *parser) LiteralExpr() *Literal {

	switch p.src.ttype {
	case IDENT, REGEXP, DURATION, NUMBER:
		defer p.Next()
		return &Literal{
			Kind:   p.src.ttype,
			LitPos: p.src.start,
			Value:  p.src.token.String(),
		}

	default:
		p.err = fmt.Errorf("%v Syntax Error: expecting one literal: Identifier, Regexp, Duration or Number; got %v instead", p.src.start, p.src.ttype)
		return nil
	}
}
