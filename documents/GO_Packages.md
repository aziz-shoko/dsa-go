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

# Package Organization

I never really understood Go's package organization, like does every package have main at root? can packages extend to multiple directories? how do subpackages work? etc.
Below is just my notes from reading up some stuff on Go package organization.

## Basic Rule of Thumb: One Directory = One Package
In Go, a directory holds a single package (ofcousrse there are exceptions, like some test patterns). Typically, the package name you declare at the top of each `.go` file 
in that directory is the same for all files in that directory. For example:
```
myrepo/
  main.go          // declares "package main"
  helpers.go       // also declares "package main" if it's in the same dir
  mypkg/
    stuff.go       // declares "package mypkg"
  anotherpkg/
    things.go      // declares "package anotherpkg"
```

## Packages vs Modules
Module: Defined by your `go.mod` file (for example, `module github.com/aziz-shoko/dsa-go`). This is the "root" that Go uses when referencing your code externally or internally.

Packages: Subdirectories within your module that hold Go code for different functionalities

For example: if `go.mod` says:
```Go
module github.com/you/myrepo
```

And you have a directory structure like:
```
myrepo/
  main.go               // package main
  mypkg/
    stuff.go            // package mypkg
  anotherpkg/
    things.go           // package anotherpkg
```

Then in `main.go`, if you want to import `mypkg`, you would have to do:
```Go
import (
    "github.com/you/myrepo/mypkg"
)
```

And ofcourse if `anotherpkg` depended on `mypkg`, then inside `anotherpkg/things.go` you'd do:
```Go
import (
    "github.com/you/myrepo/mypkg"
)
```

## The Special "package main"

"package main" is the package that produces an executable. It must have a `func main()` entry point somewhere in its file

Any package other than main is a library package. It can't be run directly, but can be imported. In other words, you can't just do `go run <gofilename>` 
unless it is specifically `package main` with entry point `func main()`. To run the library package files, you gotta use the `test` package to run tests against it 
and see if they actually work. 

## Making subpackages vs Putting Everything in `package main`

### Subpackages
You typically craete separate directories (subpackages) to keep your code clean.
For example: if you have business logic or utility functions, you don't want them all cluttered in `main.go`. Instead you would have `main.go` in the root under `package main`
and break other specific components into its packages in directories and import them into `main.go` and run it like that for best organization

### Multiple Directories But Same Package Name
Its possible to have files in multiple directories under all in the same package but it is not typical and is not recommended. You rarely see `package main` repeated in multiple
subdirectories. If you do, they have to compile together as a single main package, which is weird and can get messy.

### Everything in package main
For small projects, you just have everything in package main. This is not a good approach when your code grows because it will lack modular organization.

## Referencing Packages

### Within the same module:
Use the **module path + subdirectory** convention. For example, if your module is github.com/you/myrepo, and you have a subpackage in mypkg/,
the import path is github.com/you/myrepo/mypkg.

### External Imports
If you import some external library (as in pulling in code from seperate module outside your code), you refer to that library's module path. For example:
```Go
import "github.com/go-sql-driver/mysql"
```

## Package Interdependency
As mentioned earlier in the document, **Circular Imports Are Not Allowed**.

Go does not permit circular imports. If package A imports package B, then package B cannot import package A.
So design your code so that dependencies flow one way. If you hit a circular import issue, you typically need to 
extract common pieces into a third package that both A and B can import.

## Typical Project Layout Example
A more complete example might look like this:
```
myrepo/
  go.mod             // module github.com/you/myrepo
  cmd/
    myapp/
      main.go        // package main (entry point)
  internal/
    database/
      db.go          // package database
    config/
      config.go      // package config
  pkg/
    helpers/
      helpers.go     // package helpers
  // etc.
```

* The `cmd/myapp/` folder is where the "main" package lives. That's your actual executable.
* `internal` contains packages meant to be private to your module (a Go convention)
* `pkg` contains "public" or well-organized library code that can be imported by other packages (and potentially other modules)
  
Then in `cmd/myapp/main.go`, you might do:
```Go
package main

import (
    "github.com/you/myrepo/internal/database"
    "github.com/you/myrepo/internal/config"
    "github.com/you/myrepo/pkg/helpers"
)

func main() {
    cfg := config.Load()           // just an example
    dbConn := database.Connect(cfg)
    result := helpers.DoStuff(dbConn)
    // ...
}
```

