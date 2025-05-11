# Interfaces
An interface specifies abstract behavior in terms of methods. 
```
type Stringer interface {
    String() string
}
```

Concrete types offer methods that satisfy the interface. In Go, since any types can have methods, any type 
can implement(satisfy) the interface. In other words, interface is like a contract that just shows what it can do. 
But the actual logic for the implementations of what an interface can do is provided by concrete type's methods in the 
background. (But techinically, we dont say implements in Go, because it happens automatically. If a type satisfies an interface,
we say it **is a** member of that interface)

## But what are methods?
A method is a special type of function (syntax from Oberon-2). It has a receiver parameter before the function name parameter.
```Go
type IntSlice []int

// (is IntSlice) part tells that the function String belows to type IntSlice
func (is IntSlice) String() string {
    ...
}
```

type assertions
interface dispatch or method resolution

## Interfaces Philosophy
General advice on how interfaces should be done according to [docs](https://go.dev/wiki/CodeReviewComments#interfaces):

### Explanation of the philosophy above

Terminology Clarification:
* Consumer: The package/code that uses an interface to perform operations. It relies on the behavior defined by the interface.
* Producer: The package/code that creates concrete implementations of an interface.
* Concrete type: A specific implementation (struct, pointer) as opposed to an interface.

The Key Principle:
The core principle in Go's interface design is: Define interfeaces where they are used, not where they are implemented (distinction between using interface and implementing)

Example: File Operations
Let's say we're building file processing system:
```Go
// fileprocessor/processor.go (CONSUMER)
package fileprocessor

// DataReader is defined in the consumer package
type DataReader interface {
	Read() ([]byte, error)
	Close() error
}

// ProcessData uses the interface
func ProcessData(reader DataReader) ([]byte, error) {
	data, err := reader.Read()
	if err != nil {
		return nil, err
	}

	defer reader.Close()

	// Process data...
	return processedData, nil
}
```

Now in the other pacakge
```Go
package filesource

import (
	"os"
)

// FileSource is a concrete implementation
type FileSource struct {
	file *os.File
}

func NewFileSource(path string) (*FileSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FileSource{file: f}, nil
}

func (fs *FileSource) Read() ([]byte, error) {
	// Implementation
	return data, nil
}

func (fs *FileSource) Close() error {
	return fs.file.Close()
}
```

Main usage example
```Go
// main.go (USAGE)
package main

import (
    "fileprocessor"
    "filesource"
)

func main() {
    source, _ := filesource.NewFileSource("data.txt")
    result, _ := fileprocessor.ProcessData(source)
    // Use result...
}
```

Why This Approach Is Better

1. Decoupling: fileprocessor only knows about the behavior it needs, not implementation details.
2. Flexibility: The producer (filesource) can add new methods to FileSource without breaking consumers.
3. Testing: Consumer can easily create mock implementations.

What the Documentation Is Warning Against
The documentation warns against defining interfaces in the producer package and then returning them:
```Go
// BAD APPROACH
// filesource/source.go
package filesource

type DataReader interface {
    Read() ([]byte, error)
    Close() error
}

type fileSource struct { /* ... */ }
func (fs *fileSource) Read() ([]byte, error) { /* ... */ }
func (fs *fileSource) Close() error { /* ... */ }

// Returns interface instead of concrete type
func NewDataReader(path string) (DataReader, error) {
    // ...
    return &fileSource{/* ... */}, nil
}
```

This approach creates problems:

1. It forces all consumers to use this specific interface
2. If the producer wants to add methods, it must modify the interface (breaking all consumers)
3. It makes testing harder as consumers can't easily create their own implementations

In Practice
The recommended pattern is:
1. Consumers define interfaces with only the methods they need
2. Producers return concrete types with all their methods
3. The Go compiler automatically handles the type conversion when a concrete type is used where an interface is expected

## Why interfaces?
Without interfaces, we'd have to write (many) functions
for (many) concrete types, possibly coupled to them
```Go
func OutputToFile(f *File, ...) {...}
func OutputToBuffer(b *Buffer, ...) {...}
func OutputToSockets(s *Socket, ...) {...}
```
Better, we want to define our function in terms of abstract behavior
```Go
type Writer interface {
	Write([]byte) (int, error)
}

func OutputTo(w io.Writer, ...) {...}
```

Another useful example of interfaces in Go
```Go
package main

import (
	"fmt"
	"io"
	"os"
)

// user named type
type ByteCounter int

// reference receiver
func (b ByteCounter) Write(p []byte) (int, error) {
	l := len(p)
	// ByteCounter(l) is type casting from int to named type ByteCounter
	*b += ByteCounter(l)
	return l, nil
}

func main() {
	// Where does the state of c live?
	// The variable c lives in your program's memory. When you pass &c to create f2, you're passing the address of c. The Write
	// method receives this address as the receiver parameter b, and then uses *b to deference and modify the actual value at that address
	var c ByteCounter

	f1,_ := os.Open("a.txt")
	// Why &c and not c?
	// The Write method is defined on *ByteCounter, not on ByteCounter. This is because the methods needs to modify
	// the value (with *b += ByteCounter(l)). In Go, methods that need to modify the receiver must use pointer receivers
	f2 := &c

	n,_ := io.Copy(f2, f1)

	fmt.Println("copied", n, "bytes")
	fmt.Println(c)
}
```

## Interfaces and substitution
All the methods must be present to satisfy the interface

```Go
var w io.Writer
// Has read, write, and close methods
var rwc io.ReadWriteCloser

w = os.Stdout					// OK: *os.File has Write method
w = new(bytes.Buffer)			// OK: *bytes.Buffer has Write method
w = time.Second					// ERROR: no Write method

rwc = os.Stdout					// OK: *os.File has all 3 methods
rwc = new(bytes.Buffer) 		// ERROR: no Close method

w = rwc 						// OK: io.ReadWriteCloser has Write
rwc = w							// ERROR: no Close method
```
This is why it pays to keep interfaces small

## Interface Satisfiability
The receiver must be of the right type (pointer or value)

```Go
type IntSet struct { /*...*/}

func (*IntSet) String() string

var _ = IntSet{}.String() // ERROR: String needs the memory address of IntSet, but here IntSet doesn't have mem address

var s IntSet
var _ = s.String()		  // OK: s is a variable: &s used automatically

var _ fmt.Stringer = &s   // OK
var _ fmt.Stringer = s	  // ERROR: no String method
```

## My Library system example and what i learned from it
Below is the original code i wrote with no interfaces for a simple library system. 
```Go
package main

import (
    "fmt"
    "os"
)

type Books struct {
    Available bool
    User string
    Title string
    Author string
    PageCount int
}

func (b *Books) CheckOut(borrower string) {
    b.Available = false
    b.User = borrower
}

func (b *Books) ReturnBook(title string) {
    if title == b.Title {
        b.Available = true
    } else {
        fmt.Fprintf(os.Stderr, "Book Title %s not found", title)
    }
}

func (b *Books) Status(title string) bool {
    if b.Available {
        fmt.Printf("The book %s is available ", b.Title)
        return b.Available
    } else {
        fmt.Printf("The Book %s is not available ", b.Title)
        return false
    }
}

type DVDs struct {
    Available bool
    User string
    Title string
    Director string
    Runtime int
}

func (b *DVDs) CheckOut(borrower string) {
    b.Available = false
    b.User = borrower
}

func (b *DVDs) ReturnDVDs(title string) {
    if title == b.Title {
        b.Available = true
    } else {
        fmt.Fprintf(os.Stderr, "DVD Title %s not found ", title)
    }
}

func (b *DVDs) Status(title string) bool {
    if b.Available {
        fmt.Printf("The DVD %s is available ", b.Title)
        return b.Available
    } else {
        fmt.Printf("The DVD %s is not available ", b.Title)
        return false
    }
}

type Magazines struct {
    Available bool
    User string
    Title string
    IssueNumber int
    PublicationDate string
}

func (b *Magazines) CheckOut(borrower string) {
    b.Available = false
    b.User = borrower
}

func (b *Magazines) ReturnMagazine(title string) {
    if title == b.Title {
        b.Available = true
    } else {
        fmt.Fprintf(os.Stderr, "Magazine Title %s not found", title)
    }
}

func (b *Magazines) Status(title string) bool {
    if b.Available {
        fmt.Printf("The Magazine %s is available ", b.Title)
        return b.Available
    } else {
        fmt.Printf("The Magazine %s is not available ", b.Title)
        return false
    }
}

func main() {
    bookLibrary := []Books{
        {true, "", "Moby Dick", "Herman Melville", 100},
        {true, "", "Harry Potter", "J.K. Rowling", 300},
    }

    dvdLibrary := []DVDs{
        {true, "", "Transformers 1", "Michael bay", 90},
        {true, "", "Transformers 2", "Michael bay", 110},
    }

    magazineLibrary := []Magazines{
        {true, "", "Mens Health", 123, "Nov 12, 2024"},
        {true, "", "Womans Health", 645, "Nov 15, 2024"},
    }

    // Simulate some checkout interaction with library
    bookLibrary[0].CheckOut("Sam")
    dvdLibrary[1].CheckOut("Jill")
    magazineLibrary[0].CheckOut("David")

    // should say false
    fmt.Println(bookLibrary[0].Status(bookLibrary[0].Title))
    fmt.Println(dvdLibrary[1].Status(dvdLibrary[1].Title))
    fmt.Println(bookLibrary[0].Status(magazineLibrary[0].Title))

    // Simulate some return interaction
    bookLibrary[0].ReturnBook(bookLibrary[0].Title)
    dvdLibrary[1].ReturnDVDs(dvdLibrary[1].Title)
    magazineLibrary[0].ReturnMagazine(magazineLibrary[0].Title)

    // should say true
    fmt.Println(bookLibrary[0].Status(bookLibrary[0].Title))
    fmt.Println(dvdLibrary[1].Status(dvdLibrary[1].Title))
    fmt.Println(bookLibrary[0].Status(magazineLibrary[0].Title))
}
```

As we can see the code above absolutely sucks and is repetitive with no real world logic intuition behind it (its not like a real library
where we can go and check out stuff from one place, here we are going to different structs and calling their specific methods). 

Below is the interface iteration i did to decouple the consumers from the producers (as in decouple users of interface from the ones that satisfy it)
```Go
package main

import (
	"fmt"
)

type Books struct {
	Available bool
	User string
	Title string
	Author string
	PageCount int
}

func (b *Books) CheckOut(borrower string) {
	b.Available = false
	b.User = borrower
}

func (b *Books) Return() {
	b.Available = true
}

func (b *Books) Status() bool {
	if b.Available {
		fmt.Printf("The book %s is available ", b.Title)
		return b.Available
	} else {
		fmt.Printf("The Book %s is not available ", b.Title)
		return false
	}
}

type DVDs struct {
	Available bool
	User string
	Title string
	Director string
	Runtime int
}

func (b *DVDs) CheckOut(borrower string) {
	b.Available = false
	b.User = borrower
}

func (b *DVDs) Return() {
	b.Available = true
}

func (b *DVDs) Status() bool {
	if b.Available {
		fmt.Printf("The DVD %s is available ", b.Title)
		return b.Available
	} else {
		fmt.Printf("The DVD %s is not available ", b.Title)
		return false
	}
}

type Magazines struct {
	Available bool
	User string
	Title string
	IssueNumber int
	PublicationDate string
}

func (b *Magazines) CheckOut(borrower string) {
	b.Available = false
	b.User = borrower
}

func (b *Magazines) Return() {
	b.Available = true
}

func (b *Magazines) Status() bool {
	if b.Available {
		fmt.Printf("The Magazine %s is available ", b.Title)
		return b.Available
	} else {
		fmt.Printf("The Magazine %s is not available ", b.Title)
		return false
	}
}

type LibraryItem interface {
	CheckOut(string)
	Return()
	Status() bool
}

func checkOut(b LibraryItem, user string) {
	b.CheckOut(user)
}

func returnItem(b LibraryItem) {
	b.Return()
}

func getStatus(b LibraryItem) bool {
	return b.Status()
}

func main() {
	// bookLibrary := []Books{
	// 	{true, "", "Moby Dick", "Herman Melville", 100},
	// 	{true, "", "Harry Potter", "J.K. Rowling", 300},
	// }

	// dvdLibrary := []DVDs{
	// 	{true, "", "Transformers 1", "Michael bay", 90},
	// 	{true, "", "Transformers 2", "Michael bay", 110},
	// }

	// magazineLibrary := []Magazines{
	// 	{true, "", "Mens Health", 123, "Nov 12, 2024"},
	// 	{true, "", "Womans Health", 645, "Nov 15, 2024"},
	// }
	Library := []LibraryItem{
		&Books{true, "", "Harry Potter", "J.K. Rowling", 300},
		&Books{true, "", "Moby Dick", "Herman Melville", 100},
		&DVDs{true, "", "Transformers 1", "Michael bay", 90},
		&DVDs{true, "", "Transformers 2", "Michael bay", 110},
		&Magazines{true, "", "Mens Health", 123, "Nov 12, 2024"},
		&Magazines{true, "", "Womans Health", 645, "Nov 15, 2024"},
	}

	// Simulate some checkout interaction with library
	checkOut(Library[0], "sam")
	checkOut(Library[2], "jill")
	checkOut(Library[4], "jack")

	// should say false
	fmt.Println(getStatus(Library[0]))
	fmt.Println(getStatus(Library[2]))
	fmt.Println(getStatus(Library[4]))

	// Simulate some return interaction
	returnItem(Library[0])
	returnItem(Library[2])
	returnItem(Library[4])

	// should say true
	fmt.Println(getStatus(Library[0]))
	fmt.Println(getStatus(Library[2]))
	fmt.Println(getStatus(Library[4]))

	// messing around with type assertions
	if book, ok := Library[0].(*Books); ok {
		fmt.Println(book.Author)
	}
}
```

Most important things I learned from writing the code above for idiomatic Go.
* Interfaces represent "complete behaviors", as in LibraryItem represents the whole interaction with the Library.
Initially I thought interfaces represented "behaviors" in the sense that represents a single specific method but it doesn't. 
LibraryItem holds all the possible interactions with Library and any user of this interface (like functions) can then use its
specific methods to whatever it needs

* In Go, you can create a slice (or other collection, like maps?) of an interface type, and that slice can hold any value that
implements the interface. For example, Books, DVDs, and Magazines were all of different struct type but all implement the LibraryItem
interface so we can consolidate them under a single Library slice (this is more intuitive, its like going to one library for everything now)

* Although you can have a slice of interface type, an interface only holds the methods, so the attributes of the structs cannot be accessed
directly. To access type specific properties, you need to use type assertions. Basically type assertion is sort of like casting to the original
struct type so that you can access its attributes. The if statement example is show in the code above

Below is further refinement after learning compositions.
```Go
package main

import (
	"fmt"
)

type LibraryResource struct {
	Available bool
	User      string
	Title     string
}

func (b *LibraryResource) CheckOut(borrower string) {
	b.Available = false
	b.User = borrower
}

func (b *LibraryResource) Return() {
	b.Available = true
}

func (b *LibraryResource) Status(itemType string) bool {
	if b.Available {
		fmt.Printf("The %[2]s %[1]s is available ", b.Title, itemType)
		return b.Available
	} else {
		fmt.Printf("The %[2]s %[1]s is not available ", b.Title, itemType)
		return false
	}
}

type Books struct {
	LibraryResource
	Author    string
	PageCount int
}

func (b *Books) Status() bool {
	return b.LibraryResource.Status("book")
}

type DVDs struct {
	LibraryResource
	Director string
	Runtime  int
}

func (b *DVDs) Status() bool {
	return b.LibraryResource.Status("dvd")
}

type Magazines struct {
	LibraryResource
	IssueNumber     int
	PublicationDate string
}

func (b *Magazines) Status() bool {
	return b.LibraryResource.Status("magazine")
}

type LibraryItem interface {
	CheckOut(string)
	Return()
	Status() bool
}

func checkOut(b LibraryItem, user string) {
	b.CheckOut(user)
}

func returnItem(b LibraryItem) {
	b.Return()
}

func getStatus(b LibraryItem) bool {
	return b.Status()
}

func main() {
	Library := []LibraryItem{
		&Books{LibraryResource{true, "", "Harry Potter"}, "J.K. Rowling", 300},
		&Books{LibraryResource{true, "", "Moby Dick"}, "Herman Melville", 100},
		&DVDs{LibraryResource{true, "", "Transformers 1"}, "Michael bay", 90},
		&DVDs{LibraryResource{true, "", "Transformers 2"}, "Michael bay", 110},
		&Magazines{LibraryResource{true, "", "Mens Health"}, 150, "Nov 12, 2024"},
		&Magazines{LibraryResource{true, "", "Womands Health"}, 645, "Nov 15, 2024"},
	}

	// Simulate some checkout interaction with library
	checkOut(Library[0], "sam")
	checkOut(Library[2], "jill")
	checkOut(Library[4], "jack")

	// should say false
	fmt.Println(getStatus(Library[0]))
	fmt.Println(getStatus(Library[2]))
	fmt.Println(getStatus(Library[4]))

	// Simulate some return interaction
	returnItem(Library[0])
	returnItem(Library[2])
	returnItem(Library[4])

	// should say true
	fmt.Println(getStatus(Library[0]))
	fmt.Println(getStatus(Library[2]))
	fmt.Println(getStatus(Library[4]))

	// messing around with type assertions
	if book, ok := Library[0].(*Books); ok {
		fmt.Println(book.Author)
	}
}
```
Here, the most important thing I learned is that Go likes to follow a "has a" relationship. So rather than have a hierarchical inheritance type of structure, it follows a structure where objects 
just access each other for the data it needs. Here, a book, dvd, magazine all have a book, user, title fields, so it can be encansulated under one struct and can have other structs use its attributes/methods.
This reduced the repetitive attributes and methods shared by all three original structs of book, dvd and magazine. 


## Nil Interfaces
An interface variable is nil until initialized

It really has two parts (interface variable is also a descriptor that contains two pointers):
* a value or pointer of some type (A pointer to the concrete value (like a struct) or the copy of the value itself)
* a pointer to type information (known as interface talbe/itable) so the correct actual method can be identified

```Go
var r io.Reader 		// nil until initialized
var b *bytes.Buffer 	// ditto (ditto means same thing)

r = b 					// r is no longer nil!
						// but it has a nil pointer to a Buffer
						// Basically, r now contains:
						// - a nil pointer value
						// a non-nil pointer to *bytes.Buffer's interface table
						// so r != nil would evaluate to true
```

This two-pointer representation is why `r = b` results in a non-nil interface even though `b` itself is nil. 
Precisely:
* The interface descriptor's type pointer is non-nil (points to *bytes.Buffer's itable)
* The interface descriptor's value pointer is nil (doesn't point to an actual Buffer)

