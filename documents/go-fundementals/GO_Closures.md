# Closures in Go

Before getting into Go, lets discuss what heaps and stacks are in memory. Then talk about variable scopes vs their lifetimes

## Stack vs Heap Memory
The stack and heap are two region of memory used during program execution with different characteristics and purposes.

### Stack Memory
Stack is a region of memory that stores temporary variables created by each function during execution
* How it works: Variables are allocated in LIFO (Last In, First Out) order
* Allocation: Memory allocation is automatic when you declare variables 
* Deallocation: Memory is automatically freed when variables go out of scope
* Size: Usually fixed and relatively small (often MB range)
* Speed: Very fast allocation/deallocation (just pointer movement)
* Content: Local variables, function parameters, return addresses

Pros of Stack:
* Extremely fast allocation/deallocation
* Memory management is automatic
* No fragmentation issues
* Predictable memory usage

Cons of Stack:
* Limited size (stack overflow if exceeded)
* Variables can't live beyond their scope/function
* Can't resize variables after declaration
* Not suitable for large data structures

### Heap Memory
Heap is a region of memory used for dynamic allocation, where variables can be allocated and freed in any order
* How it works: Memory can be allocated at runtime
* Allocation: Manually requested (using `new`, `malloc`, etc)
* Deallocation: Must be manually freed in many languages (or handled by garbage collection like in Go)
* Size: Much larger than stack (often in GB range)
* Speed: slower than stack operations
* Content: Dynamically allocated data, objects that need to persist beyond function calls

Pros of Heap
* Flexible size limited only by system memory
* Variables can persist beyond the scope they were created in
* Supports dynamic data structures (size can change)
* Good for large dat athat needs to persist

Cons of Heap
* slower allocation/deallocation than stack
* Potential memory leaks if not managed properly
* Can become fragmented over time
* More overhead for memory management

## Understanding Compiler vs Runtime

### Compile Time
This is the phase when the source code that someone wrote is analyzed and converted into an executable (or object code)
During this phase, the following usually happens:
* Syntax parsing: The compiler checks for grammar/structure
* Semantic checks: The compiler checks for type correctness, variable scoping, and other language rules
* Optimizations: The compiler may optimize code (in Go, it performs "escape analysis" to decide if something goes on the stack or heap)

After this, typically an executable program or a compiler artifact is generated to be run

### Runtime
This is the pahse when the compiler program actually executes
During this phase, the following happens:
* Memory allocation: the program rquests, uses, and frees memory on the fly
* Control flow: Branches, loops, concurrency, goroutine scheduling, etc
* Garbage collection: In Go, the garbage collection runs during runtime and automatically frees unused memory

This is basically when the real time execution of your code happends, handling all the inputs and products outputs.

## Scope vs Lifetime
Scope is static, based on the code at compile time because it can be determined by the compiler at compile time from the `{}` braces (i think)
Lifetime refers to the actual lifetime of the variables and it depends program execution time (runtime) because this is when variables are put into stacks or heaps 
when it is actually running
```Go
package xyz

func doIt() *int {
    var b int
    ...

    return &b
}
```
Variable `b` can only be seen inside doIt function, but its value will live on

Typically in languages like C/C++, you wouldnt be able to do this because the variable b would be allocated in the stack memory and when function doIt terminates,
the memory allocated in Stack for b would be released and you would lose the b value. Once losing the value from the stack, you wouldn't be able to use its memory address
to reference an empty space of stack that was already released. However, in Go, Go just takes that variable and puts it heap if you decide to return its memory address 
like in the example above. Initially b lived in stack, then it would be put into heap memory block to live as long as it is referenced by something. 


The value (object) will live on as long as part of the program keeps a pointer to it

## What is a closure?
A closure is when a function inside another function "closes over" one or more
local variables of the outer function (explanation provided below)
```Go
func fib() func() int {
    a, b := 0, 1

    return func() int {
        a, b = b, a + b
        return b
    }
}
```

