package ql

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"regexp"
)

// eval a valid AST against a Record

// AsType return a string description of xval runtime type.
//
// values are:
//
//    'nil'
//    'val'
//    'bool'
//    'number'
//    'decimal'
//    'duration'
//    'regexp'
func AsType(xval interface{}, xerr error) (val string) {
	if xerr != nil {
		return "'err'"
	}

	switch x := xval.(type) {

	case *string:
		if x == nil {
			return "'nil'"
		}
		return "'val'"

	case bool:
		return "'bool'"

	case int64:
		return "'number'"

	case float64:
		return "'decimal'"

	case time.Duration:
		return "'duration'"

	case *regexp.Regexp:
		return "'regexp'"

	default:
		return fmt.Sprintf("'invalid %T'", x)
	}
	return
}

// AsLiteral convert xval runtime type to a literal string
//
// if xerr is not nil the special value '<err>' is used
//
// if xval is nil the special value '<nil>' is used
//
// Otherwise:
//
//     *string: is quoted
//     bool: is represented as 'true' or 'false'
//     int64: as base10 integer
//     float64: using the %g format
//     time.Duration: String() method is called
//     *regexp.Regexp: String() method is called between '/'
func AsLiteral(xval interface{}, xerr error) (val string) {
	if xerr != nil {
		return fmt.Sprintf("<err:%v>", xerr)
	}
	if xval == nil {
		return "<nil>"
	}

	switch x := xval.(type) {

	case *string:
		return strconv.Quote(*x)

	case bool:
		return strconv.FormatBool(x)

	case int64:
		return strconv.FormatInt(x, 10)

	case float64:
		return strconv.FormatFloat(x, 'g', -1, 64)

	case time.Duration:
		return x.String()

	case *regexp.Regexp:
		return "/" + x.String() + "/"

	default:
		return fmt.Sprintf("<invalid %T>", x)
	}
	return
}

// AsBool check the runtime type and convert it.
//
// if xerr is not nil, it is returned in 'err'
func AsBool(xval interface{}, xerr error) (val bool, err error) {
	//eval it as it come
	if xerr != nil {
		err = xerr
		return
	}

	if xval == nil {
		val = false
		return
	}
	switch x := xval.(type) {

	case *string:
		val, err = strconv.ParseBool(*x)

	case bool:
		val = x

	case int64:
		val = x != 0

	case float64:
		val = x != 0

	default:
		err = fmt.Errorf("got %s expecting 'bool'", AsType(xval, xerr))
	}
	return
}

// AsRegexp check the runtime type and convert it
//
// if xerr is not nil, it is returned in 'err'
func AsRegexp(xval interface{}, xerr error) (val *regexp.Regexp, err error) {
	//eval it as it come
	if xerr != nil {
		err = xerr
		return
	}
	//check the xval type, it has to be a bool
	var isregexp bool
	val, isregexp = xval.(*regexp.Regexp)
	if !isregexp {
		err = fmt.Errorf("got %s expecting 'regexp'", AsType(xval, xerr))
		return
	}
	return
}

// AsDecimal try to convert runtime value 'xval' to float64
//
//    if xerr is not nil it is returned
//    *string: nil is converted to NaN, otherwise an attempt is made to parse it as a ql literal.
//    int64, float64: use natural type conversion
//    time.Duration: is converted in float64 number of seconds
func AsDecimal(xval interface{}, xerr error) (val float64, err error) {
	if xerr != nil {
		err = xerr
		return
	}
	switch v := xval.(type) {

	case *string:
		if v == nil {
			val = math.NaN()
			return
		}
		//scan the string
		parser := &parser{src: newScanner(strings.NewReader(*v))}
		parser.Next()
		var xval Expr // the parsing of the value itself
		xval, err = parser.LiteralExpr(), parser.err
		if err != nil {
			err = fmt.Errorf("cannot evaluate record value %q as decimal: %v", *v, err)
			return
		}
		return AsDecimal(Eval(xval, nil))

	case int64:
		val = float64(v)

	case float64:
		val = v

	case time.Duration:
		val = v.Seconds()

	default:
		err = fmt.Errorf("cannot convert %s to decimal", AsType(xval, xerr))
	}
	return
}

// AsValue tries to convert the runtime value to *string runtime type
//
//    *string: is left unchanged
//    bool, int64, float64, duration: literal value is used
//    *regexp.Regexp, the regexp definition, undercorated is used
func AsValue(xval interface{}, xerr error) (val *string, err error) {
	if xerr != nil {
		err = xerr
		return
	}
	switch v := xval.(type) {

	case *string:
		val = v
		return

	case bool:
		s := strconv.FormatBool(v)
		val = &s

	case int64:
		s := strconv.FormatInt(v, 10)
		val = &s

	case float64:
		s := strconv.FormatFloat(v, 'g', -1, 64)
		val = &s

	case time.Duration:
		s := v.String()
		val = &s

	case *regexp.Regexp:
		s := v.String()
		val = &s

	default:
		err = fmt.Errorf("cannot convert %s to value", AsType(xval, xerr))
	}
	return
}

