# What are regular expressions?
Its just regex for Go, have to be careful when using it.
Go's reglar expression syntax is a subset of what some other languages have.
This is to avoid the performance impact of catastrophic backtracking or other malicious intent.

## Searching in Strings
Use the `strings` package for simple searches

Carefully use the `regexp` package for complex searches and validation

### Simple String Searches
Boolean searches:
* `strings.HasPrefix(s, substr)`
* `strings.HasSuffix(s, substr)`
* `strings.Contains(s, substr)`

Location searches:
* `strings.LastIndex(s, substr)`
* `strings.LastIndexByte(s, char)`

Search and replace:
* `strings.Replace(s, substr, replacement, count)`
* `strings.ReplaceAll(s, substr, replacement)`

## Regex Expressions
Syntax for repetition and character classes:
* `.` is any character
* `.*` is zero or more
* `.+` is one or more
* `.?` is zero or one (prefer one)
* `a{n}` is n repetitions of the letter "a"
* `a{n,m}` is n to m repetitions of the letter "a"
* `[a-z]` is a character class (here letters a-z)
* `[^a-z]` is an negated class (here anything except a-z)

Syntax for location:
* `xy` is "x" followed by "y" `(a[sub]string!)`
* `x|y` is either "x" or "y"
* `^x` is "x" at the beginning
* `x$` is "x" at the end
* `^x$` is "x" by itself (its the whole thing)
* `\b` is a word boundary
* `\bx\b` is the "x" by itself (inside the string)
* `(x)` is a capture group

Some built-in character classes:
* `\d` is a decimal digit
* `\w` is a word character ([0-9A-Za-z_])
* `\s` is whitespace
Below character classes are for unicodes
* `[[:alpha]]` is any alphabetic character
* `[[:alnum:]]` is any alphanumeric character
* `[[punct:]]` is any punctuation character
* `[[:print:]]` is any printable character
* `[[:xdigit:]]` is any hexadecimal character

See [https://golang.org/pkg/regexp/syntax/](https://golang.org/pkg/regexp/syntax/)

UUID validation
Be careful when validating UUID with regular expressions