# Concurrency

[Youtube Video](https://www.youtube.com/watch?v=A3R-4ZYBqvE&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=22)

Concurrency Definition:
Parts of the program may execute independently in some none deterministic order (sometimes partial) order

Parallelism:
Parts of a program execute independently at the same time

You can have concurrency with a single-core processor (as in interrupt handling in the operating system)

Parallism can happen only on a multi-core processor

Concurrency doesn't make the program faster, parallelism does

Concurrency Vs. Parallelism
Concurrency is about dealing with things happening out of order (like your program listens
for requests but rather than waiting, you can use concurrency for it to do something else
while its waiting for any requests)

Parallelism is about things actually happening at the same time

A single program won't have parallelism without concurrency

We need concurrency to allow parts of the program to execute independently
(And that's where the fun begins)

Race Conditions:
System behavior depends on the (non-deterministic) sequence or timing of parts of the program
executing independently, where some possible behaviors (orders of exeuction) produce invalid results


In other words, its a bug and should not work that way. 

Some ways to solve race conditions:
Race conditions involve independent parts of the program changing things that are shared

Solutions making sure operations produce a consistent state to any shared data:
* don't share anything
* make the shared things read-only
* allow only for one writer to the shared things
* make the read-modify-write operations atomic (as in non divisible, perform all the at the same time)

# Concurrency in Go
Sections below will be about Go channels and goroutines

## Channels 

A channel is a one way communcations pipe
* things go in one end, come out the other 
* in the same order they went in
* until the channel is closed
* multiple readers & writers can share it safely

### Sequential Process
Looking at a single independent part of the program, it appears to be sequential
```Go
for {
    read()
    process()
    write()
}
```
This is perfectly natural if we think of reading & writing files or network sockets

### Communicating sequential processes
Now put the parts together with channels to communicate
* each part is independent
* all they share are teh channels between them
* the parts can run in parallel as the hardware allows

Concurrency is always hard
(the human brain has hard time thinking that way)

CSP (Communicating sequential processes) provides a model for thinking about it that makes it less hard 
(take the program apart and make the pieces talk to each other)

Basically, to my understanding, CSP allows us to think synchronously where the sequential processes we create,
talk to each other through channels and these channels act as buffers or synchronization points between these sequential processes

"Go doesn't force developers to embrace the asynchronous ways of event-driven programming. ... That lets you write
asynchronous code in a synchronous style. As people, we're much better suited to writing things in a synchronous style."
- Andrew Gerrand

### Goroutines
A goroutine is a unit of independent execution (coroutine)

It's easy to start a goroutine: put go in front of a function call

The trick is knowing how the goroutine will stop:
* you have a well defined loop terminating condition, or
* you signal completion through a channel or context, or 
* you let it run until the program stops

But you need to make sure it doesn't get blocked by mistake

### Channels
Is a data type that represents a way of communicating

A channel is like a one-way socket or a Unix pipe
(except it allows multiple readers & writers)

It's a method of synchorinization as well as communication
We know that a send (write) always happens before a receive call (read)

It's also a vehicle for transferring ownership of data, so that only one goroutine at a time
is writing the data (avoid reace conditions)

"Don't communicate by sharing memory; instead, share memory by communicating." - Rob Pike

## Go Channels and Goroutines Notes

### Abstract idea
* Goroutines are lightweight, independent workers (like mini threads)
* Channels are typed pipes that goroutines use to talk to each other
* Instead of "shared memory by locks", Go encourages "share memory by communicating"

Think: Divide the problem into independent steps. Let each part do its own job, and let them talk only through channels

Go's model is based on CSP(Communicating Sequential Processes) -- you write seemingly sequetial code, but it runs concurrently with others

### Goroutines
A goroutine is a function that runs independently in the background, it doesn't blcok the current function

Syntax:
```Go
go someFunction()

// Example code
func printMessage() {
    fmt.Println("Hello from goroutine")
}

func main() {
    go printMessage()       // starts a new goroutine
    time.Sleep(time.Second) // wait so goroutine can finish
}

// NOTE: You wont see antying if main() exits before the goroutine runs -- goroutines run asynchronously
```

Behind the scenes:
* Goroutines are managed by Go's scheduler, not the OS
* They're cheap (a few KB stack) -- you can spawn thousands
* Go uses a work stealing scheduler to distribute goroutines across available CPU threads (`GOMAXPROCS`)

### Channels
A channel is a typed conduit for sending and receiving values between goroutines

Syntax:
```Go
ch := make(chan int)        // unbuffered channel
ch := make(chan string, 10) // buffered channel with size 10

// Sending and receiving 
ch <- 42    // send 42 into channel
x := <-ch   // receive from channel

// NOTE: Unbuffered channel forces synchronization between sender and receiver

// Basic Example
func worker(ch chan string) {
    msg := <-ch
    fmt.Println("Received:", msg)
}

func main() {
    ch := make(chan string)
    go worker(ch)
    ch <- "Hello Channel"
}
```
Whats happening above:
1. ch := make(chan string) creates an unbuffered channel
2. go worker(ch) launches a goroutine that runs concurrently with the main goroutine
3. The worker immediately tries to receive from the channel with msg := <-ch
4. Since the channel is unbuffered, the worker goroutine blocks until data is sent
5. The main goroutine continues and sends "Hello Channel" with ch <- "Hello Channel"
6. This unblocks the worker, which can now receive the message and print it

The key insight is that channels in Go provide synchronization between goroutines. With unbuffered channels:

A send operation blocks until another goroutine is ready to receive
A receive operation blocks until another goroutine sends data

Even with multiple cores (true parallelism), the program won't crash. The worker goroutine will still 
block on the receive operation until data is available, regardless of whether it's running on a separate core.

When you say "goroutines are asynchronous by nature," it means they run independently of the main program flow, 
but they still synchronize with each other when they communicate through channels.

#### Buffered vs Unbuffered channels
Unbuffered by default
```Go
ch := make(chan int)
```
* Sender blcoks until receiver is ready
* Forces tight synchronization
* Good for coordinating actions

Buffered
```Go
ch := make(chan int, 3)
```
* Sender only blocks if buffer is full 
* Receiver blocks if buffer is empty
* Allow decoupling producer/consumer speeds

### Closing and ranging over channels
```Go
close(ch)
```
* Only senders should close.
* Can't send after closing. But you can receive from closed channels.

```Go
for val := range ch {
    fmt.Println(val)
}
```
Stops when channel is closed and all values are received.

### Select Statement -- Waiting on Multiple Channels
```Go
select {
case msg1 := <-ch1:
    fmt.Println("Received", msg1)
case msg2 := <-ch2:
    fmt.Println("Received", msg2)
default:
    fmt.Println("No channel ready")
}
```
* Think of it like switch, but for channels.
* Useful for multiplexing, timeouts, or non-blocking communication.

### Gotchas and Pitfalls
1. Deadlocks
* When goroutines wait forever on each other.
* Often caused by missing receiver/sender or forgetting to close channels.

2. Closing channels too early
* Only the sender should close; closing from the wrong place causes panics.

3. Race conditions
* Goroutines don't magically solve state-sharing issues — avoid shared memory or use sync primitives (sync.Mutex).

4. Leaky goroutines
* If a goroutine waits forever on a channel that no one sends to, it stays alive. This creates memory leaks.

### Intuition
1. Break Problems Into Steps
"What are the independent units of work in this problem?"

Each one becomes a goroutine.

2. Design the Flow of Data, Not Control
"Who produces data? Who consumes it?"

Channels connect producers and consumers. Don't pass control, pass data.

3. Let the System Coordinate Itself
Don't micromanage goroutines. Let channels determine timing.

If a goroutine needs data, it'll block. When it's ready, the data comes through.

4. Think of Channels as Conversation Lines
A goroutine waits on the phone until someone picks up (synchronization).

With a buffered channel, it's like leaving a voicemail.

### Patterns and Use cases
Fan Out/Fan In
* Fan-Out: Multiple goroutines read from same input channel.
* Fan-In: Multiple goroutines send to one channel.
```Go
out := make(chan int)
for i := 0; i < 3; i++ {
    go worker(out)
}
```

Worker Pools
* Goroutines waiting for jobs from a jobs channel.

Pipelines
* Chain of goroutines passing data via channels.
```
stage1 → stage2 → stage3
```

Timeouts
```Go
select {
case <-time.After(time.Second * 2):
    fmt.Println("Timeout!")
}
```

### Summary Cheatsheet
Here's a markdown table summarizing the key Go concurrency concepts:

| Concept          | Keyword/Syntax        | Purpose                        |
| ---------------- | --------------------- | ------------------------------ |
| Goroutine        | `go fn()`             | Lightweight async worker       |
| Channel          | `make(chan T)`, `<-`  | Communication pipe             |
| Buffered Channel | `make(chan T, size)`  | Decouples sender/receiver      |
| Closing Channel  | `close(ch)`           | Signals no more sends          |
| Select           | `select { case ... }` | Waits on multiple channels     |
| Range Channel    | `for x := range ch`   | Receive till channel is closed |

### Philosophy
"Model your program like a conversation. Who talks to whom? What do they say? Let them speak only through channels."

Avoid shared state — it's messy and leads to bugs.
Concurrency is not about speed, it's about structure and clarity.
Design around data flow, not just function flow.

Experienced Go developers follow these principles:

* Design for independence: Make goroutines as independent as possible
* Communicate with channels: Use channels for synchronization and data exchange
* Don't rely on timing: Never assume one goroutine will run before another


### Select 
`select` allows any "ready" alternative to proceed among
* a channel we can read from
* a channel we can write to
* a default action that's always ready

Most often `select` runs in a loop so we keep trying

We can put a timeout or "done" channel into the select
* We can compose channels as synchronization primitives!
* Traditional primitives (mutex, condition variable) can't be composed

```Go
package main

import (
	"time"
	"log"
)

func main() {
	log.Println("start")

	const tickRate = 2 * time.Second

	stopper := time.After(5 * tickRate)
	ticker := time.NewTicker(tickRate).C

loop:
	for {
		select {
		case <- ticker:
			log.Println("tick")
		case <- stopper:
			break loop
		}
	}
	log.Println("finish")
}
```

In the example above, you can use select statement to simultaneously to listen to different channels
(the select statement executes whichever channel it receives from right away) and have a stopper channel
to break out of the loop.

Select also has a default case. In a select block, the default case is always ready and will be chosen if no other case is:
```Go
func sendOrDrop(data []byte) {
    select {
    case ch <- data:
        // sent ok; do nothing
    default:
        log.Printf("overflow: drop %d bytes", len(data))
    }
}
```
Don't use default inside a loop, the select will busy wait and waste CPU
(once i get good at concurrency, a recommended read is: concurrency in go by katherine cox-buday)

### Fan-In strategy example
```Go
package main

import (
	"math/rand"
	"time"
	"fmt"
)

func main() {
	// Create channels for each service
	authLogs := make(chan string)
	orderLogs := make(chan string)
	paymentsLogs := make(chan string)
	
	// Start services that continously generate logs
	go generateAuthLogs(authLogs)
	go generateOrderLogs(orderLogs)
	go generatePaymentLogs(paymentsLogs)

	// Fan-in to merge all logs into a single channel
	mergedLogs := fanIn(authLogs, orderLogs, paymentsLogs)

	// Consume from the merged channel
	for log := range mergedLogs {
		fmt.Println("LOG:", log)
	}
}

func generateAuthLogs(out chan<- string) {
	for {
		time.Sleep(time.Duration(rand.Intn(3))* time.Second)
		out <- "auth: user login"
	}
}

func generateOrderLogs(out chan<- string) {
	for {
		time.Sleep(time.Duration(rand.Intn(3))* time.Second)
		out <- "order: new order placed"
	}
}

func generatePaymentLogs(out chan<- string) {
	for {
		time.Sleep(time.Duration(rand.Intn(3))* time.Second)
		out <- "payment: transaction successful"
	}
}

func fanIn(ch1, ch2, ch3 <-chan string) <-chan string {
	merged := make(chan string)

	go func() {
		for {
			select{
			case log := <- ch1:
				merged <- log
			case log := <- ch2:
				merged <-log
			case log := <- ch3:
				merged <- log
			}
		}
	}()
	
	return merged 
}
```
Characteristics of Fan-In Pattern

Definition: Fan-in consolidates data from multiple sources into a single channel
Key Components:

Multiple input channels (producers)
* A single output channel (consumer)
* A multiplexer function that forwards messages from inputs to output

Benefits:

* Simplifies consumer code (reads from just one channel)
* Decouples producers from consumers
* Allows independent, concurrent producers to feed into a single processing pipeline

Implementation Details:

* Uses select to wait on multiple channels simultaneously
* Runs in its own goroutine to avoid blocking
* Can handle any number of input channels (though reflect.Select is needed for dynamic channel counts)

Common Use Cases:

* Log aggregation from multiple services
* Combining results from parallel workers
* Merging events from different sources
* Implementing publish-subscribe patterns

Go-Specific Features:

* Natural implementation using Go's channel and select mechanism
* Can be made generic to handle any number and type of channels
* Often paired with the fan-out pattern for parallel processing

Lifecycle Management:

* Must handle channel closing properly
* Can include timeout or cancellation mechanisms
* Typical implementations run indefinitely until input channels close or context cancels

### Time-Out example 
```Go
 
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	services := []string{"auth", "payments", "inventory"}
	results := make(chan string, 2)

	for _, s := range services {
		go func(s string) {
			resultCh := make(chan string)

			// Launch the check itself
			go checkService(s, resultCh)

			select {
			case res := <-resultCh:
				results <- res
			case <-time.After(2 * time.Second):
				results <- fmt.Sprintf("%s: timed out", s)
			}
		}(s)
	}

	for msg := range results {
		fmt.Println(msg)
	}

}

func checkService(service string, resultCh chan<- string) {
	delay := time.Duration(rand.Intn(4)) * time.Second
	time.Sleep(delay)

	resultCh <- fmt.Sprintf("%s: healthy (responded in %v)", service, delay)
}
```
The code above is asimple timeout strategy example, PROBABLY NOT BEST PRACTICE!
Just done for learning purposes

### Pipeline strategy example
```Go
package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	logged := make(chan string)
	cleaned := make(chan string)
	filtered := make(chan string)
	results := make(chan []string)

	go GenerateLogs(logged)
	go CleanLogs(logged, cleaned)
	go FilterLogs(cleaned, filtered)

	go func() {
		results <- StoreLogs(filtered)
	}()

	logs := <- results
	fmt.Println("Stored Logs:", logs)
}

func GenerateLogs(out chan<- string) {
	logs := []string{
		"  ERROR Disk full  ",
		"INFO system rebooted",
		"Warning: low battery",
		" error unable to write file ",
		"Something went wrong",
	}

	for _, log := range logs {
		time.Sleep(300 * time.Microsecond)
		out <- log
	}
	close(out)
}

func CleanLogs(in <-chan string, out chan<- string) {
	for rawLog := range in {
		trimmedLog := strings.TrimSpace(rawLog)
		cleanLog := strings.ToLower(trimmedLog)
		out <- cleanLog
	}
	close(out)
}

func FilterLogs(in <-chan string, out chan<- string) {
	for log := range in {
		if strings.Contains(log, "error") {
			out <- log
		}
	}	
	close(out)
}

func StoreLogs(in <-chan string) []string {
	var logs []string
	for log := range in {
		logs = append(logs, log)
	}
	return logs
}
```
Pipeline Strategy Characteristics
Definition: A pipeline processes data through a series of stages, where each stage receives values from the previous stage, performs some transformation, and sends results to the next stage.

Key Components:
* Input source (producer/generator)
* One or more intermediate processing stages
* Output sink (consumer/collector)
* Channels connecting each stage

Benefits:
* Enables concurrent processing of different data items at different stages
* Creates clear separation of concerns between processing steps
* Allows for efficient resource utilization
* Simplifies complex data transformations into discrete, manageable steps

Implementation Details:
* Each stage runs in its own goroutine
* Stages are connected via channels
* Each stage ranges over its input channel and produces values to its output channel
* Channels are typically closed by the sender when no more data will be sent

Common Use Cases:
* Data processing workflows
* ETL (Extract, Transform, Load) operations
* Stream processing
* Log processing (as in your example)

Go-Specific Features:
* Channel directionality (chan<- for send-only, <-chan for receive-only)
* Range over channels to process all values until closure
* Explicit channel closing to signal completion

Lifecycle Management:
* Upstream stages close outbound channels when done
* Downstream stages detect completion through channel closure
* Pipeline naturally terminates when all data is processed