// Eval the expr using a 'rec' map of 'value'
func Eval(expr Expr, rec map[string]*string) (val interface{}, err error) {

	switch x := expr.(type) {

	case *BinaryExpr:
		switch x.Op {

		case AND:
			//recursively eval X and convert it as boolean
			bval, berr := AsBool(Eval(x.X, rec))
			if !bval || berr != nil {
				val, err = bval, berr
				return
			}
			// no shortcut possible, I need to eval the next one
			// val from X is already true
			val, err = AsBool(Eval(x.Y, rec))

		case OR:
			//recursively eval X and convert it as boolean
			bval, berr := AsBool(Eval(x.X, rec))
			if bval || berr != nil { // possible shortcut already true  or there is an error
				val, err = bval, berr
				return
			}

			// no shotcut
			// val from X is already true
			val, err = AsBool(Eval(x.Y, rec))

		default:
			err = fmt.Errorf("unsupported binary operator %q", x.Op)
		}

	case *UnaryExpr:
		switch x.Op {

		case NOT:
			bval, berr := AsBool(Eval(x.X, rec))
			if berr != nil {
				err = berr
				return
			}
			val = !bval

		default:
			err = fmt.Errorf("unsupported unary operator %q", x.Op)
		}

	case *PostCompExpr:

		switch x.Op {

		case EXISTS: // handles 'exist' the only one !
			//exist can only be evaluated on IDENT
			if x.X.Kind == IDENT {
				// okay, test if it exists
				_, exists := rec[x.X.Value[1:]]
				// fmt.Printf("exists(%v) = %v\n", Fmt(x.X), exists)
				val = exists
				return // ok
			}
			// invalid type
			err = fmt.Errorf("cannot test existence on %q. Only %q is supported", x.X.Kind, Token(IDENT))

		default: //any other op
			err = fmt.Errorf("unsupported postfix operator %q", x.Op)
		}

	case *ParenExpr:
		return Eval(x.X, rec)

	case *CompExpr:
		switch x.Op {

		case MATCH: // regexp matching
			var lhs *string
			// convert the lhs as much as possible as a value (*string)
			lhs, err = AsValue(Eval(x.X, rec))
			if err != nil {
				err = fmt.Errorf("invalid left hand side of '~' comparison: cannot convert to 'value': %v", err)
				return
			}
			// then convert the rhs as a regexp
			var rhs *regexp.Regexp
			rhs, err = AsRegexp(Eval(x.Y, rec))
			if err != nil {
				err = fmt.Errorf("invalid right hand side of '~' comparison: cannot convert to 'regexp': %v", err)
				return
			}
			// and ... match them both
			if lhs == nil {
				val = rhs.Match(nil) // a reg can match nil, right? or shall I make it "false" don't know, yet
			} else {
				val = rhs.Match([]byte(*lhs))
			}

		case LT, GT, EQ:
			var lhs, rhs float64
			lhs, err = AsDecimal(Eval(x.X, rec))
			if err != nil {
				return
			}
			rhs, err = AsDecimal(Eval(x.Y, rec))
			if err != nil {
				return
			}
			//ok both hs have been evaluated
			switch x.Op {

			case LT:
				val = lhs < rhs

			case GT:
				val = lhs > rhs

			case EQ:
				val = lhs == rhs
			}
		}

	case *FuncExpr:
		err = fmt.Errorf("unsupported function %q", x.Func)

	case *Literal:
		//all literal have a different runtime type
		switch x.Kind {

		case IDENT: // *string
			ident := x.Value[1:]
			val = rec[ident]
			return

		case NUMBER: // int64
			return strconv.ParseInt(x.Value, 10, 64)

		case DECIMAL: // float64
			return strconv.ParseFloat(x.Value, 64)

		case DURATION: // time.Duration
			val, err = time.ParseDuration(x.Value)
			if err != nil {
				err = fmt.Errorf("Invalid duration %q.", x.Value)
				return
			}

		case REGEXP: // *regexp.Regexp
			pattern := x.Value
			pattern = pattern[1 : len(pattern)-1] //strip away the first and last '/'
			// unescape the escaped '/'
			pattern = strings.Replace(pattern, "\\/", "/", -1)
			val, err = regexp.Compile(pattern)
		}

	default:
		err = fmt.Errorf("Unknown syntax tree node %T", x)

	}
	return
}
