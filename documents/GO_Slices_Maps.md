# Arrays

Arrays are typed by size, which is fixed at compile time

```Go
// All these are equivalent
var a [3]int
var b [3]int{0,0,0}
var c [...]{0,0,0}  // sized by initializer

var d [3]int
d = b               // elements are copied (array has no string descriptor like in strings, so its copied)

var m [...]int{1,2,3,4}

c = m               // TYPE MISMATCH (its sizes different)
```

Arrays are passed in by value, thus elements are copied

# Slices
Slices have variable length, backed by some array; and they behave very much like string and string descriptors.
Slice variable itself is a descriptor (descriptor has pointer to slice, length, and capacity) and assigning them to 
different variables simply overwrites the descriptor itself, not changing the underlying array

```Go
var a []int             // nil, no storage
var b = []int{1,2}      // initialized

a = append(a , 1)       // append to nil OK
b = append(b, 3)        // []int{1,2,3}

a = b                   // overwrites a

e[0] == b[0]            // true
```

Slices are passed by reference; no copying, updating OK

## The off-by-one bug

Slices are index like [8:11]
(read as the starting element and one past the ending element, so
this way we have 11 - 8 = 3 elements in our slice)

For loops work the same way in most cases:
```Go
for i := 8; i < 11; i++ {       // in math written [8, 11)
    ...
}
```

So basically, include 8 is inclusive and 11th element is not. 

## Slices vs Arrays
Most GO APIs take slices as inputs, not arrays

| Slice                         | Array                        |
| ----------------------------- | ---------------------------- |
| Variable length               | Length fixed at compile time |
| Passed by Reference           | Passed by value              |
| Not comparable                | Comparable (==)              |
| Cannot be used as map key     | Can be used as map key       |
| Has copy & append helpers     | ---                          |
| Useful as function parameters | Useful as "pseudo" constants |

Slices can't be used as keys for maps because the keys must be comparable and slices are not comparable
because of variable length sizes


# Maps
Maps are dictionaries: indexed by key, returning a value

You can read from a nil map, but inserting will panic

```Go
var m map[string]int        // nil, no storage
p := make(map[string]int)   // non-nil but empty

a := p["the"]               // returns 0
b := m["the"]               // same thing
m["and"] = 1                // PANIC - nil map
m = p
m["and"]++                  // Ok, same map as p now
c := p["and"]               // returns 1
```

Maps are passed by reference; no copying, updating OK
(so basically, m and p above are some descriptors with pointers in them that point to the underlying hashTable)
The type used for the key must have == or != defined (not slices, maps, or funcs)


Maps can't be compared to one another; maps can be compared to nil as a special case.
```Go
var m = map[string]int{
    "and": 1,
    "the": 1,
    "or":  2,
}

var n map[string]int

b := m == n             // SYNTAX ERROR (can compare maps)
c := n == nil           // true
d := len(m)             // 3
e := cap(m)             // TYPE MISMATCH ERROR (can set capacities for maps, but not really view them)
```

Maps have a special two-result lookup function

The second variable tells you if the key was there 
(can have a case where is the result for key 0 because it doesn't exist, or because it exists and value is 0)

```Go
p := map[string]int{}       // non-nil but empty
a := p["the"]               // returns 0
b, ok := p["and"]           // returns 0 and false, ok used by convention

p["the"]++

c, ok := p["the"]           // returns 1 and true

if w, ok := p["the"]; ok {
    // we know w is not the default value
    ...
}
```

# Built In Functions

Each type has certain built-in functions

```Go
len(s)          //string      string length
len(a), cap(a)  //array       array length, capacity (constant)

make(T, x)      //slice       slice of type T with length x and capacity x
make(T, x, y)   //slice       slice of type T with length x and capacity y

copy(c, d)      //slice       copy d to c; # = min of the two lengths
c=append(c, d)  //slice       append d to c and return a new slice result

len(s), cap(s)  //slice       slice length and capacity

make(T)         //map         map of type T
make(T, x)      //map         map of type T with space enough for x elements

delete(m, k)    //map         delete key k (if present, else no change) for map m

len(m)          //map         map length
```

# Make nil useful
Nil is a type of zero: it indicates the absence of something

Many built-ins are safe: len, cap, range

```Go
var s []int
var m map[string]int

l := len(s)             // length of nil slice is 0
i, ok := m["int"]       // 0, false for any missing key

for _, v := range s {   // skip if s is nil or empty
    ...
}
```

"Make the zero value useful" - Rob Pike
[Understanding Nil](https://www.youtube.com/watch?v=ynoY2xz-F8s)

# Example code
```Go
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	scan := bufio.NewScanner(os.Stdin)
	words := make(map[string]int)

	scan.Split(bufio.ScanWords)

	for scan.Scan() {
		words[scan.Text()]++
	}

	fmt.Println(len(words), "unique words")

    // struct to group map data (needed for sorting)
	type kv struct {
		key string
		val int
	}

    // can't sort map itself, so create slice of structs and sort the structs
	var ss []kv

	for k, v := range words {
		ss = append(ss, kv{k, v})
	}

    // need to define how we want our struct to be sorted, here we decide on largest to smalles based on val of our map
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].val > ss[j].val
	})

    // give only the first three popular words
	for _, s := range ss[:3] {
		fmt.Println(s.key, "appears", s.val, "times")
	}
}
```

# Slices in Detail
To understand nil slices vs empty slices, length vs capacity, go to: [link](https://www.youtube.com/watch?v=pHl9r3B2DFI&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=12)