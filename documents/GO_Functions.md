# Functions

Functions are "first class" objects; you can:
* Define them -  even inside another function
* Create anonymous function literals
* Pass them as function parameters/return values
* Store them in variables
* Store them in slices and maps (but as keys)
* Store them as fields of a structure type
* Send and receive them in channels
* Write methods against a function type
* Compare a function var gainst nil

Basically it can be treated as a variable, in other languages, functions is typically just defined and can only be called

## Function signatures
The signature of a function is the order & type of its parameters and return values

It does not depend on the names of those parameters or returns
```Go
var try func(string, int) string

func Do(a string, b int) string {
    ...
}

func NotDo(x string, y int)(a string) {
    ...
}
```

These funcitons have the same stuctural type. So basically the every parameter order and type must match if the functions
have the same signature, but its names does not have to match.

## Parameter terms
A function declaration lists **formal parameters** (so a and b are formal parameters)
`func do(a, b int) int {...}`

A function call has **actual parameters** (aka "arguments) (so 1 and 2 are actual parameters)
`result := do(1, 2)`

A parameter is passes **by value** if the function gets a copy;
the caller can't see changes to the copy

A parameter is passed **by reference** if the function can modify the actual parameter such that the caller sees the changes
(done by passing in memory address that points to the underlying location of data allowing the function to change caller values)

By value:
* numbers
* bool
* arrays
* structs

By reference:
* things passed by pointer `(&x)`
* strings (but they're immutable)
* slices
* maps
* channels

### Param example
Below is an example of a function to demonstrate how pass by reference changes the caller map
```Go
package main

import (
	"fmt"
)

func do(m1 map[int]int) {
	m1[3] = 0
	m1 = make(map[int]int)
	m1[4] = 4
	fmt.Println("m1", m1)
}
func main() {
	m := map[int]int{4: 1, 7: 2, 8: 3}

	fmt.Println("m", m)
	do(m)
	fmt.Println("m", m)
}

// Output:
// m map[4:1 7:2 8:3]
// m1 map[4:4]
// m map[3:0 4:1 7:2 8:3]
```
Whats happening here is that `m` is created in main and is a descriptor (with pointer inside) to wherever that map lives.
When it is passed to `do` function, the m1 is another alias descriptor that points to the same underlying map and the change
`m1[3] = 0` is made against the caller map m. However, we then make a new map which returns a different descriptor and m1 
become a descriptor to a whole new map that exists inside the function. Since `m1[4] = 4` is done with new file descriptor, 
the changes happen to the new map and thats why we see in Println `m1 map[4:4]`. Essentially, `m` never gets copied to function
scopre, but instead is just referenced by `m1` within the function allowing it to change the caller map before it gets a new
descriptor.


Below is an example of how reference can be used in another way with different behavior
```Go
package main

import (
	"fmt"
)

// param must change its type to pointer type * in order to hold the mem address of m
func do(m1 *map[int]int) {
    // when indexing, dereference is wrapped in parenthesis
    // m1 hold the address m itself (not the file descriptor), so dereferencing m1
    // gives us the location of m itself, and we are making a change [3] = 0 at the place m is
	(*m1)[3] = 0
    // since dereferenced m1 is m itself, make returns a new descriptor and the descriptor of m gets changes this time
    // In the example above, m1 was a descriptor itself and make reassigned m1 value, but here m1 just holds m address
    // and the changes are refelected to m directly and m is now a new map
	*m1 = make(map[int]int)
    // add key 4 val 4 to new map
	(*m1)[4] = 4
	fmt.Println("m1", *m1)
}
func main() {
    // m is a descriptor of the created map
	m := map[int]int{4: 1, 7: 2, 8: 3}

	fmt.Println("m", m)
    // we pass in the mem address of m itself, not the descriptor this time
	do(&m)
	fmt.Println("m", m)
}
```

The biggest difference between the previous example and this one is that m1 now is a variable that points directly to m whereas
before, m1 was a descriptor itself that pointed to the same underlying map and not the variable m itself. So changes in the
example before could reflect on the map itself but changes directly to m1 would only affect m1. In this exmaple, m1 is now a 
pointer m and we are working directly with m to my understanding and any changes done to m1 is like changing the variable m

### Parameter passing: the ultimate truth
Technically, the parameters in Go are passed by value for reference types such as slices, maps, structs, channels, strings, etc.
The descriptors get copied into the function, but the underlying data is not copied. These copy descriptors then point to the 
same underlying data but the descriptors are ultimately copied. The true pass by references in given in the second example above
where we actually passedin the mem address of the original descriptor and the changes made in the function directly affected
the original descriptor.

Therefore, in Go we think of it as "by reference" when it technically isn't


## Function Return Values
Functions can have multiple return values
```Go
func doIt(a int, b []int) int {
    ...
    return 1
}

func doItAgain(a string) (int, error) {
    ...
    return 1, nil
}
```

Every return statement must have all the values specified and multiple return values are wrapped in parenthesis


## Deferred Execution
How do we make sure someting gets done?
* close a file we opened
* close a socket/HTTP request we made
* unlock a mutex we locked
* make sure somethings gets saved before we're done
* ...

The defer statement captures a function call to run later

### Defer example
We need to ensure the file closes no matter what
```Go
func main() {
    f, err := os.Open("my_file.txt")
    
    if err != nil {
        ...
    }

    defer f.Close()

    // and do someting with the file
}
```
Defer only takes functions and it calls them at function exit. So here, Close is guaranteed to run right after the function exits.
(dont defer closing the file until we know it really happened!)

You can also have multiple defers
```Go
defer a()
defer b()
defer c()
```

But defers work in the order of LIFO stack because everytime a defer is encountered, Go pushed them into stack and executes
them in the order of LIFO by popping it. So at the end of the function, c gets executed, then b, and then finally a

### Defer Scope
We need to ensure the file closes no matter what
```Go
func main() {
    f := os.Stdin

    if len(os.Args) > 1 {
        if f, err := os.Open(os.Args[1]); err != nil {
            ...
        }
        defer f.Close()
    }

    // and do something else with the file
}
```

Defer only executes when the function ends no matter where it is defined. So here, the defer will not execute when we leave 
the if block and only when main ends.

### Defer gotcha #1
The scope of a defer statement is the function
```Go
func main() {
    for i := 1; i < len(os.Args); i++ {
        f, err := os.Open(os.Args[i])
        ...
        defer f.Close()
        ...
    }
}
```
Here the defer is not closing files at the end of every loop, defer only runs once function exits (in this case main). So
the deferred calls to Close must wait unitl function exit and we might run out of file descriptors before that. In this case, 
just do f.Close() at the end of loop instead of relying on defer here.

### Defer gotcha #2
Unlike a closure, defer copies arguments to the deferred call
```Go
func main() {
    a := 10

    defer fmt.Println(a)

    a = 11
    fmt.Println(a)
}
// prints 11, 10
```

Here, the parameter `a` gets copied at the defer statement location when `a=10` (not a reference)

Also, a defer statement runs before the return is "done". Meaning, the return value is calculated at the end of function, but
not yet returned. After the calculation, the defer is evaluated and any changes made by defer is then finally returned by the 
function. Although defer accepts functions, you can create anonymous functions 
to execute tasks that do not necessarily need a function.
```Go
// notice that it is naked return statement and returns a by default
func doIt() (a int) {
    defer func() {
        a = 2
    }()

    a = 1
    return
}
// returns 2
```

We have a named return value and a "naked" return
The deferred anonymous function can update that variable and return the updated value. 
To my understanding, a return value is calculated at the end of funciton execution, and then deferred func runs 
and can update that value. 