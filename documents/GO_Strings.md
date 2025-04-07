# Strings

Types related to strings:
* byte: a synonym for uint8
* rune: a synonym for int32 for characters
* string: an immutable sequence of "characters"
  * physically a sequence of bytes (UTF-8 encoding)
  * logically a sequence of (unicode runes)

Runes (characters) are enclosed in a single quotes: 'a'
"Raw" strings use backtick quotes: `string with "quotes"`
They also don't evaluate escape characters such as \n

Go uses "runes" (which are actually just int32 values) to represent Unicode code points internally. This allows Go to handle any character in the Unicode standard.
UTF-8 encoding is indeed more space-efficient than using 32 bits for every character. In UTF-8:

* ASCII characters (English letters, numbers, basic punctuation) use just 1 byte
* Most European and Middle Eastern characters use 2 bytes
* Most Asian characters (like Chinese, Japanese) use 3 bytes
* Rare characters might use 4 bytes

Go's string type is actually a sequence of bytes, not runes. When you need to process individual characters, you convert strings to runes. This design lets Go handle text efficiently while still supporting the full Unicode character set.

## String Structure
The internal string representation is a pointer and a length
Strings are immutable and can share the underlying storage

Immutability in Go strings means that once a string is created, you cannot change the characters within it. If you try to modify a string, Go actually creates a new string with the changes. This is different from languages like C where strings are mutable character arrays.

A string descriptor in Go is a small data structure that contains:
1. A pointer to the underlying byte array
2. The length of the string

It's different from a simple pointer because it includes metadata(the length). In Go's implementation, a string is effectively a struct with these two fields.

```Go
s := "hello world"
hello := s[:5]
world := s[7:]
```
What's happening here is:
* `s` is a string descriptor pointing to the full "hello, world" bytes in memory
* `hello` is a new string descriptor that points to the same memory location as `s`, but with a length of 5
* `world` is another string descriptor pointing to the 7th position in the same memory, with its own length

The key insigh is that no new copies are made during slicing operation. All three variables reference the same underlying byte array in memory and its immutable, changing one string would affect other strings that share the same storage.

Below is another example to further reinforce understanding
```Go
s := "the quick brown fox"

a := len(s)                     // 19
b := s[:3]                      // "the"
c := s[4:9]                     // "quick"
d := s[:4] + "slow" * s[9:]     // replaces "quick"

s[5] = 'a'                      // SYNTAX ERROR
s += "es"                       // now plural (copied)
```
Strings are passed by reference, thus they aren't copied.

s is a string descriptor that points to the string "the quick brown fox"
b and c are also string descriptors that points to same underlying string "the quick brown fox"
d however, is a string descriptor that points to entirely new location wherever the new string "the slow brown fox" is created

s[5] would fail because strings are immutable and can't be changed
s += "es" doesn't fail, but s string descriptor no longer points to original string array, it points to new array that copies original contents s and adds "es" at the end. 

## String Functions
Package `strings` has many functions on strings

```Go
s := "a string"

x := len(s)                 // built in, = 8
strings.Contains(s, "g")    // returns true
strings.Contains(s, "x")    // returns false

strings.HasPrefix(s, "a")   // returns true
strings.Index(s, "string")  // returns 2

s = strings.ToUpper(s)      // returns "A STRING", but makes new copy of string, since immutable
```

len(s) - O(1)

This is a constant time operation because the string descriptor stores the length, so it's just returning a pre-computed value.

strings.Contains(s, "g") - O(n+m)

Where n is the length of the string being searched (s), and m is the length of the substring ("g")
This uses the Boyer-Moore algorithm or a variant, which has to scan through the string in the worst case.

strings.Contains(s, "x") - O(n+m)

Same complexity as above, but in practice might terminate earlier since it returns as soon as a match is found.

strings.HasPrefix(s, "a") - O(m)

Where m is the length of the prefix being checked ("a")
This only needs to check the first m characters of the string.

strings.Index(s, "string") - O(n+m)

Where n is the length of s and m is the length of "string"
