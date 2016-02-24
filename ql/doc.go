// Package ql implements a query langage for logfmt records.
//
// The langage offers:
//
//    - a simple bool arithmetic ('or', 'and' and 'not') to combine test all together
//    - comparison operators ( '<', '>', '=', '~')
//    - existential operator ('?') in postfix
//    - literal for record attributes: regexp, numbers, decimals, durations
//    - a set of builtin functions
//
// A ql statement can be evaluated on any given Record, it might return one of the following runtime type:
//
//    - *string: for the raw attribute value
//    - bool: as the result of any comparison
//    - *regexp.Regexp: for regexp Literal
//    - int64: for numbers
//    - float64: for decimals
//    - time.Duration: for durations
//
//
// For instance it is possible to match 'username' key against a regular expression.
//
//       .username ~ /eric\..*/   : keep all "eric" from the logs
//       .load  > 153ms         : keep big loads
//       .count  < 5            : keep pages with small visit number.
//       ...
//
package ql

const (
	// text to describe ql syntax
	WriteAQuery = `
  How to write a query ?

    'Key' are evaluated to the key's value. Key name is prefixed by '.'
      	  record   user=John mail=john@doe.com
          query    .user
          result   John
    
    Comparison is to compare two values
      	  record   in=120 out=125
          query    .in < .out
          result   true
      Comparison available operators are: '<', '=', '>'
    
    Matching: to match a value against a regular expression
      	  record   user=johndoe@mail.com
          query    .user ~ /john.*/
          result   true
      Regular expression literal are delimited by '/' character. To write a '/' 
      inside the regular expression you need to escape it : '/path\/subpath/'
    
    Logic arithmetic: Comparisons and matchings can be combined using usual boolean arithmetic
      	  record   user=johndoe@mail.com age=20
          query    .user ~ /john.*/ and .age < 40
          result   true
      'AND' operator has priority over 'OR'
          '.a OR .b AND .c' is equivalent to '.a  OR  ( .b  AND  .c )'
    
    Space delimiter: logfmt keys can be anything but ' ', therefore key names *must* be 
      delimited by space.
          query    '(.a AND .b)'  is not valid '.b)' is a single key name.
          query    '(.a AND .b )' is valid.
      Be careful.  
`
)