The interface descriptor contains: 
1. Value pointer: Points to the actual data of the concrete type(or contains the data itself for the small values)
2. Type pointer: Points to the interface table(itable) which contains:
	* Type information about the concrete type
	* Pointers to the method implementations for this specific concrete type
	* Runtime type information for reflection

So when you call a method on an interface, Go uses the type information to find the correct implementation to call. 
Even if the value pointer is nil, the methods can still get called and panic can only happen if a method tries to access the value itself.

## Error type in Go
In Go, `error` is a special type, but its really an interface
```Go
type error interface {
	func Error() string // a function that handles the error and returns string message of the error
}
```
We can compare it to nil unless we make a mistake
The mistake is to store a nil pointer to a concrete type in the error variable

For example:
```Go
package main

import (
	"fmt"
)

type errFoo struct {
	err error
	path string
}

func (e errFoo) Error() string {
	return fmt.Sprintf("%s: %s", e.path, e.err) 
}

func XYZ(a int) *errFoo { // instead of returning the concrete type, should return 'error'
	return nil
}

func main() {
	var err error = XYZ(1) 	// BAD: interface gets a nil concret ptr

	if err != nil {
		fmt.Println("oops")
	} else {
		fmt.Println("OK!")
	}
}
```

As mentioned in section above, a nil assignemnt to nil interface makes the interface not nil.
So, because the error is an interface, we dont return concrete types for errors because then the 
error interface is none nil and can trigger the err message when comparing `err != nil`. 
Therefore, it is Go idiom to always return the interface type error from functions and not concrete 
error types. 

