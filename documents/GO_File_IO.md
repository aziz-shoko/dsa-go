# Standard I/O

## What are Streams?
A "stream" in computing is a sequence of data elements made available over time. Its like water flowing through a pipe, the data flows from one point to another. The key characteristics are:
1. Sequential access - Data is processed in order, one element after another
2. Flow - Data moves from a source to a destination
3. Abstraction - The mechanism of how data moves is hidden from you

## The Three Standard Streams
Unix systems define three standard streams that every program has access to by default:
1. Standard Input (Stdin) - The stream of data going into your program
2. Standard Output (Stdout) - The stream of data coming out of your program (for normal output)
3. Standard Error (Stderr) - A seperate output stream specifically for error messages

### Input Stream (Stdin)
Standard input is where your program reads data from. By default, this is ocnnected to the keyboard, but it's much more than that:
* Can be the characters you type on the terminal
* Can be the dat piped from another program
* Or it can be redirected from a file

The niceness of streams is that your program doesn't need to care where the data comes from (abstraction). The program just looks and reads from stdin, and the 
operating system handles the complexity behind it.

### Output Stream (Stdout)
Standard output is where your program writes its normal output. By default, this displays on the terminal, but:
* It could be redirected to a file
* It could be piped as input to another program
* It could be sent to a device like a printer

Here too, the program just writes to stdout and doesn't have to worry about where that data ultimately goes.

### Error Stream (stderr)
Standard error is a separate output stream specifically for error messages. It exists separately from stdout so that:
* Error messages can be visible even if normal output is redirected
* Errors can be logged separately from normal output
* Error handling can be managed differently

## Why this abstraction Matters
The stream concept gives the separation of concerns:
1. Your program just reads from stdin and writes to stdout/stderr
2. The shell/OS decided where those streams actually connect to

This means you can:

* Chain programs together using pipes: program1 | program2 | program3
* Redirect input from files: program < input.txt
* Redirect output to files: program > output.txt
* Redirect only errors: program 2> errors.log

And your program's code doesn't change at all!

## Buffering
One more important concept: streams are often buffered. Thsi means:
* Data might be collected in chunks before being processed
* Output might be held temporarily before being displayed
* This improves efficiency by reducing the overhead of handling each tiny piece separately

So basically, abstraction is a big thing in Unix systems. With streams, our program just has to worry about looking at Stdin and see which data it got. The complex parts of
where the data then goes is just abstracted away by the OS itself. Same thing for the other two outputs (Stdout, Stderr)

# Formatted I/O
Formatted I/O in Go means the ability to read and write data with specific formatting, basically controlling exactly how the data is presented when it is outputted
and how it's interpreted when the data in inputted. With formatted I/O, instead of working with raw bytes moving through the stream, you're working with structured, formatted data.

In Go, formatted I/O is primarily handled through the fmt package, which provides functions that let you control the format of data that you read form input streams 
or write to output streams.

## Output Formatting
The most common formatted output functions are:
* `fmt.Printf()` - Formats and writes to standard output (stdout)
* `fmt.FprintF()` - Formats and writes to a specified writer (could be a file, network, connections, etc.)
* `fmt.Sprintf()` - Formats and returns a string (doesn't write to any stream)

These functions use format specifiers (like %d, %s, %f) to control how values are formatted:
```Go
name := "Alice"
age := 30
fmt.Printf("Name: %s, Age: %d\n", name, age)
// Output: Name: Alice, Age: 30
```

## Input Formatting
For reading formatted input:

* `fmt.Scanf()` - Scans formatted text from standard input (stdin)
* `fmt.Fscanf()` - Scans formatted text from a specified reader
* `fmt.Sscanf()` - Scans formatted text from a string

These also use format specifiers to interpret incoming data:
```Go
var name string
var age int
fmt.Scanf("%s %d", &name, &age)
// If user inputs "Bob 25", name will be "Bob" and age will be 25
```

# File I/O
Package `os` has functions to open or create files, list directories, etc. and hosts the File type

Package `io` has utilities to read and write; `bufio` provides the buffered I/O scanners, etc.

Package `io/ioutil` has extra utilities such as reading an entire file to memory, or writing it out all at once

Package `strconv` has utilities to convert to/from string representations


## Cat Command Clone
```Go
package main

import (
	"fmt"
	"os"	
	"io"
)

func main() {
	for _, fname := range os.Args[1:] {
		file, err := os.Open(fname)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue // continue with the next file if this file fails to open
		}

		if _, err := io.Copy(os.Stdout, file); err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}

		file.Close()
	}
}
```
In Go, always check the error!!! Functions in Go often return 2 types, value and error, make sure to always check and handle the errors.

Basically, to my understanding so far, reading from Files works this way. Typicall you need the file descriptor that is a pointer to the actual file itself
and you get this by os.Open method that returns an object of type `*os.File` that has the file descriptor and metadatas related to the file. The actual
content reading from the disk to memory happens with other methods like `io.Copy` from `io` package. There are different approaches to how the data is read,
one approach is to read load the entire file to memory but thats inefficient or another approach is load it by individual bytes that is also inefficient as each
read will have to make a read call to OS. The most common and more efficient way to read is through buffers which read data from files in chunks. Typically
programmers know what type of data format and type they are working with in Files, if they dont, they typically use some input stream methods that allows
them to see what type of data they are workings with, how to parse it and format them into what kind of data they want it to be outputted. Please checks docs
on the specifics on these types of methods.

Fprintf is also used to specify which input stream you want it to be outputted. Common stdout methods like `fmt.Printf` can only output to Stdout and if you
want to specify a different input stream, you must use `Fprintf` statements with `F` in front that allows you to choose which stream to output. So for example,
`fmt.Fprintln(os.Stderr, err)` outputs err to the Stderr stream and Fprintf statements is commonly used with dealing with errors.


## wc command clone
```Go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for _, fname := range os.Args[1:] {
		var lc, wc, cc int

		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		scan := bufio.NewScanner(file)

		for scan.Scan() {
			s := scan.Text()

			wc += len(strings.Fields(s))
			cc += len(s)
			lc++
		}

		fmt.Printf("%5d %5d %5d %s\n", lc, wc, cc, fname)
		file.Close()
	}
}
```

Create a scanner with bufio.NewScanner(someReader), which sets up a buffer behind the scenes
Use scanner.Scan() in a loop - it returns true while there's more data and false when done
Inside the loop, use scanner.Text() to get the current line as a string