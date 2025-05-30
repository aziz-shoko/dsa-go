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

### Fan-Out code example
```Go
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func main() {
	tasks := make(chan string, 10)
	results := make(chan string)

	// start worker pool
	for range 3 {
		go worker(tasks, results)
	}

	// Generate tasks in the background
	go generateTasks(tasks)

	// Collect and display results
	expectedResults := 6 // 6 images
	for range expectedResults {
		img := <-results
		fmt.Println(img)
	}

}

// generateTasks simulates creating a batch of image processing tasks
func generateTasks(out chan<- string) {
	images := []string{
		"vacation.jpg",
		"family.jpg",
		"dog.jpg",
		"graduation.jpg",
		"wedding.jpg",
		"party.jpg",
	}

	for _, img := range images {
		out <- img
	}
	close(out)
}

// processImage simulates image processing (resize, filter, etc.)
func processImage(img string) string {
	fmt.Printf("Processing image: %s\n", img)

	// Simulate processing time (random duration)
	processingTime := time.Duration(1+rand.Intn(3)) * time.Second
	time.Sleep(processingTime)

	// return processed image (just a string transformation in this case)
	return strings.ToUpper(img) + " [PROCESSED]"
}

// Implement a worker function that takes tasks from a channel
// and sends results to another channel
func worker(in <-chan string, out chan<- string) {
	for task := range in {
		processedImage := processImage(task)
		out <- processedImage
	}
}
```
The code above is an example of Fan-Out strategy because from one channel of tasks, we are spawning
out multiple goroutines to handle tasks

Fan-Out Pattern in Go - Key Characteristics
Definition: Fan-out distributes work across multiple goroutines to process data concurrently, allowing parallel execution of similar tasks.
Key Components:
* A single source channel of tasks/data
* Multiple worker goroutines that all consume from this channel
* Each worker performs the same operation on different data items
* (Optional) A results channel for collecting processed data

Benefits:
* Parallelizes CPU-bound operations for better performance
* Distributes workload evenly across available resources
* Speeds up batch processing operations
* Provides natural load balancing (faster workers process more items)


Implementation Details:
* Workers compete for tasks from the shared input channel
* Channel ensures each task is processed exactly once
* Workers typically run identical code but on different data
* Each worker runs until input channel is closed and drained


Common Use Cases:
* Batch processing (images, documents, etc.)
* CPU-intensive calculations on multiple items
* Parallel API requests
* Data transformations on large datasets


Go-Specific Features:
* Uses goroutines for lightweight concurrent execution
* Works well with Go's channel mechanics
* Naturally balances work across available cores
* Often paired with fan-in for collecting results


Lifecycle Management:
* Workers exit when input channel closes
* Results typically collected by counting expected items or using sync.WaitGroup

