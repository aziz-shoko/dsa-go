# 🗂 Repository Pattern in Go (with HTTP + In-Memory Example)

## 📌 Overview

The **Repository Pattern** is a way to **decouple your application's data access logic from the business logic** (e.g., HTTP handlers). In Go, it's a powerful but lightweight tool for keeping code clean, modular, and testable.

---

## 💡 Why Use a Repository?

### ✅ Problem with Global State
- Global variables (e.g., `var ID = 1`) are shared, mutable state — they break encapsulation and are unsafe in concurrent programs.
- They tightly couple logic and make testing hard.

### ✅ Solution: Abstraction via Repository
- Create an **interface** that defines what operations can be done (e.g., `AddTweet(t Tweet) (int, error)`).
- Implement that interface in a separate struct (`TweetMemoryRepository`).
- The HTTP handler calls the repository via this interface — **it doesn't care** how the data is stored.

---

## 🧱 Repository Pattern Structure

```go
type TweetRepository interface {
    AddTweet(t Tweet) (int, error)
}

// implementation for in memory data logic (as in store data in memory for now)
type TweetMemoryRepository struct {
    tweets []Tweet
}

func (r *TweetMemoryRepository) AddTweet(t Tweet) (int, error) {
    t.ID = len(r.tweets) + 1
    r.tweets = append(r.tweets, t)
    return t.ID, nil
}
```

🧭 Clean Architecture Principles Applied
1. Separation of Concerns
* HTTP handlers focus only on request/response logic.

* Data storage logic lives in the repository.

2. Dependency Inversion
* High-level modules (e.g., handler) depend on abstractions (TweetRepository), not concrete types.

* Swap TweetMemoryRepository with TweetPostgresRepository later — no handler changes required.

3. Testability
* You can pass a mock implementation of TweetRepository to the server struct during tests.

🧰 Dependency Container: The server Struct
```Go
type server struct {
    repository TweetRepository
}
```
Why use it?
* It acts as a container for external dependencies (repositories, loggers, etc.).

* Keeps function signatures clean.

* Makes the app scalable — just add fields as needed.

🚫 "Don't Use HTTP Types in Repositories"
Repositories should be pure domain logic — no references to:

* http.Request

* http.ResponseWriter

* JSON-specific structs

✅ Accept and return plain Go structs (like tweet), decoupled from protocol-level concerns.

🧠 Practical Workflow (Request to Repository)
HTTP handler receives a POST request.

It unmarshals the JSON body into a userPayload.

It creates a tweet from that payload.

It calls AddTweet(tweet) on the repository.

It gets back an assigned ID.

It sends the ID in the HTTP response.

🛠 Tips for Applying the Pattern
✅ Design your interfaces around what the app needs, not how it's implemented.

✅ Start with in-memory storage for simplicity.

✅ Later, implement the same interface with a DB or external service.

✅ Avoid leaking HTTP concerns into your data logic.

✅ Use dependency injection (like the server struct) to keep things flexible.

🔄 Example Evolution
Phase	What Changes	What Stays the Same
In-Memory	TweetMemoryRepository	Interface + Handler
PostgreSQL	TweetPostgresRepository	Interface + Handler
Testing	MockTweetRepository	Interface + Handler

📚 Further Reading
[https://threedots.tech/post/repository-pattern-in-go/]

✅ Final Thought
Treat your repository like a "data access contract" — not a dumping ground for logic.

Keep layers clean, dependencies injected, and interfaces abstract — your code will stay flexible, readable, and robust.

The full code for reference:
```Go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type userPayload struct {
	Message  string `json:"message"`
	Location string `json:"location"`
}

//	type Payload struct {
//		ID int `json:"ID"`
//	}
//
// I think the best way to store userPayload tweets is to compose it with existing userPayload
type tweet struct {
	userPayload
	ID int `json:"ID"`
}

// var ID = 1
// Below is me trying to decouple data logic (storage) from the http stuff
// the http stuff will add the tweets data through this interface and not directly
// manipulate the underlying storage directly, but by just calling this AddTweet method
type TweetRepository interface {
	AddTweet(t tweet) (int, error)
}

// Below, is now the actual data logic implementation for the interface
// and the implementation of the actual data storage, which is just an in memory Slice of tweet
type TweetMemoryRepository struct {
	tweets []tweet
}

func (d *TweetMemoryRepository) AddTweet(t tweet) (int, error) {
	t.ID = len(d.tweets) + 1
	d.tweets = append(d.tweets, t)
	return t.ID, nil
}

// Server container for dependendcies code below:
// As of right now, i have no idea why we are using a container for dependencies
// and i have no idea what is meant by:
// "It's not a good practice to use HTTP types within the repository. You shouldn't do that in the production code."
type server struct {
	repository TweetRepository
}

type Payload struct {
	ID int `json:"ID"`
}

func main() {
	s := server{
		repository: &TweetMemoryRepository{},
	}

	http.HandleFunc("/tweets", s.addTweet)
	// Your solution goes here. Good luck!
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s server) addTweet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u := userPayload{}

	if err := json.Unmarshal(body, &u); err != nil {
		log.Println("Failed to unmarshal payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.repository.AddTweet(tweet{
		userPayload: u,
	})
	fmt.Printf("Tweet: `%s` from %s\n", u.Message, u.Location)

	// payload := map[string]int{"ID": ID}
	payload := Payload{
		ID: id,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)

}
```