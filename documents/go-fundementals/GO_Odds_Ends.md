# Enumerated Types

Docs based on this video:
[Youtube](https://www.youtube.com/watch?v=oTtYtrFv3gw&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=31)

There no real enumerated types in Go

You can make an almost-enum type using a named type and constants:
```Go
type shoe int

const (
    tennis shoe = iota
    dress
    sandal
    clog
)
```

`iota` starts at 0 in each const block and increments on each line;
here 0,1,2

# Variable argument list
You can have the function accept as many params as it is given.
```Go
package main

import(
    "fmt"
)

// ONLY THE LAST PARAM CAN BE VARIABLE PARAM
func sum(nums ...int) {     // nums is treated as a slice
    var total int

    for _, num := range nums {
        total += num
    }

    return total
}

func main() {
    fmt.Println(sum())
    fmt.Println(sum(1))
    // fmt.Println(sum(1, 2, 3, 4))

    s := []int{1 ,2, 3, 4}
    fmt.Println(sum(s...))  // the variable... syntax unpacks the slice s into individual arguments for function sum
}
```

Basically the rest of the video talks about BitWise operations and different int types. For more info, go watch video