## Pointer vs Value receivers 
A method can be defined on a pointer to a type

```Go
type Point struct {
	x,y float32
}

func (p *Point) Add(x, y float32) { // Pointer receiver changes the struct values itself directly
	p.x, p.y = p.x + x, p.y + y		
}

func (p Point) OffSetOf(p1 Point) (x float32, y float32) { // Value reciever can read values from struct but any changes don't change the struct value itself
	x, y = p.x - p1.x, p.y - p1.y
	return // nakes return of x and y
}

// Pointer methods may be called on non-pointers and vice versa
// Go automatically use * or & as needed

p1 := new(Point) 		// *Point, at (0.0)
p2 := Point{1, 1}		// holds value of Point{1,1} struct directly, not a descriptor

p1.OffSetOf(p2)			// same as (*p1).OffSetOf(p2)
p2.Add(3, 4)			// same as (&p2).Add(3, 4)
```

But, & may only be applied to objects that are addressable (meaning anything you can obtain the memory address of using &)

### Compatability between objects and receiver types

| Compatability   | Pointer | L-Value | R-Value |
|-----------------|---------|---------|---------|
|pointer receiver | OK		| OK&	  | Not OK	|
|value receiver   | OK*		| OK      | OK		|

* L-value: An expression that can appear on the left side of an assignemtn. These are
addressable and have a memory location. Varialbes like `var p Point` are L-Values
* R-Value: An expression that can only appear on the right side of an assignment. These are
temporary values that don't have persistent memory address. Literals and composite literals like `Point{1,1}` are R-values