# Channels in Detail
[Youtube](https://www.youtube.com/watch?v=fCkxKGd6CVQ&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=26)

## Channel state
Channels block unless ready to read or write

A channel is ready to write if
* it has buffer space, or 
* at least one ready is ready to read (rendezvous)

A channel is ready to read if
* it has unread data in its buffer, or
* at least one writer is ready to write (rendezvous), or
* it is closed

Channels are unidirectional, but have two ends
(which can be passed separately as parameters)
```Go
// An end for writing & closing
func get(url string, ch chan<- result) {		// write-only end
	...
}

// An end for reading
func collect(ch <-chan result) map[string]int {	// read-only end
	...
}
```

### Closed Channels
Closing a channel causes it to return the "zero" value
(Unless it is a buffered channel with many values, then all the values in channel 
will be processed and the default 0 value be sent at the end after the channel closes)
We can receive a second value: is the channel closed?
```Go
// Buffered channels can hold values even if no one is receiving them immediately. 
// Unbuffered channels require a receiver to be ready at the same time as the sender,
// or the program blocks (deadlocks if nobody ever reads).
func main() {
	ch := make(chan int, 1)	// NOTICE: It is buffered channel
	ch <- 1					// If it was unbuffered channel, then would be blocked here

	b, ok := <-ch
	fmt.Println(b, ok)		// 1 true

	close(ch)

	c, ok := <-ch
	fmt.Println(c, ok)		// 0 false
}
```

A channel can only be closed once (else it will panic)
But why close a channel? if the channel doesn't have a value, then it is
unreadable and gets blocked. By closing the channel, the default readable value
of 0 is returned, so closing it makes it readable and signals to the receiver that 
the channel is done and no more messages

One of the main issues in working with goroutines is ending them
* An unbuffered channel requires a reader and writer (a writer blocked on a channel with no reader will "leak")
* Closing a channel is often a signal that work is done
* Only one goroutine can close a channel (not many)
* We may beed some way to coordinate closing a channel or stopping goroutines (beyon the channel itself)


Nil Channels

Reading or writing a channel that is nil always blocks *

But a nil channel in a select block is ignored
This can be powerful tool:
* Use a channel to get input 
* Suspend it by changing the channel variable to nil
* You can even un-suspend it again
* But close the channel if there really is no more input (EOF)

### Channel State Reference
## Go Channel State Reference

| **State**        | **Receive**       | **Send**      | **Close**            |
| ---------------- | ----------------- | ------------- | -------------------- |
| **Nil**          | Block\*           | Block\*       | Panic                |
| **Empty**        | Block             | Write         | Close                |
| **Partly Full**  | Read              | Write         | Readable until empty |
| **Full**         | Read              | Block         |                      |
| **Closed**       | Default Value\*\* | Panic         | Panic                |
| **Receive-only** | OK                | Compile Error | Compile Error        |
| **Send-only**    | Compile Error     | OK            | OK                   |

---

* `select` ignores a nil channel since it would always block  
* Reading a closed channel returns `(<default-value>, !ok)`

### Rendezvous (unbuffered)
By default, channels are unbuffered (rendezvous model)
* the sender blocks until the receiver is ready (and vice versa)
* the send always happens before the receive
* the receive always returns before the send
* the sender & receiver are synchronized

### Buffering
Buffering channels allows the sender to send without waiting
* the sender deposits its items and returns immediately
* the sender blocks only if the buffer is full
* the receiver blocks if the buffer is empty
* the sender & receiver run independently

Common uses of buffered channels:
* avoid goroutine leaks (from an abandoned channel)
* avoid rendezvous pauses (performance improvement)

Don't buffer until it's needed: buffereing may hide a race condition
(Premature optimization, first always get the program to work properly and the optimize)
Some testing may be required to find the right number of slots

Special uses of buffered channels:
* counting semaphore pattern

Counting Semaphores
A counting semaphore limits work in progress (or occupancy)

Once it's "full" only one unit of work can enter for each one that leaves

We model this with a buffered channel:
* attempt to send (write) before starting work
* the send will block if the buffer is full (occupancy is at max)
* receive (read) when the work is done to free up a space in the buffer (this allows the next worker to start)


# Concurrent File Processing
[Youtube](https://www.youtube.com/watch?v=SPD7TykYy5w&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=27)
[Raw Code](https://github.com/matt4biz/go-class-walk/blob/trunk/walk1/walk1.go)

The video above shows a program to remove duplicate files from a given directory.
It has 4 approaches for doing so.

## Sequential Processing
First approach is simple sequential processing, nothing more and just a straightforward program.

The program flow:
main -> searchFile -> hash -> searchFile -> main

Basically main calls searchFile and searchFile calls hash function to hash content and adds it back to hashMap 
in searchFile and is returned to main for main to start again and print the content of hashMap

## A concurrent approach (like map-reduce)
Use a fixed pool of goroutines and a collector and channels

Basically, the hashing of files is a independent process. Because it is an independent process, goroutines can be used.
One good way to think about developing a goroutine architecture is to start at the independent process and figure out the
inputs and outputs. Then, treat the independent process as a goroutine and its input as a channel and its output as a channel.
For example:
```
                 ┌────────────┐
                 │  searchTree│
                 └─────┬──────┘
                       ↓  (file paths)
                [paths channel]
                       ↓
               ┌───────┴────────┐
               │ Multiple Worker│ goroutines
               │ (processFiles) │
               └───────┬────────┘
                       ↓ (hash + path)
                [pairs channel]
                       ↓
               ┌───────┴────────┐
               │  collectHashes│
               └───────┬────────┘
                       ↓
                [results channel]

```

In this case, processFiles is a independent process and to hash a file, we need the file's path inorder to gets its file 
descriptor. So a paths channel was created to actually pass this into the processfile goroutines and these goroutines output
the hashed pairs which then also have to be outputted to a channel so that others can read from. Its outputted 
to a channel and read from a centralized collectHashes function, instead of goroutines writing directly to the hashMap is because
the hashMap is a shared memory and all goroutines accessed it could cause a race condition. To prevent that, goroutines write to a 
channel and something else can process the outputs from the channel. This is known as 'owner goroutine pattern', where only one 
goroutine owns the mutable state. Finally, its important to actually close the channels and spin down the goroutines properly.
First, once the Walk finsishes sending all the file paths, the paths channel can be closed and the goroutines reading from it will
eventually finish processing and die. We have to wait before closing the pairs channel because these goroutines could still be 
writing to the pairs channel and if pairs is prematurely closed, then these goroutines can get stuck waiting for a channel to write to.
So only after the goroutines finish, the pairs channel can also be closed and once that is closed, the collectHash goroutine can die too.

CPU-bound work → focus on parallelism (distribute work across cores)
I/O-bound work → focus on concurrency (manage multiple operations in progress)


## Sync and Concurrent Subdirectories
The next optimization used was putting subdirectories in its own goroutines. So when a subdirectory is encountered,
we just launch another goroutine and move on to next file/directory. But this would spawn unknown number of goroutines
because we dont know how many subdirectories in can encounter. So using done channel to signal that all goroutines are done 
is not possible because we simply dont know how many goroutines will spawn. So the best approach is to still use the done logic
to know number of goroutines, for example keep it for the processFile workers. But for the dynamic goroutine creation for 
subdirectories, use syn.WaitGroup

Although its not in the code, you can also buffer pairs channel. That way multiple goroutiens can write to the pairs channel
and they wouldn't get blocked when multiple goroutines are trying to write to the same channel. But thats an example of a 
semaphore optimization. 


## Channels as counting semaphores
The optimization idea here was that in the previous optimized code, searchTree function can get blocked
if paths channel is blocked. So one way to solve that, is why not have a goroutine for each file that is 
found? The problem with that, it can exhaust your computers CPU and memory if potentialy tens of thousand
files is found and each tries to process and hash it at the same time. 

So the idea to create one goroutine per file was kept, but what changed is how many actually get to do
work. This is controlled by having the limits buffered channel. Limits channel is buffered up to workers
and only allows up to workers amount to actually do the processFiles work. For example, take a look at 
the different between previous optmized code vs the current:
```Go
// previous optmized code
func processFiles(paths <-chan string, pairs chan<- pair, done chan<- bool) {
	for path := range paths {
		pairs <- hashFile(path) // calls hashFile and then blocks, now bunch of cpu and mem allocated for hashfile while it idly stays blocked
	}

	done <- true
}

// current optimized code
func processFile(path string, pairs chan<- pair, wg *sync.WaitGroup, limits chan bool) {
	defer wg.Done()

	limits <- true 	// blocks before calling hashFile, so goroutine uses minimal CPU and mem

	defer func() {
		<-limits
	}()

	pairs <- hashFile(path)
}
```

Even though both codes can potentially launch as many goroutines as files, the optimized blocks before
any resource usage so the scheduler sort of ignores it while working only with the unblocked ones. 

So yea, the idea is pretty straightforward, just create a buffered limits channel and for each goroutine
have it try to write to it, if its full, its blocked and can't get resources to do the process, and if
its not blocked, just write to it to signal that a worker is in progress and finally drop a data from 
limits to signal that a new spot is free. 

Evalutations:
NOTE: increasing the buffere size to allow for more workers in progress will actually degrade the program.
There is a sweet spot for the number of workers allowed to work because it is bound by disk contetion and I/O.

Amdahl's Law: speedup of a whole program is limited by part that is parallelized
Speedup = 1/(1 - p + (p/n))

Let:
p be the proportion of the program that can be made parallel (e.g., if 70% of a program's execution time can be parallelized, p = 0.70)
(1 - p) is the proportion that is inherently sequential (must be run in order, p = 0.30 in the example)
n is the number of processors (or concurrent workers)

So you can find out how much of ur program is parallel. For example, benchmark Speedup shows 6.25s, n = 8 processors,
that would give p = approximately 96% parellel program


# Conventional Synchronization
Package `sync` contains the conventional synchronization stuff (as in how concurrency is handled in most languages)
For example:
* Mutex
* Once
* Pool
* RWMutex
* WaitGroup

Package `sync/atomic` for atomic scalar reads & writes

We saw a use of WaitGroup in the "file walk" example (above)

The reason why conventional synchronization exists along with Go's CSP model (goroutines and channels model)
is because sometimes multiple goroutines might write or read from a shared data. When this happens, we must make sure 
only **one** of them can do so at any instant (in the so called "critical section")

We accomplish this with some type of lock:
* acquire the lock before accessing the data
* any other goroutine will block waiting to get the lock
* release the lock when done

## Example
```Go
// counting semaphore example
package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println(do())
}

func do() int {
	buffer := make(chan bool, 1)	// we are making a channel to only allow one work at a time
	var n int64			// shared data variable
	var w sync.WaitGroup
	
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			buffer <- true			// limit workers to allow only one worker to do work at a time		
			n++						// all goroutines accessing and modifying the sharen n variable
			<-buffer				// clear buffer once work is done
			defer w.Done()
		}()
	}

	w.Wait()
	return int(n)
}


// the exact idea is achieved by the sync package using Mutex, which locks the variable and only allows one 
// worker to do work at a time. Since this package is meant just for that, it is slightly more efficient compared
// to the manual counting semaphore example above
package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println(do())
}

func do() int {
	var m sync.Mutex		// create the mutex
	var n int64
	var w sync.WaitGroup
	
	for i := 0; i < 1000; i++ {
		w.Add(1)
		go func() {
			m.Lock()		// Lock the state
			n++
			m.Unlock()		// unlock when done incrementing n
			defer w.Done()
		}()
	}

	w.Wait()
	return int(n)
}
```

Mutexes in action:
```Go
type SafeMap struct {
	sync.Mutex			// not safe to copy
	m map[string]int
}

// so methods must take a pointer, not a value
func (s *SafeMap) Incr(key string) {
	s.Lock()
	defer.s.Unlock()

	// only one goroutine can execute this code at the same time, guaranteed
	s.m[key]++
}
// Using defer is good habit, avoids mistakes
```

Mutexes are typically embedded into structs so that the Lock and Unlock becomes available to SafeMap struct
and good practice is to just write Lock and defer Unlock at the beginning so that it locks and unlocks when 
method is done doing its work.

The type of mutex you use can depends on your writes and reads. If you are doing a lot of read operations, there is a 
specific mutex for that. Check docs for mutex optimizations and stuff.

## Only-once execution
A `sync.Once` object allows us to ensure a functino runs only once (only the first call Do will call the function passed in)
```Go
var once sync.Once
var x *singeton

func initialize() {
	x = NewSingleton()
}

func handle(w http.ResponseWriter, r *http.Request) {
	once.Do(initialize)
	...
}
```

Checking x == nil in the handler is unsafe! Best to use the sync.Once method to make sure initialize runs only once

## Pool
A pool provides for efficient & safe reuse of objects, but it's a container of interface{}
```Go
var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func Log(w io.Writer, key, val string) {
	b := bufPool.Get().(*bytes.Buffer) 	// reflection or type assertion?
	b.Reset()
	// write to it
	w.Write(b.Bytes())
	bufPool.Put(b)
}
```

Ofcourse there are also stuff in sync package but above are the most useful to Go probably


# Concurrency Gotchas

## Concurrency Problems
1. Race conditions, where unprotected read & writes overlap
   * must be some data that is written to 
   * could be a read-modify-write operation
   * and two goroutines can do it at the same time

2. deadlock, when no goroutine can make progress
   * goroutines could be all blocked on empty channels
   * goroutines could all be blocked waiting on a mutex
   * GC could be prevented from running (busy loop)

Go detects some deadlocks automatically; with `-race` it can find some data races

3. goroutine leak
   * goroutine hangs on a empty or blocked channel
   * not deadlock: other goroutines make progress
   * often found by looking at `pprof` output

When you start a goroutine, always know how/when it will end

4. channel errors
   * trying to send on a closed channel
   * trying to send or receive on a nil channel
   * closing a nil channel
   * closing a channel twice

5. other errors
   * closure capture
   * misuse of Mutex
   * misuse of WaitGroup
   * misuse of select

A good taxonomy of Go concurrency errors may be found in this paper:
http://cseweb.ucsd.edu/~yiying/GoStudy-ASPLOS19.pdf

Many of the errors are basic & should easily be found by review;
maybe we'll get static analysis tools to help find them


Closure Capture Problem:
A goroutine closure shouldn't capture a mutating variable
```Go
for i := 0; i < 10; i++ { // WRONG
	go func() {
		fmt.Println(i)
	}()
}

// Instead, pass the variable's value as a parameter
for i := 0; i < 10; i++ {
	go func (i int) {
		fmt.Println(i)
	}(i)
}
```

Select Problems:
`select` can be challenging and lead to mistakes
* default is always active (if you put select in for loop with default, it might eat up a lot of CPU)
* a nil channel is always ignored
* a full channel (for send) is skipped over
* a "done" channel is just another channel
* available channels are selected at random (NOT TOP TO BOTTOM!)


Four considerations when using concurrency
1. Don't start a goroutine without knowing how it will stop
2. Acquire locks/semaphores as late as possible; release them in the reverse order
(In other words, "Acquire locks as late as possible" means you should delay grabbing
a lock until the exact moment you need it, and hold it for the shortest time possible.
And releasing them in reverse order is for keeping the order of locks and unlocks 
consistent across other goroutines to prevent deadlocks (dining philosophers problem))
3. Don't wait for non-parallel work that you could do yourself (as in only use goroutines 
when concurrency actually makes sense, dont abuse goroutines)
4. Simplify! Review! Test! 