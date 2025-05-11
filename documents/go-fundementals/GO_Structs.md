# What are Stucts?
Structs are basically a data structure that holds multiple data types together. 

```Go
type Employee struct {
    Name string
    Number int
    Boss *Employee
    Hired time.Time
}

func main() {
    var e Employee

    e.Name

    // the + sign prints the Field names along with the value
    fmt.Printf("%T %+[1]v\n", e)
}
// Outputs 
// main.Employee {Name: Number:0 Boss:<nil> Hired:0001-01-01 00:00:00 +000 UTC}
```

In Go, be careful with how you create a hashMap of structs. Typically, hashMap value is a pointer to structs.
```Go
type Employee struct {
    Name string
    Number int
    Boss *Employee
    Hired time.Time
}

func main() {
    c := map[string]*Employee{}

    c["Lamine"] = &Employee{"Lamine", 2, nil, time.Now()}

    c["Lamine"].Number++

    c["Matt"] = &Employee{
        Name: "Matt",
        Number: 1, 
        Boss: c["Lamine"],
        Hired: time.Now(),
    }

    fmt.Printf("%T %+[1]v\n", c["Lamine"])
    fmt.Printf("%T %+[1]v\n", c["Matt"])
}

// func main() {
//     c := map[string]Employee{}

//     c["Lamine"] = Employee{"Lamine", 2, nil, time.Now()}

//     c["Lamine"].Number++

//     c["Matt"] = Employee{
//         Name: "Matt",
//         Number: 1, 
//         Boss: &c["Lamine"],
//         Hired: time.Now(),
//     }

//     fmt.Printf("%T %+[1]v\n", c["Lamine"])
//     fmt.Printf("%T %+[1]v\n", c["Matt"])
// }
```

The uncommented section would work as expected, but if here we can't have hashMap of struct values (only pointers)
because, for example, `&c["Lamine"]` the memory address can change in hashMap when dealing with collisions and resizing
or hashing. The memory address wouldn't stay constant and can change and corrupt data. But with storing pointers as values, 
the holders of the pointers can change in the hashMap all the want but the actual underlying structs will live in a constant
place in memory. 

## Struct Compatibility
A struct has two compatability issues that must be addressed if trying to assign one to another, its struct fields
and name compatability

```Go
type album1 struct {
    title string
}

type album2 struct {
    title string
}

func main() {
    var a1 = album1{
        "the white album"
    }

    var a2 = album2{
        "the black album"
    }

    // a1 = a2 wouldnt work because they have different names declaraed (name compatability)
    // the reason why we can assign them in the first place is because we first cast a2 to album1
    // and both structs have the same fields of title
    a1 = album1(a2)
    fmt.Println(a1, a2)
}
```

Two struct types are compatible if:
* the fields have the same types and names
* in the same order
* and the same tags (*)

A struct may be copied or passed as a parameter in its entirety
A struct is comparable if all its fields are comparable
The zero value for a struct is "zero" for each field in turn

## Struct Tags
Stuct tags are basically extra information in the fields that help other programs deal with Go stuff.
The most common example of this is Go structs working with JSON.
So its mostly used for data conversion from Go to external formats (JSON, XML, etc)

```Go
package main

import (
    "fmt"
)

type Response struct {
    // the json: is the struct tag part, notice its wrapped in backticks
    Page int    `json:page`
    Words []string `json:words,omitempty`
}

func main() {
    r := &Response(Page: 1, Words: []string{"up", "in", "out"})
}
```