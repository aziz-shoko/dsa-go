# Packages
In Go, everything lives inside a package. Every standalone program has a main package.
```Go
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
}
```

Nothing is "global"; it's either in your package or in another
It's either at package scope (as in stuff declared outside of a function) 
or function scope (declared in a function), nothing global

## Package-level declarations
```Go
package secrets
const DefaultUUID = "000000000-0000-0000-0000-000000000000"
var secretKey string

type k8secret struct {
    ....
}

func Do(it string) error {
    ...
}
```

But you can't use the short declration operator := at package scope (makes it hard to parse stuff)

## Packages control visibility
Every name that's capitlized is exported
```Go
package secrets

import ...

// private, not visible to whatever that exports it (since its lowercase)
type internal struct {
    ...
}

// public, visible to whatever that exports to it (since its uppercase)
func GetAll(space, name string) (map[string]string, error) {
    ...
}
```
This means another package in the program can import it
(within a package, everything is visible even across files)

## Imports
Each source file in your package must import what it needs
```Go
package secrets

import (
    "encoding/base64"
    "encoding/base64"
    "fmt"
    "os"
    "strings"
)
```
It may only import what it needs; unused imports are an error
Generally, files of the same package live together in a directory

## No cycles
Go doesn't want cyclic dependencies, for exmaple:
A package "A" cannot import a package that imports A
```Go
package A

import "B"

//--------
package B

import "A"  // WRONG
```
Move common dependencies to a third package, or eliminate them

Another problem that cyclic dependencies create is Initialization
Items within a package get initialized before main
```Go
const A = 1

var B int = C
var C int = A

func Do() error {
    ...
}

func init() {
    ...
}
```
Only the runtime can call init, and init is something that is sort of hidden away and runs before main

## What makes a good package?
A Philosophy of Software Design - by John Ousterhout
A package is essentially means of hiding information.
A package should embed deep functionality behind a simple API

```Go
package os

func Create(name string) (*File, error)
func Open(name string) (*File, error)

func (f *File) Read(b []byte) (n int, err error)
func (f *File) Write(b []byte) (n int, err error)
func (f *File) Close() error
```
The Unix File API is perhaps the best example of this model, the stuff that happen in the background to create, open, and read
files is complicated and all the complicated parts hides behind a set of simple APIs like Read, Write, Create, etc

Roughly five functions hide a lot of complexity from the user

# Declaration
There are six ways to introduce a name:
* Constant declaration with const
* Type declaration with type
* Variable declaration with var (must have type or initial value, sometimes both)
* short, initailized variable decration of type `:=` (only inside a function)
* Function declaration with func (methods may only be declared at package level)
* Formal parameters and named returns of a function

## Structural typing
It's the same type if it has the same structure or behavior
```Go
a := [...]int{1, 2, 3}
b := [3]int{}

a = b                   // OK

c := [4]int{}

a = c                   // NOT OK
```

Go uses structural typing in most cases
It's the same type if it has the same structure or behavior:
* arrays of the same size and base type
* slices with the same base type
* maps of the same key and values type
* structs with the same sequence of field names/types
* functions with the same parameter & return types

## Named typing
It's only the same type if it has the same declared type name, for example:
```Go
type x int          // althought x is underlyingly an int, it is not logically the same as int, its type of user declared type x

func main() {
    var a x         // x is a defined type; base int

    b := 12         // b defaults to int

    a = b           // TYPE MISMATCH

    a = 12          // OK, untyped literal
    a = x(b)        // OK, type conversion
}
```

Go uses named typing for non-function user-declared types