The inner function gets a reference to the outer function's vars

Those variables may end up with a much longer lifetime than expected - as long as there's a reference to the inner function. 

If you were to do `f := fib`, f would just be a descriptor that has a pointer to function fib and doing `f()` 
would just call the function fib. However, if you were to `f := fib()`, f would NOT be the value of returned b variable,
but instead a descriptor too that has a pointer to the anonymous function AND a pointer to the original variables, thats what
a closure is and f would be a closure in the case of `f := fib()` (closure has both the function and the referenced variables
in the anonymous function, but this is just one specific type of closure!!!). Basically, f closure would have both 
the address of anonymous function and variables a, b. As soon as fib returns the anonymous function, a and b are no longer
inside stack memory but get moved to heap. Then the closure points  to heap to persist the lifetime of a and b. 
With every call now, `f()`, the anonymous function executes and keeps updating a and b all the while the states of a and b
exist in heap and we can keep getting next sequences of fib sequence by called `f()`. This way we don't have to declare global
variables for a and b, they were just created in the function and we were able to extend their life through closures.

## Point of Closures Explained with Sort.Slice
Example code:
```Go
sort.Slice(ss, func(i, j int) bool {
    return ss[i].val > ss[j].val
})
```

### The Probelm being solved
Go needs a way to sort any type of slice. But how can the sort package know how to compare custom structs or complex types?
It can't possibly know in advance what fields you want to sort by.

### Clearing up my misconceptions on Closures in Go
The documentation above only discusses one type of closure. Below is the more general definition of closures and practical use of it.

What are Closures in Go?

A closure is a function that can access variables from its containing scope, even after that scope would normally be gone (stack -> heap scenarios). In Go,
closures happen whenever an anonymous function accesses variables defined outside its body.

1. Different forms of closures:
   * Closures aren't just functions that return other functions (though that's one common pattern)
   * A closure can be any function that accesses variables from its surrounding scope
   * In the sort example, the anonymous comparison function is a closure because it accesses `ss` from the outer scope 
     and the state of ss is persisted in the anonymous function when it gets sent to `sort.Slice` method's scope
2. Closures don't always invole heap allocation:
   * Sometimes, closures may involve variables being moved to the heap from the stack in some cases where the containing scope dies (like function ending)
   * This implementation detail isn't essential to understanding what makes something a closure (like discussed above in the documentation)
   * A closure is generally defined by its behavior of accessing outer scope variables, not by how its implemented

### Anonymous Functions and Closures
To me, it seems like anonymous functions exist as a way to "inject logic" into other functions. Typically in functions, you can pass in variable parameters but you
can't change its logic inside that function. So Go, has anonymous functions as a way to inject specific code logic that would change some behavior of the function. 
Ofcourse what kind of logic that can be injected is designed by the designer of the funciton by providing a specific function signature in the params. For example,
in the `sort.Slice`, the method slice doesn't know how to sort certain types like structs. So that logic part is left to the user and we can tell Slice to sort based
on logic `ss[i].val > ss[j].val`. In other words, `sort.Slice` method demonstrates a powerful programming pattern where one function accepts another function as a 
parameter to customize its behavior (or logic). Below is a more organized and summarized version of what is meant here.

1. Separation of conernts:
   * The `sort.Slice` function knows HOW to sort (the algorithm)
   * You provide the logic for WHICH element comes before another
2. Anonymous function enable "logic injection":
   * You're literally inserting a piece of logic (the comparison) into another function
   * This allows cusotmizations without complex inheritance or interface implementation
   * It's like giving `sort.Slice` instructions on how to compare your specifc data type


Closures are perfect for Logic Injection and make them more powerful because:
1. They carry context:
   * The closure can "see" and use variables from its surrounding scope
   * In the `sort.Slice` example, it uses `ss` without needing it to be passed as a parameter (Slice method knows about `ss` because of closure, even tho it wasnt passed as param) 
   * This creates more concise, readable code
