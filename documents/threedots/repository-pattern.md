# Repository Pattern in Go (with HTTP + In-Memory Database Example)

Further Reading
[https://threedots.tech/post/repository-pattern-in-go/]

## Overview
The respository pattern is basically a way of decoupling your application logic from the database logic.
Typically its done by:
* having an interface in between application logic and data logic to represent the repository contract
* then the database logic is separated into a different file and they become the producers for the interface
* since the application logic are the consumers, they will define the interface, because they will be determining what behaviors they need
* Then to tie them together, typically a dependency struct is created using that interface (and a constructor to initailize the struct)
* and then a concrete implementation is defined and given to the dependency struct so its methods can basically call the interface's methods
  
Below is a code example of a simple database server that handles list, add, update, fetch and drop
operations. The database logic and application logic was separated out into different files for clarity

```Go
// main.go file
package main

import (
	"fmt"
	"log"
	"net/http"
)

type DataRepository interface {
	AllItems() (map[string]dollars, error)
	CreateItem(string, string) error
	UpdateItem(string, string) error
	GetItem(string) (string, error)
	DeleteItem(string) error
}


// For holding all application dependencies
type Server struct {
	repo DataRepository // hold something that is a DataRepository
	// could have other dependencies here, like logger or something
}

// common practice to return mem addresses in constructors
// because the instances will be passed around and any changes
// during the passing around should reflect on that instance
func NewServer(r DataRepository) *Server {
	return &Server{repo: r}
}

func (s *Server) list(w http.ResponseWriter, req *http.Request) {
	db, _ := s.repo.AllItems()
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (s *Server) add(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	err := s.repo.CreateItem(item, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	err := s.repo.UpdateItem(item, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("Updated %s for %s\n", item, price)
}

func (s *Server) fetch(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	price, err := s.repo.GetItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item found: %s\n", item)
	fmt.Fprintf(w, "Price of item %s is %s\n", item, price)
}

func (s *Server) drop(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	err := s.repo.DeleteItem(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item removed: %s\n", item)
	log.Printf("Item removed: %s\n", item)
}

func main() {
	// create the conrete database
	// the goal is to then pass it to the Server struct that holds all the dependencies
	// and have the handlers call handlers through the Server struct
	db := database{
		"shoes": 50,
		"socks": 5,
	}
	serverHandlerService := NewServer(db)

	http.HandleFunc("/list", serverHandlerService.list)
	http.HandleFunc("/create", serverHandlerService.add)
	http.HandleFunc("/update", serverHandlerService.update)
	http.HandleFunc("/delete", serverHandlerService.drop)
	http.HandleFunc("/read", serverHandlerService.fetch)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
```

Based on the repository pattern, the code above is the consumers side and separates the http logic 
from the business logic. All the http handlers interact with the database based on the interface's 
methods, so this part of the code doesn't have to worry at all about the database logic. It just knows
that the interface will give it something to listitems, somthing to additems, etc. 

Notice how to tie the logical parts together, the dependency struct `Server` was created. It holds 
the dependency to an interface, as in it holds anything that can implement the interface. To satisfy 
this struct, we create the actual concrete type (database) which has all the methods that the interface
is looking for and use a constructor to implicitly launch the dependency struct with the concrete type.
(refer to comments for some finer key details)

```Go
// database.go
package main

import (
	"fmt"
	"strconv"
	"log"
)

// NOTE: don't do this in real life
type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) AllItems() (map[string]dollars, error) {
	return db, nil
}

func (db database) CreateItem(item string, price string) error {
	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		log.Printf("Invalid price format received: %s\n", price)
		return fmt.Errorf("Invalid price format: %s\n", price)
	}

	db[item] = dollars(priceFloat)
	log.Printf("Added %s for %s\n", item, dollars(priceFloat).String())	

	return nil
}

func (db database) UpdateItem(item, price string) error {
	if _, ok := db[item]; ok == false {
		log.Printf("Item not found: %s\n", item)
		return fmt.Errorf("Item not found: %s\n", item)
	}

	priceFloat, err := strconv.ParseFloat(price, 32)
	if err != nil {
		log.Printf("Invalid price format received: %s\n", price)
		return fmt.Errorf("Invalid price format: %s\n", price)
	}

	db[item] = dollars(priceFloat)
	return nil
}

func (db database) GetItem(item string) (string, error) {
	if _, ok := db[item]; ok == false {
		log.Printf("Item not found: %s\n", item)
		return "", fmt.Errorf("Item not found: %s\n", item)
	}

	return db[item].String(), nil
}

func (db database) DeleteItem(item string) error {
	delete(db, item)
	return nil
}
```
The code above is now the database logic, it literally just handles the database logic part and thats it.
Notice how nothing in the database logic uses HTTP types, just purely data sent from the http handlers 
in the application logic.

This implementation of database and http logic server is very basic. The point is to just show the 
repository pattern from threedots.tech and thats why it has very basic functionality and no concurrency