# Customizing Errors
Most of the time, errors are just strings
```go
func (h HAL9000) OpenPodBayDoors() error {
    if h.kill {
        return fmt.Errorf("I'm sorry %s, I can't do that", h.victim)
    }
}
```

## Error Types
Errors in Go are objects satisfying the error interface:
```Go
type error interface {
    Error() string
}

// Any concrete type with Error() can represent an error
type Fizgig struct{}

func (f Fizgig) Error() string {
    return "Your fizgig is bent"
}
```

So the basic idea in Go with this approach is how can we basically add more data to
errors so that its easier to debug or something (i think)

## A custom error type
Here is an example of building out a custom error type
```go
type errKind int

// Go's way of doing enum
const (
    _   errKind = iota  // so don't start at 0
    noHeader
    cantReadHeader
    invalidHdrType
    ...
)

type WaveError struct {
    kind    errKind
    value   int
    err     error
}

// We use different formats depending on the situation 
func (e WaveError) Error() string {
    switch e.kind {
    case noHeader:
        return "no header (file too short?)"
    case cantReadHeader:
        return fmt.Sprintf("can't read header[%d]: %s", e.value, e.err.Error())
    case invalidHdrType:
        return "invalid header type"
    case invalidChkLength:
        return fmt.Sprintf("invalid chunk length: %d", e.value)
    }
}

// Can have a couple of helper methods to generate errors

// with returns an error with a particular value (e.g., header type)
func (e WaveError) with(val int) WaveError {
    e1 := e
    e1.value = val
    return e1
}

// from returns an error with a particular location and 
// underlying error (e.g., from the standard library)
func (e WaveError) from(pos int, err error) WaveError {
    e1 := e
    e1.value = pos
    e1.err = err
    return e1
}

// And we have some prototype errors we can return or customize
var (
    HeaderMissing          = WaveError{kind: noHeader}
    HeaderReadFailed       = WaveError{kind: cantReadHeader}
    InvalidHeaderType      = WaveError{kind: invalidHdrType}
    ...
)

// Here's an example of those errors in use
func DecodeHeader(b []byte) (*Header, []byte, error) {
    var err error
    var pos int

    header := Header{TotalLength: uint32(len(b))}
    buf := bytes.NewReader(b)

    if len(b) < headerSize {
        return &header, nil, HeaderMissing
    }

    if err = binary.Read(buf, binary.BigEndian, &header.riff); err != nil {
        // prototype paradigm (IMPORTANT)
        return &header, nil, HeaderReadFailed.from(pos, err)
    }
    ...
}
```

## Wrapped Errors
Starting with Go 1.13, we can wrap one error in another
```Go
func (h HAL9009) OpenPodBayDoors() error {
    ...
    if h.err != nil {
        return fmt.Errorf("I'm sorry %s, I cant: %w", h.victim, h.err)
    }
    ...
}
```
The easiest way to do that is to use the `%w` format verb with fmt.Errorf()

Wrapping errors gives us an error chain we can unravel
`top-level error` -> `intermediate error` -> `original error`

```go
package main

import (
    "fmt"
    "errors"
)

// Level 1: Deep down in your code, something basic fails
func readConfigFile() error {
    // This simulates a low-level error (like from the OS)
    return errors.New("permission denied")
}

// Level 2: A higher function calls it and adds context
func loadUserSettings() error {
    err := readConfigFile()
    if err != nil {
        // WRAP the error with more context
        return fmt.Errorf("failed to load config: %w", err)
        //                                        ↑ %w wraps the error
    }
    return nil
}

// Level 3: An even higher function adds more context  
func initializeApp() error {
    err := loadUserSettings()
    if err != nil {
        // WRAP again with even more context
        return fmt.Errorf("app startup failed: %w", err)
    }
    return nil
}

func main() {
    err := initializeApp()
    if err != nil {
        fmt.Println(err)
        // Output: "app startup failed: failed to load config: permission denied"
    }
}

// In other words:
// DON'T DO THIS (loses the original error):
return fmt.Errorf("failed to load config: %s", err)

// DO THIS (wraps and preserves the original error):
return fmt.Errorf("failed to load config: %w", err)
//                                        ↑ This keeps the original error
```

Back to the WaveError example, with wrapped errors, you can have your Custom error types unwrap their 
internal errors.
```Go
type WaveError struct {
    kind errKind
    value int
    err error
}

// So the wrapped error can be set in err error field and we can use a method to unwrap it
func (w *WaveError) Unwrap() error {
    return w.err
}
```

Now, for unwrapping the errors, there are some useful utilities to help with that.
1. `errors.Is` 
We can check whether an error has another error in its chain

`errors.Is` compares with an error variable, not a type
```Go
if audio, err = DecodeWaveFile(fn); err != nil {
    if errors.Is(err, os.ErrPermission) {
        // let's report a security violation
        ...
    }
    ...
}
```