A method requiring a pointer receiver may only be called on an addressable object
```Go
var p Point

p.Add(1, 2)				// OK, &p
Point{1, 1}.Add(2,3)	// Not Ok, can't take address
```

Breaking Down the Table: 
Pointer Receiver (func (p *Point) Method())
* Pointer: ✓ OK
  * When you call ptr.Method() where ptr is a pointer like *Point, it works directly.
* L-Value: ✓ OK (with automatic address-taking)
  * When you call p.Method() where p is an L-value like a variable, Go automatically takes the address with &p.
* R-Value: ❌ Not OK
  * When you call Point{1,1}.Method(), it fails because you can't take the address of a temporary value or literal.
  * You can't do &Point{1,1}.Method() directly (though you could assign to a variable first).

Value Receiver (func (p Point) Method())
* Pointer: ✓ OK (with automatic dereferencing)
  * When you call ptr.Method() where ptr is a pointer, Go automatically dereferences with *ptr.
* L-Value: ✓ OK
  * When you call p.Method() where p is an L-value, it works directly.
* R-Value: ✓ OK
  * When you call Point{1,1}.Method(), it works because value receivers accept any value of the correct type.

## Consistency in receiver types
If one method of a type takes a pointer receiver, then all its methods shuold take pointers*

and in general objects of that type are probably not safe to copy.
```Go
type Buffer struct {
	buf []byte
	off int
}

func (b *Buffer) ReadString(delim byte) (string, error) {
	...
}
```
But of course there are cases when they shouldn't for reasons

