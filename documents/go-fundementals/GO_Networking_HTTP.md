# Networking with HTTP
The Go was made to work with the cloud and it has many standard library packages for making web servers.

That includes:
* client & server sockets
* route multiplexing
* HTTP and HTML, including HTML templates
* JSON and other data formats
* cryptographic security
* SQL database access
* compression utilities
* image generation

There are also lots of 3rd-party packages with improvements

## Web Server Fundamentals

### What is a webserver
A web server is a program that listens for requests over a network (typically the internet), processes those requests,
and returns responses. The most common protocol used for this communication is HTTP. 

### HTTP: The foundation
HTTP is a request-response protocol that forms the foundations of data communiations on the web:
1. Request: A client (usually a browser) sends a request to a server
2. Response: The server processes the request and sends back a response

Each HTTP contains:
* A method (GET, POST, PUT, DELETE, etc.)
* A url path (like `/users`)
* Headers (metadata about the request)
* Sometimes a body (data sent with the request)

Each HTTP response contains:
* A status code (200 OK, 404 Not Found, 500 server error, etc)
* Headers (metadata about the response)
* Often a body (the actual content being returned)

### Basic Components of a Web Server
1. `Listener`: A component that continously listens for incoming connections on a specific port
2. `Router/Handler`: Logic that determines what code runs based on the incoming request path
    * In the code below, `http.HandleFunc("/", handler)` maps the root path to your handler function
3. `Request Handler`: Code that processes the request and generates a response
    * In the code below, the `handler` function that writes "Hello, World!" to the response
4. `Response Writer`: Component that sends data back to the client
    * In the code below, the `w http.ResponseWriter` parameter

```Go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world! from %s\n", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### HTTP Request-Response Lifecycle
1. Client initiates request: A browser or other clients sends an HTTP request to your server
2. Server receives request: You web browser program receives the raw HTTP
3. Server routes request: The server determines which handler should process the request based on the URL path
4. Handler processes request: Your code runs, processes the request, perhaps interacts with a database
5. Response creation: Your code creates a response (HTML, JSON, plain text, etc.)
6. Server sends response: The response is sent back to the client
7. Client processes response: The browser renders the HTML, or an app processes the JSON, etc.

### Middleware
Middleware is code that runs between receiving the request and executing your route handlers. It's like a pipeline of functions that each request flows through.
Common uses for middleware:

* Logging: Recording information about each request
* Authentication: Verifying user identity
* Authorization: Checking if the user has permission
* CORS handling: Managing cross-origin requests
* Body parsing: Converting request bodies (like JSON) into usable data

In Go's standard library, middleware can be implemented using handler wrapping patterns.

### REST APIs
REST (Representational State Transfer) is an architectural style for designing networked applications. A RESTful API uses HTTP methods explicitly and is stateless.
Key concepts:

* Resources: Identified by URLs (e.g., /users, /products)
* HTTP Methods:
  * GET: Retrieve a resource
  * POST: Create a new resource
  * PUT/PATCH: Update a resource
  * DELETE: Remove a resource


* Statelessness: Each request contains all the information needed to complete it
* Representation: Resources can be represented in different formats (JSON, XML, etc.)

For example, a RESTful API might have these endpoints:

* GET /users - List all users
* POST /users - Create a user
* GET /users/123 - Get user with ID 123
* PUT /users/123 - Update user 123
* DELETE /users/123 - Delete user 123

### Building a More Complete Web Server
A more complete web-server would typically have:
1. Structured routing: Handling different paths and HTTP methods
2. Data parsing: Processing request bodies (JSON, form data)
3. Data serialization: Converting data to JSON or other formats for responses
4. Error handling: Properly managing and reporting errors
5. Database connectivity: Storing and retrieving persistent data
6. Middleware chain: For cross-cutting concerns like logging


# Go HTTP Design
An HTTP handler function is an instance of an interface
```Go
// internal HTTP package
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}

// user defined function (outside of HTTP package)
// The HTTP framework can call a method on a function type
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World! from %s\n", r.URL.Path[1:])
}
```
So here is an explanation based on my understanding of Go's OOP principles:

Internally in http package: 
1. the package defines the `handler` interface
2. the package defines the `HandlerFunc` type with a `ServeHTTP` method
3. this allows functions to be treated as handlers inside the http package

As the user of the package:
* the user simply writes handler functions like `func handler(w http.ResponseWriter, r *http.Request)`
* You register this with `http.HandleFunc("/", handler)
* Behind the scenes, `HandleFunc` converts your function to a `HandleFunc` type, which then implements the `Handler` interface by making the ServeHTTP method

The main idea to understand about this design is that in Go, you have to Program to interfaces, not implementations.
Basically, other components inside the http package also expect interfaces and not concrete types. The user gives the concrete type to 
http package and then this concrete type gets converted to a Handler interface by having a func type HandlerFunc implement the serveHTTP method
using users concrete type. Once the Handler is implemented, the HTTP server works with the user function through the Handler interface.
So in summary:
1. The HTTP server only cares that someting can handle a request (via `ServeHTTP`)
2. It doesn't care how that's implemented
3. This allows for flexibility of different components like middleware, routers, and simple functions to all be handlers

### The function types
So interfaces accept only methods. The user is giving a concrete type of function and instead of creating a seperate struct and then assigning a method to it,
in Go you can have function types and assign methods to it. Basically, you make the function type adapt to interface by creating a type and a method for that type.

In most OOP languages, to adapt a function to an interface, you'd need to:
1. Create a struct/class
2. Add your function as a method
3. Implement the interface

Go's function types let you skip this boilerplate by:
1. Defining `HandlerFunc` as a function type
2. Adding the `ServeHTTP` method to that type
3. Having `ServeHTTP` simply call the function

Go allows a simply function to be easily adapted to an interface. It's clean, lightweight solution that avoids unnessary structs

Example Code:
```Go
package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type todo struct {
	UserID	  int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var form = `
<h1>Todo #{{.ID}}</h1>
<div>{{printf "User %d" .UserID}}</div>
<div>{{printf "%s (completed: %t)" .Title .Completed}}</div>`

func handler(w http.ResponseWriter, r *http.Request) {
	const base = "https://jsonplaceholder.typicode.com/"
	resp, err := http.Get(base + r.URL.Path[1:])

	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	defer resp.Body.Close()
	var item todo


	if err = json.NewDecoder(resp.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	tmpl := template.New("mine")

	tmpl.Parse(form)
	tmpl.Execute(w, item)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
The code above is a simple server that when url with path `/todos/3` is typed, it makes a request to typicode, gets json, translates it to Go struct, 
and parses it to web browser html to show it. 
