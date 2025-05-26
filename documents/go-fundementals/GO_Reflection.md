# Type assertion
[youtube](https://www.youtube.com/watch?v=T2fqLam1iuk&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=33)

"interface{} says nothing" since it has no methods

It's a "generic" thing, but sometimes we need its "real" type

We can extract a specific type with a type assertion (aka "downcasting")
Basically, type assertions let you extract the concrete type from an
interface value.

This has the form `value.(T)` for some type T
```go
var w io.Writer = os.Stdout

f := w.(*os.File)      // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer
```

Type assertions have the two result version and we can avoid panic
```Go
var w io.Writer = os.Stdout

f, ok := w.(*os.File)       // success: ok = true, f == os.Stdout
b, ok := w.(*bytes.Buffer)  // failure: ok = false, b == nil
```

* Someone already put a concrete value (like *os.File) into an interface variable (io.Writer)
* The type assertion lets you get that concrete value back out
* You're essentially saying "I believe this interface actually holds a *os.File, let me get it"

# Reflection
Type Assertions vs Reflection:

Type assertions: Simple extraction of concrete types from interfaces
Reflection: A much broader system that lets programs examine and manipulate their own structure at runtime

A bit of history of CS:
In traditional compiled languages (like C), when you compile your code:
```
Source Code -> Machine code
int x = 5; -> mov eax, 5
```
The machine code has no idea that `eax` used to be an `int` named `x`. All
type information is completely lost- its just raw assembly instructions

Moden Languages with Reflection:
Languages like Go, Java do something different. They embed extra metadata 
alongside the machine code
```
Source Code -> Machine Code + Type Metadata
```
So when you compile a Go program, the binary contains:
1. The actual executable machine instructions
2. Extra data tables describing types, method names, field named, etc

This is why Go binaries are larger than equivalent C programs, they're carrying around this type information.

How Reflection Uses This:
When you use reflection in Go:
```go
var x interface{} = 42
t := reflect.TypeOf(x)  // "int"
```
The `reflect.TypeOf()` function looks up that embedded type metadata to tell
you "this value is an int"

So the type information in modern languages solved the "lost type information'
problem by just... not losing it! They keep it around in the binary for reflection to use

## Deep Equality
We can use the `reflect` package in UTs to check equality
```Go
want := struct{
    a: "a string",
    // Not comparable with ==
    b: []int{1, 2, 3},
}

got := gotGetIt(...)
if !reflect.DeepEqual(got, want) {
    t.Errorf("bad response: got=%#v, want=%#v", got, want)
}
```
You can use `github.com/kylelemon/godebug/pretty` to show a deep diff

## Switching on type
We can also use type assertion in a switch statement (matching a type not a value).
Below is an example of how a part of Println works:
```Go
func Println(args ...interface{}) {
    buf := make([]byte, 0, 80)

    for arg := range args {
        switch a := arg.(type) {
            case string:                            // concrete type
                buf = append(buf, a...)
            case Stringer:                          // interface
                buf = append(buf, a.String()...)
            ...
        }
    }
}
```
Here the switch variable a has a specific type if the case has a single type

## Hard JSON
Not all JSON messages are well-behaved
What if some keys depend on others in the message?

Like this bullshit, (notice album is a key and a value)
```json
{
    "item": "album",
    "album": {"title": "Dark side of the moon"}
}

{
    "item": "song",
    "song": {"title": "bella donna", "artist": "Stevie Nicks"}
}
```

### Custom JSON Decoding
We'll make a wrapper and a custom decoder

```go
type response struct {
    Item    string  `json:"item"`
    Album   string
    Title   string
    Artist  string
}

type respWrapper struct {
    response
}
```
We need respWrapper because it must have a separate unmarshal method from the response type (see below)

```Go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type response struct {
    Item    string  `json:"item"`
    Album   string
    Title   string
    Artist  string
}

type respWrapper struct {
    response
}

var j1 = `{
    "item": "album",
    "album": {"title": "Dark side of the moon"}
}`

var j2 = `{
    "item": "song",
    "song": {"title": "bella donna", "artist": "Stevie Nicks"}
}`

func main() {
	var resp1, resp2 respWrapper
	var err error

	if err = json.Unmarshal([]byte(j1), &resp1); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", resp1.response)

	if err = json.Unmarshal([]byte(j2), &resp2); err != nil {
		log.Fatal()
	}

	fmt.Printf("%#v\n", resp2.response)
	
}

func (r *respWrapper) UnmarshalJSON(b []byte) (err error) {
	// map[string]interface{} represents a JSON
	var raw map[string]interface{}

	// only unmarhsals item because only the Item field in struct has json tag
	err = json.Unmarshal(b, &r.response) 	// ignore errors
	err = json.Unmarshal(b, &raw)

	switch r.Item {
	case "album": 
		inner, ok := raw["album"].(map[string]interface{})
		if ok {
			if album, ok := inner["title"].(string); ok {
				r.Album = album
			}
		}
	case "song":
		inner, ok := raw["song"].(map[string]interface{})
		if ok {
			if title, ok := inner["title"].(string); ok {
				r.Title = title
			}
			
			if artist, ok := inner["artist"].(string); ok {
				r.Artist = artist
			}
		}
	}

	return
}
```