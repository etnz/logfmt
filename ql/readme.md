# logfmt/ql [![GoDoc](https://godoc.org/github.com/etnz/logfmt/ql?status.svg)](https://godoc.org/github.com/etnz/logfmt/ql)


'ql' is a simple query langage and interpreter to evaluate a simple expression on top of a logfmt record.


## How to write a query ?

'Key' are evaluated to key's value. Key name is prefixed by '.'

      record   user=John mail=john@doe.com
      query    .user
      result   John

Comparison is to compare two values

      record   in=120 out=125
      query    .in < 200
      result   true
      query    .in < .out
      result   true

Comparison available operators are: '<', '=', '>'

Matching a value against a regular expression

      record   user=johndoe@mail.com
      query    .user ~ /john.*/
      result   true

Regular expression literal are delimited by '/' character. To write a '/'  inside the regular expression you need to escape it : '/path\/subpath/'

Logic arithmetic: Comparisons and matchings can be combined using usual boolean arithmetic

      record   user=johndoe@mail.com age=20
      query    .user ~ /john.*/ and .age < 40
      result   true

'AND' operator has priority over 'OR'

`.a OR .b AND .c` is equivalent to `.a  OR  ( .b  AND  .c )`

Space delimiter: logfmt keys can be anything but ' ', therefore key names *must* be delimited by space.

      query    '(.a AND .b)'  is not valid '.b)' is a single key name.
      query    '(.a AND .b )' is valid.

Be careful.  