## Currying Functions
Currying functions come from functional programming. Currying takes a ufnction and reduces its argument
count by one (one argument gets bound, and a new function is returned)

```Go
func Add(a, b int) int {
	return a + b
}

func AddToA(a int) func(int) int {
	return func(b int) int {
		return Add(a, b)
	}
}

addTo1 := AddToA(1)
fmt.Println(Add(1,2) == AddTo1(2)) // true
```

## Method Values
A method value is a method closed over its receiver. 

A selected method may be passed similar to a closure;
the receiver is closed over at that point

```Go
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X - p.X, q.Y - q.X)
}

p := Point{1, 2}
q := Point{4, 6}

distanceFromP := p.Distance		// this is a method value
fmt.Println(distanceFromP(q)) 	// and can be called later
```

## Interfaces in Practice
1. Let consumers define interfaces (what minimal behavior do they require?)
2. Re-use standard interfaces wherever possible
3. Keep interface declarations small (The bigger the interface, the weaker the abstraction)
4. Compose one-method interfaces into larger interfaces (if needed)
5. Avoid coupling interfaces to particular type/implementations
6. Accept interfaces, but return concrete types (let the consumer of the return type decide how to use it)
"Be liberal in what you accept, be conservative in what you return"
* Put the least restriction on what parameters you accept (the minimal interface)
Dont require ReadWriteCloser if you only need to read

* Avoid restricting the use of your return type (the concrete value you return might fit with many interfaces!)
Returning *os.File is less restrictive than returning io.ReadWriteCloser because files have other useful methods

Returning error is a good example of an 

## Empty Interface
The interface{} has no methods

So it is satisfied by anything! (so it represents anything)

Empty interfaces are commonly used; they're how the formatted I/O routines can print any type
```Go
func fmt.Printf(f string, args ...interface{})
```

Reflection is needed to determine what the concrete type is. (i will learn reflection later)