2. `errors.As`
We can get an error of an underlying type if it's in the chain
`errors.As` looks for an error type, not a value
```Go
...
if audio, err = DecodeWaveFile(fn); err != nil {
    var e os.PathError  // a struct
    
    if errors.As(err, &e) {
        // let's pass back the underlying file error
        return e
    }
}
...
```

## Philosophy of Errors in Go
When it comes to errors, Go devs fall into one of these categories:
1. they hate constantly writing if/else blocks or,
2. they think writing if/else blocks makes things clearer
3. they dont care and just care about writing code

But, there two kinds of errors:
1. Normal errors

Normal errors result from input or external conditions (for example, wrong user input(input) or is the network available(external))
Go handles this case by returning the error type
```Go
// Not exactly os.Open, but shows the basic logic
func Open(name string, flag int, perm FileMode) (*File, error) {
    r, e := syscall.Open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))

    if e != nil {
        // notice how PathError is a struct
        return nil, &PathError{"open", name, e}
    }

    return newFile(uintptr(r), name, kindOpenFile), nil
}
```

2. Abnormal errors

The other kind of error is Abnormal errors that result from invalid program logic (for example, a nil pointer)

For program logic errors, Go code does a panic
```Go
func (d *digest) checkSum() [Size]byte {
    // finishing writing the checkSum
    ...
    if d.nx != 0 {          // panic is there's a data left over
        panic("d.nx != 0")
    }
    ...
}
```
Basically, you write codes to not have a bug and if you have a bug (invalid program logic), then Go's panics and 
crashes the program. It crashes because, well you have to go fix that bug and rerun the program so that it doesn't have 
a bug. 

So overall, when your program has a logic bug:
"Fail Hard, Fail Fast"

In more detail:
If you server crashes, it will get immediate attention
* logs are often noisy 
* so proactive log searches for "problems" are rare

We want evidence of the failure as close as possible in time and space to the orignial defect in the code
* connect the crash to logs that explain the context
* traceback from the point closest to the broken logic 

In a distributed system, crash failures are the safest type to handle
* it's better to die than to be a zombie or babble or corrupt the DB
* not crashing may lead to Byzantine failures

### When should we panic?
Only when the error caused by our own programming defect, e.g. 
* we can't walk a data structure we built
* we have an off-by-one bug encoding bytes

In other words,
**Panic should be used when our assumptions of our own programming design or logic are wrong**
These cases mgiht use an "assert" in other programming languages

So for example:
A B-tree data structure satisfies several invariants:
1. every path from the root to a leaf has the same length
2. if a node has n children, it contains n - 1 keys
3. every node (except the root) is at least half full
4. the root has at least two children if it is not a leaf
5. subnode keys fall between the keys of the parent node that lie on either side of hte subnode pointer 

If any of these rules (assumptions in logic) is ever false, the B-tree should panic!

### Exception handling
Although Go panics for Abnormal errors, other programming languages have exception handling, which is basically
letting the program do something else when it does catch an abnormal error (sometimes good, for example, if the
airplane's engine stops working, we dont want to stop the plane immedietly, but let it run as long as possible 
until safe landing.)

Excepting handling was popularized to allow "graceful degradation" of safety-critical systems (e.g., Ada 
and flight control software)

Ironically, the most safety-critical systems are built without using exceptions! 

Exception handling introduces invisible control paths through code
So code with exceptions is harder to analyze (automatically or by eye)

Officially, Go doesn't support exception handling as in other languages
But practically, it does - in the form of `panic & recover`

panic in a function will still cause deferred function calls to run

Then it will stop only if it finds a valid recover call in a defer as it unwinds the stack

Recovery from panic only works inside defer:
```Go
func abc() {
    panic("omg")
}

func main() {
    defer func() {
        // can only do recover inside defer
        if p := recover(); p != nil {
            // what else can you do?
            fmt.Println("recover:", p)
        }
    }

    abc()
}
```

Panic and recover is typically not used in Go (mostly for Unit Tests) because we want Go to handle erors the Go 
way and not bring in other languages philosophies like exceptions

### Define errors out of existence (How to handle errors)
Error (edge) cases are one of the primary sources of complexity

The best way to deal with many errors is to make them impossible

Design your abstractions so that most (or all) operations are safe:
* reading from a nil map
* appending to a nil slice
* deleting a non-existant item from a map
* taking hte length of a uninitialized string

Try to reduce edge cases that are hard to test or debug (or even thing about!)

### Proactively Prevent Problems
Every piece of data in your software should start life in a valid state

Every transformation should leave it in a valid state
* break large programs into small pieces you can understand
* hide information to reduce the chance of corruption
* avoid clever code and side effects
* avoid unsafe operations
* asert your invariants
* never ignore errors
* test, test, test

Never, ever accept input from a user (or environment) without validation