# Pointers vs Values
Pointers are shared values, but not copied
Values are copied values, but not shared

Value semantics lead to higher integrity, particularly with concurrency
(don't share)

Pointer semantics may be more efficient

## Common uses of pointers
* some objects can't be copied safely (mutex)
* some objects are too large to copy efficiently (consider pointers when size > 64 bytes)
* some methods need to change (mutate) the receiver (as in changes in method must reflect object situtations)
* when decoding protocal data into an object (JSON, etc.; often in a variable argument list)
```Go
var r Response
err := json.Unmarshal(j, &r)
```
* when using a pointer to signal a "null" object

Any small struct under 64 bytes probably should be copied. You simply pass in struct, func copies it and returns changed struct.
```Go
type Widget struct {
    ID int
    Count int
}

func Expend(w Widget) Widget {
    w.Count--

    return w
}
```
Note that Go routinely copies string & slice descriptors. Descriptors are basically structs that hold data like the pointer 
to underlying string or array, and any additional info like length of string or capacity for slice etc. These descriptor structs
are less than 64 bytes and are routinely copied in Go.

## Stack Allocation
Stack allocation is more efficient.

Accessing a variable directly is more efficient than following a pointer so Go tries to allocated stuff in the stack
and avoid pointers. 

Accessing a dense sequence of dta is more efficient than sparse data (an array is faster than a linked list, etc.)

Older languages also put things in the stack because the stack is faster and more efficient for stuff like garbage collection, etc.
Newer interprated languages typically allocated stuff in Heap at great efficiency expense. 

## Heap Allocation
Go would prefer to allocate on the stack, but sometimes can't:
* a function returns a pointer to a local object
* a local object is captured in a function closure
* a pointer to a local object is sent via a channel
* any object is assigned into an interface
* any object whole size is a variable at runtime (slices)

The use of `new()` has nothing to do whether something gets allocated in stack or heap, and therefore generally not used.

Build with the flag `-gcflags` `-m=2` to see the escape analysis (but some things are too complicated to escape analysis)
Escape analysis is to see which things went to stack or heap. 

## For loops
The value returned by range is always a copy.
```Go
for i, thing := range things {
    // thing is a copy
    ...
}
```
So if we were iterating over a slice of structs, thing would be a copy of the struct and things.SomeAttribute
reassigned would only reflect the change within the for loop scope and the changes wouldn't reflect once for loop ends.

Use the index if you need to mutate the element:
```Go
for i := range things {
    things[i].which = whatever
    ...
}
```
By using the index, you are directly accessing the struct itself and not getting the copy of it. So now you can make changes
to reflect on the original struct.

## Slice safety
Anytime a function mutates a slice that's passed in, we must return a copy. As in anytime slice size grows, and if it doesn't,
we don't really have to worry about returning it. The reason why we have to return it is because with bigger size means new 
memory address and this new memory address needs to be reassigned back to caller or else caller slice will still keep pointing
to original array. 
```Go
func update(things []thing) []thing {
    ...
    things = append(things, x)
    return things
}
```
Thats because the slice's backaing array may be reallocated to grow.
```Go
// another example
package main

import "fmt"

func main() {
    // Case 1: Caller's slice gets updated (no reallocation)
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Println("Before modifyElements:", numbers)
    modifyElements(numbers)
    fmt.Println("After modifyElements:", numbers) // Changes are visible
    
    // Case 2: Caller's slice does NOT get updated (reallocation happens)
    numbers = []int{1, 2, 3, 4, 5}
    fmt.Println("\nBefore addElementWrong:", numbers)
    addElementWrong(numbers, 6)
    fmt.Println("After addElementWrong:", numbers) // No change visible
    
    // Case 3: Caller's slice gets updated correctly (with return value)
    numbers = []int{1, 2, 3, 4, 5}
    fmt.Println("\nBefore addElementCorrect:", numbers)
    numbers = addElementCorrect(numbers, 6)
    fmt.Println("After addElementCorrect:", numbers) // Change is visible
}

// Case 1: Simple modification - no reallocation
func modifyElements(s []int) {
    for i := range s {
        s[i] *= 2
    }
    // No return needed since we're just modifying existing elements
}

// Case 2: Append operation but NOT returning the new slice
func addElementWrong(s []int, x int) {
    s = append(s, x) // Creates new array but result is never returned
    fmt.Println("Inside addElementWrong:", s) // [1 2 3 4 5 6]
    // Not returning s means caller still has the original slice
}

// Case 3: Append operation WITH proper return
func addElementCorrect(s []int, x int) []int {
    return append(s, x) // Returns new slice header pointing to new array
}
```

Another thing to keep in mind is that keeping a pointer to an element of a slice is risky.
In the example below, users is a descriptor that holds the mem address of an array of structs. 
alice is then a variable that holds the address for first element. But if slice get resized, it can 
completely change its memory address where itself and its elements live and alice might point to a location
that is no longer there. 
```Go
type user struct {
    name string
    count int
}

func addTo(u *user) {
    u.count++
}

func main() {
    users := []user{{"alice", 0}, {"bob", 0}}
    alice := &users[0]  // risky
    amy := user{"amy", 1}

    users = append(users, amy)

    addTo(alice)        // alice is like a stale pointer
    fmt.Println(users)  // so alice's count will be 0
}
```

## Capturing the loop variable
Be careful with using memory address in Go's for range loop

```Go

func (r OfferResolver) Changes() []ChangeResolver {
	var result []ChangeResolver

	for _, c := range r.d.Status.Changes {
		change := c // make unique

		result = append(result, ChangeResolver{&change})
	}
}
```

In this case, you would still need to keep the change := c // make unique line.
The optimization in modern Go only applies to the specific case of creating slices from arrays with expressions like item[:]. 
It doesn't change the fundamental behavior of how loop variables work in for range loops.
Your understanding is exactly right:

The loop variable c is still allocated at a specific memory location that stays constant throughout the loop
Only the value stored at that location changes with each iteration
If you directly used c in your append (without the local variable assignment),
you'd be creating ChangeResolver structs that all contain pointers to the same memory location
By loop's end, all these pointers would point to the last change in the list

The line change := c // make unique creates a new variable with a new memory address for each iteration. 
Each ChangeResolver will then contain a pointer to a different memory location, preserving the correct data.
Modern Go has not changed this behavior and likely never will, as it would be a fundamental change to how variables
and memory work in the language. The array-to-slice optimization was a special case that
could be implemented without breaking backward compatibility.
So to summarize: Yes, you definitely need to keep the change := c line in your code for it to work correctly,
even in the latest versions of Go.

This is the same for traditonal `for i := 0; i < n; i++ {}` loops where i has same address but changing values. But modern Go
handles some special cases to automatically. 