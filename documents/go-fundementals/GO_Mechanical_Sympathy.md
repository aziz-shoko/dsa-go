# Mechanical Sympathy
[Youtube](https://www.youtube.com/watch?v=7QLoOd9HinY&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=34)

"The most amazing achievement of the computer software industry is its
continuing cancellation of the steady and staggering gains made by the
computer hardware industry." - Henry Peteroski

We got similar perceived performance 30 years ago with:
* 100 times less CPU
* 100 times less memory
* 100 times less disk space

Mechanical sympathy is a term that means understanding how the underlying
hardware and systems work so you can write software that works with them
and not against them

## Performance in the cloud
But today, we've made a deliberate choice to accept some overhead

We have to trade off performance against other things:
* choice of architecture
* quality, reliability, scalability
* cost of development & ownership

We need to optimize where we can, given those choices
We still want simplicity, readability, & maintability of code

## Optimization
When optimizing, we should think top down. 

Top-down refinement:
Architecture - latency, cost of communication
Design - algorithms, concurrency, layers
Implementation - programming language, memory use

Mechanical sympathy plays a role in our implementation
Interpreted languages may cost 10x more to operate due to their inefficiency

Some unfortunate realities:
* CPUs aren't getting faster any more
* the gap between memory and CPU isn't shrinking 
* software gets slower more quickly than CPUs get faster

Software development costs exceed hardware costs

Go (and the Go philosophy) encourages good design:
you can choose
* to allocate contiguously
* to copy or not copy
* to allocate on the stack or heap (sometimes)
* to be synchronous or asynchronous
* to avoid unnecessary abstraction layers
* to avoid short/forwarding methods

Go doesn't get between you and the machine
Good code in Go doesn't hide the costs involved

## Memory hierarchy
As memory capacity increases, access latency also increases
cpu core -> cache -> memory -> ssd -> cloud

"Computational" cost is often dominated by memory access cost

Caching takes advantage of access patterns to keep frequently used code 
and data "close" to the CPU to reduce access time

Caching imposes some costs of its own
* Memory access by the cache line, typically 64 bytes
* Cache coherency to manage cache line ownership


# Becnhmarking
Go Benchmarks: [Youtube](https://www.youtube.com/watch?v=nk4rALKLQkc&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=35)
Benchmarks live in test files ending with _test.go
You run benchmarks with go test -bench
Go only runs the BenchmarkXXX functions

## How Benchmarking works in Go
Go's benchmarking framework is built into the `testing` package. Here's how it works:

Time Control: The framework controls when timing starts and stops, allowing for setup code that doesn't affect measurements
Adaptive iteration: Go automatically determines how many times to run your code to get stable measurements. It starts with 1 
iteration and increases until it can measure reliably.
Statistical Analysis: Multiple runs provide statistical confidence in results
Benchmarks run alignside tests using the same `go test` command

## Basic Syntax
```Go
func BenchmarkFunctionName(b *testing.B) {
    // Setup code (not timed)
    
    b.ResetTimer() // Optional: reset timer after setup
    
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        functionToTest()
    }
}
```
* Function name must start with `Benchmark`
* Must accept `*testing.B` parameter
* Must contain a loop that runs `b.N` times
* Place the code you want to measure inside the for loop

## Basic commands
```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkMyFunction

# Run benchmarks matching pattern
go test -bench=BenchmarkString.*

# Run benchmarks multiple times for better accuracy
go test -bench=. -count=5

# Set minimum benchmark time
go test -bench=. -benchtime=10s

# Set specific number of iterations
go test -bench=. -benchtime=1000000x
```

## Memory Profiling
```bash
# Include memory allocation stats
go test -bench=. -benchmem

# Generate memory profile
go test -bench=. -memprofile=mem.prof

# Generate CPU profile
go test -bench=. -cpuprofile=cpu.prof
```

## Understanding Benchmark Output
```
BenchmarkStringConcat-8    1000000    1234 ns/op    48 B/op    2 allocs/op
```
* `BenchmarkStringConcat`: Function name
* `-8`: Number of CPU cores (GOMAXPROCS)
* `1000000`: Number of iterations (b.N)
* `1234 ns/op`: Nanoseconds per operation
* `48 B/op`: Bytes allocated per operation (with -benchmem)
* `2 allocs/op`: Number of allocations per operation (with -benchmem)


# Profiling
For profiling instructions, go to:
[Youtube](https://www.youtube.com/watch?v=MDB2x1Di5uM&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=36)


# Static Analysis
Static Analysis, or also known as linting, are tools that are used to make the code better before 
actually running the code. 

"Static" means the program isn't running ("compile time")

Static analsysis transfers effor from people to tools
* mental effort while coding 
* code review effort

Static analysis improves code hygiene 
* correctness
* efficiency
* readability
* maintainability

If our code compiles & passes static analysis, we can have a lot of confidence in it even before
running unit tests.

I run these tools in my IDE every time i save a file:
* format the code
* fix the imports
* looks for issues

`go fmt` will put your code in standard form (spacing, indentation)
`go imports` will do that and also update import lists
`go lint` will check for non-format style issues, for example:
* exported names should have comments for `godoc`
* names shouldn't have `under_scores` or be in `ALLCAPS`
* `panic` shouldn't be used for normal error handling
* the error flow should be indented, the happy path not
* variable declarations shouldn't have redundant type info

The "rules" are based on Effective Go and Google's Go Code Review Comments

`go vet` will find some issues the compiler won't
* suspicious "printf" format strings
* accidentally copying a mutex type
* possibly invalid integer shifts
* possibly invalid atomic assignments
* possibly invalid struct tags
* unreachable code

`go cyclo` reports high cyclomatic complexity in functions 
(general rule is to only break up functions if cyclomatic complexity gets too high)

No static analysis tool can find all possible errors
And there are many other go commands

## One tool to rule them all
All the tools can be runned using `golangci-lint` package
It can be configured with `.golangci.yml`

It can be used in the CI/CD pipeline
Issues must be fixed for the build to pass

False positivescan be marked with //nolint

## Same VSC settings
```json
{
    "go.vetOnSave": "package",
    "go.formatTool": "goimports",
    "go.formatFlags": [
        "-local github.com/xxx,github.com/yyy"
    ],
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ],
    "go.lintOnSave": "package"
}