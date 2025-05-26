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
