# Context
[Youtube](https://www.youtube.com/watch?v=0x_oUlxzw5A&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=25)
The `Context` package offers a common method to cancel requests
* explicit cancellation
* implicit cancellation based on a timeout or deadline

A context may also carry request-specific values, such as trace ID

Many network or database requests, for example, take a context for cancellation

## Cancellation and timeouts
A context offers two controls:
* a channel that closes when the cancellation occurs
* an error that's readable once the channel closes

The error value tells you whether the request was cancelled or timed out

We often use the channel from `Done()` in a `select` block

Contexts form an immutable tree structure
(goroutine-safe; changes to a context do not affect its ancestors)

Cancellation or timeout applies to the current context and its subtree
(Ditto for a value)

A subtree may be created with a shorter timeout (but not longer)

## Context Example
The `Context` value should always be the first parameter
```Go
// First runs a set of queries and returns the result from the 
// the first to respond, canceling the others
func First(ctx context.Context, urls []string) (*Result, error) {
    c := make(chanResult, len(urls))        // buffered to avoid orphans
    ctx, cancel := context.WithCancel(ctx)

    defer cancel() // cancel the other queries when we're done

    search := func(url string) {
        c <- runQuery(ctx, url)
    }
    ...
}
```

Another example
```Go
package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func main() {
	results := make(chan result)
	list := []string{
		"https://amazon.com",
		"https://google.com",
		"https://nytimes.com",
		"https://wsj.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	for _, url := range list {
		go get(ctx, url, results)
	}


	for range list {
		r := <- results

		if r.err != nil {
			log.Printf("%-20s %s\n", r.url, r.err)
		} else {
			log.Printf("%-20s %s\n", r.url, r.latency)
		}
	}
}

func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if resp, err := http.DefaultClient.Do(req); err != nil {
		ch <- result{
			url:     url,
			err:     err,
			latency: 0,
		}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		ch <- result{
			url:     url,
			err:     nil,
			latency: t,
		}
		resp.Body.Close()
	}
}
```
In the code above, we create a new context with `context.Background()` (which is an emtpy
root tree node) and it returns a new context, which is basically a another node from the 
root and this new node actually points up to its parent node. Anyway, we then use this new
context to create a context with timeout and cancels the request after 3 seconds. To do so,
we pass it to the get function and make the request with the context that we created and 
passed in. Now if a requgest takes longer than three seconds, it gets caneled with context
deadline exceeded error.

## Context Values
Context can also pass values and it should be data specific to a request, such as:
* a trace ID or start time (for latency calculation)
* security or authorization data

Avoid using the context to carry "optional" parameters

Use a package-specific, private context key type (not string) to avoid collisions






