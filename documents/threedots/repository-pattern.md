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


# Further improvements
The repository pattern was nice, but the code above is oversimplified and is done in a single directory with separate files.
In reality, the projects will be in packages and when i learned about Go error handling, i had a really hard time implementing 
a good error handling solution because i didn't know if i should put in the database or http application side. After a lot of 
struggle with coding and arguing with claudeAI and Gemini, here is what is probably the best approach and more clear way to 
decouple stuff. 

(I didnt actually code out the code below (yet), its more to show the conceptual way of decoupling the different logic parts,
i will validate it as i learn more about these concepts)
```Go
// ===== PROJECT STRUCTURE =====
/*
inventory-app/
├── main.go                 // Entry point
├── domain/
│   ├── inventory.go        // Core business logic & errors
│   └── repository.go       // Interface definition
├── storage/
│   └── memory.go          // Database implementation
└── http/
    └── handlers.go        // HTTP layer
*/

// ===== domain/inventory.go =====
// This is the CORE of your application - pure business logic
package domain

import "fmt"

// Domain errors - these represent business problems, not technical problems
type ItemNotFoundError struct {
	Item string
}

func (e ItemNotFoundError) Error() string {
	return fmt.Sprintf("item '%s' not found", e.Item)
}

type InvalidPriceError struct {
	Price string
}

func (e InvalidPriceError) Error() string {
	return fmt.Sprintf("invalid price format: '%s'", e.Price)
}

type EmptyItemNameError struct{}

func (e EmptyItemNameError) Error() string {
	return "item name cannot be empty"
}

// Business entities - what your domain works with
type Item struct {
	Name  string
	Price float64
}

// Business rules - this is where your core logic lives
type InventoryService struct {
	repo Repository
}

func NewInventoryService(repo Repository) *InventoryService {
	return &InventoryService{repo: repo}
}

// These methods contain BUSINESS LOGIC, not implementation details
func (s *InventoryService) CreateItem(name, priceStr string) error {
	// Business rule: names can't be empty
	if name == "" {
		return EmptyItemNameError{}
	}

	// Business rule: prices must be valid numbers
	price, err := parsePrice(priceStr)
	if err != nil {
		return InvalidPriceError{Price: priceStr}
	}

	// Business rule: items must be unique (in this simple example)
	if s.repo.Exists(name) {
		return fmt.Errorf("item '%s' already exists", name)
	}

	// Now delegate to storage layer
	return s.repo.Save(Item{Name: name, Price: price})
}

func (s *InventoryService) UpdateItem(name, priceStr string) error {
	if name == "" {
		return EmptyItemNameError{}
	}

	price, err := parsePrice(priceStr)
	if err != nil {
		return InvalidPriceError{Price: priceStr}
	}

	// Business rule: can only update existing items
	if !s.repo.Exists(name) {
		return ItemNotFoundError{Item: name}
	}

	return s.repo.Save(Item{Name: name, Price: price})
}

func (s *InventoryService) GetItem(name string) (Item, error) {
	if name == "" {
		return Item{}, EmptyItemNameError{}
	}

	item, err := s.repo.FindByName(name)
	if err != nil {
		return Item{}, ItemNotFoundError{Item: name}
	}

	return item, nil
}

func (s *InventoryService) ListItems() ([]Item, error) {
	return s.repo.FindAll()
}

func (s *InventoryService) DeleteItem(name string) error {
	if name == "" {
		return EmptyItemNameError{}
	}

	if !s.repo.Exists(name) {
		return ItemNotFoundError{Item: name}
	}

	return s.repo.Delete(name)
}

// Helper function - business logic for price validation
func parsePrice(priceStr string) (float64, error) {
	// Your business rules for what constitutes a valid price
	// Maybe negative prices aren't allowed, or you have max limits, etc.
	// This is DOMAIN logic, not parsing logic
	return strconv.ParseFloat(priceStr, 64)
}

// ===== domain/repository.go =====
// The interface that defines what storage capabilities the domain needs
package domain

// This interface belongs to the DOMAIN, not the storage layer
// This is the "Dependency Inversion Principle" in action
type Repository interface {
	Save(Item) error
	FindByName(string) (Item, error)
	FindAll() ([]Item, error)
	Delete(string) error
	Exists(string) bool
}

// ===== storage/memory.go =====
// Implementation details - how we actually store data
package storage

import (
	"sync"
	"your-app/domain" // Import the domain to implement its interface
)

type dollars float64

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

// This implements domain.Repository
type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]domain.Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]domain.Item),
	}
}

// Implementation of domain.Repository interface
func (m *MemoryStore) Save(item domain.Item) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.data[item.Name] = item
	return nil // In real life, this might fail due to disk issues, etc.
}

func (m *MemoryStore) FindByName(name string) (domain.Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	item, exists := m.data[name]
	if !exists {
		// Return a generic error - let the domain layer decide what this means
		return domain.Item{}, fmt.Errorf("item not found in storage")
	}
	
	return item, nil
}

func (m *MemoryStore) FindAll() ([]domain.Item, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	items := make([]domain.Item, 0, len(m.data))
	for _, item := range m.data {
		items = append(items, item)
	}
	
	return items, nil
}

func (m *MemoryStore) Delete(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	delete(m.data, name)
	return nil
}

func (m *MemoryStore) Exists(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	_, exists := m.data[name]
	return exists
}

// ===== http/handlers.go =====
// HTTP implementation details - how users interact with your domain
package http

import (
	"errors"
	"net/http"
	"your-app/domain"
)

type Server struct {
	inventory *domain.InventoryService // Uses the domain service
}

func NewServer(inventory *domain.InventoryService) *Server {
	return &Server{inventory: inventory}
}

func (s *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	price := r.URL.Query().Get("price")

	// Call domain service - this is where business logic happens
	err := s.inventory.CreateItem(name, price)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Item created: %s\n", name)
}

func (s *Server) UpdateItem(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	price := r.URL.Query().Get("price")

	err := s.inventory.UpdateItem(name, price)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Item updated: %s\n", name)
}

// Convert domain errors to HTTP responses
func (s *Server) handleError(w http.ResponseWriter, err error) {
	var itemNotFound domain.ItemNotFoundError
	if errors.As(err, &itemNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var invalidPrice domain.InvalidPriceError
	if errors.As(err, &invalidPrice) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var emptyItem domain.EmptyItemNameError
	if errors.As(err, &emptyItem) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Unknown error
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// ===== main.go =====
// Wire everything together
package main

import (
	"log"
	"net/http"
	
	"your-app/domain"
	"your-app/storage"
	// httphandlers "your-app/http"
)

func main() {
	// Create storage layer
	store := storage.NewMemoryStore()
	
	// Create domain service with storage dependency
	inventoryService := domain.NewInventoryService(store)
	
	// Create HTTP layer with domain dependency
	server := httphandlers.NewServer(inventoryService)
	
	// Wire up routes
	http.HandleFunc("/create", server.CreateItem)
	http.HandleFunc("/update", server.UpdateItem)
	// ... other routes
	
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

For Clarification:

What is a Domain? 
A domain refers to the core business problem your software is trying to solve. 

In the case above, the domain is "managing an inventory of items with prices". 
Thats it. Its not about HTTP, databases, or JSON - its more about the business 
rules like:
* Items must have names
* Prices must be valid numbers
* You can't update something that doesn't exist

Think of domain as what your software does, not how it does it.

Domain vs Implementation Layers (Key ideas):
1. Domain is the center
HTTP Layer -> Domain Layer <- Storage Layer

* The domain defines what errors can happen (`ItemNotFoundError`)
* Both HTTP and storage layers know about this domain, but not each other
* The domain doesn't know about HTTP or storage implementations details

2. Dependencies Flow Inward
   * http pacakge imports domain package
   * storage package imports domain package
   * domain package imports nothing application specific (maybe just standard library stuff)

3. Error Handling Flows Through Layers
```Go
// Storage layer: technical error
return fmt.Errorf("item not found in storage")

// Domain layer: business error  
if !s.repo.Exists(name) {
    return ItemNotFoundError{Item: name}  // Domain error!
}

// HTTP layer: presentation error
if errors.As(err, &itemNotFound) {
    http.Error(w, err.Error(), http.StatusNotFound)  // HTTP 404
}
```

So overall, Go Error Handling Princples Are:
1. Errors are values - defined at the right layer (domain)
2. Explicit Error Handling - each layer handles what it can, passed up what it can't
3. Error wrapping - technical errors get wrapped in business errors
4. Type Assertions - `errors.As` to check for specific business errors