2. They enable on-the-fly customization:
   * You don't need to create a named function or a new type
   * The logic can be defined right where it's needed
   * This makes the code's intent clearer by keeping related pieces together
3. They reduce boilerplat:
   * Without closures, you'd need to create custom types with methods
   * Compare the closure approach with implementing a custom sort interface

```Go
// Without closures - requires new type definition
// the user would literally have to implement all these before using the sort function
type ByValue []kv
func (s ByValue) Len() int           { return len(s) }
func (s ByValue) Swap(i, j int)      { s[i], s[j] = s[i], s[j] }
func (s ByValue) Less(i, j int) bool { return s[i].val > s[j].val }

// Usage
sort.Sort(ByValue(ss))

// With closures - more concise and flexible
// the user just has to give logic with anonymous function
sort.Slice(ss, func(i, j int) bool {
    return ss[i].val > ss[j].val
})
```

This pattern (higher-order functions with closures) is powerful because:
1. It enables generic algorithms:
   * The same sorting code works for ANY type of slice
   * No need for seperate functions for different data types
2. It supports flexible customization:
   * Sort by different fields without rewriting the algorithm
   * Change sort direction (ascending/descending) with a simple change
3. It creates more readable, maintable code:
   * Intent is clear at the point of use
   * No need for separate type definitions just for sorting
4. It's more concise:
   * Fewer lines of code
   * Less cognitive overhead


## Common Mistakes with Closures
An importing thing to clarify:
```Go
func main() {
    name := "Alice"
    
    // This anonymous function is a closure because it accesses 'name'
    // from outside its own scope
    printGreeting := func() {
        fmt.Println("Hello,", name)
    }
    
    name = "Bob"  // Change the value
    
    printGreeting()  // Prints "Hello, Bob" not "Hello, Alice"
}
```
The closure captures a reference to the variable, not just its value at a point in time. 

Below is an another example that could be considered a bug in a code:
```Go
package main

import (
    "fmt"
)

func main() {
    s := make([]func(), 4)

    for i := 0; i < 4; i++ {
        s[i] = func() {
            fmt.Printf("%d @ %p\n", i, &i)
        }
    }

    for i := 0; i < 4; i++ {
        s[i]()
    }
}
// Output
// 4 @ 0x1400000e0a8
// 4 @ 0x1400000e0a8
// 4 @ 0x1400000e0a8
// 4 @ 0x1400000e0a8
``` 
Before Go 1.22, the i in the for loop would always be the same variable, so whenever the anonymous function captured the reference
(its mem address of i), it would always point to the same location in memory and at the end all functions in slice end up
referencing the same variable. To fix this, closure capture method was used.

```Go
package main

import (
	"fmt"
)

func main() {
	s := make([]func(), 4)

	for i := 0; i < 4; i++ {
		i2 := i // closure capture
		s[i] = func() {
			fmt.Printf("%d @ %p\n", i2, &i2)
		}
	}

	for i := 0; i < 4; i++ {
		s[i]()
	}
}
// Output
// 0 @ 0x1400000e0a8
// 1 @ 0x1400000e0b0
// 2 @ 0x1400000e0b8
// 3 @ 0x1400000e0c0
```
Now with each iteration, i2 is declared and a new disctinct mem address is given and the function captures the address
at that point in time (NOT THE VALUE, just the reference). This would give an output where each func in the slice points to
different values located in different memory locations. But now, after Go version 1.22, each iteration in for loop creates a 
separate variable with separate memory addresses (no longer referencing the same i varaible, but new one each time). Because 
of this, closure captures is NO LONGER NEEDED in the context of for loops. 

The most common Gotcha moment in closures is that the functions capture reference and not the value. That means in asynchronous
moments (when functions get executed later on), the referenced variable could have changed and can lead to issues. 