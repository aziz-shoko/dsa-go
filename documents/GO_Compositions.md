# Composition

## Philosophy behind Go Composition
Go's composition model follows the principle "favor composition over inheritance".
This approach focuses on what things DO rather than what they ARE.

The core reasons behind this philosophy include:
1. Simplicity and explicitness: Go values code that is straightforward to understand with minimal hidden machanisms
2. Reducing coupling: Inheritance can create tight coupling between parent and child classes
Instead of saying "Type B inherits from Type A" (like in object oriented languages), Go says "Type B has a Type A". This is composition.
3. Practical problem-solving: Go was designed for practical software engineering at scale
4. Interface satisfaction by behavior: Types satisfy interfaces implicitly based on implemented methods

## Embedding: Go's approach to Composition
Go uses a mechanism called "embedding" to achieve composition. You basically embed one struct inside another,
and the outer struct gets access to the inner struct's fields and methods.
```Go
type Engine struct {
    Power int
}

func (e *Engine) Start() {
    fmt.Println("Engine started with power: ", e.Power)
}

type Car struct {
    Engine              // Embedded type (has a relationship)
    Model string
}

// Now Car has access to Engine's fields and methods
```
The example above is not inheritance, the Car is not inheriting Engine's fields and structs 
but just getting access to it, its delegation via embedding and not subclassing.

Under the Hood, when you embed Engine into Car, Go is not saying "Car is an Engine", but instead:
* It stores an anonymoud field of type Engine inside Car
* When you call `car.Start()`, Go delegates that call to car.Engine.Start() behind the scenes.

So Car is still a composite object, and not substype of Engine. A good analogy is a robot with an arm:
* In OOP inheritance: "A robot is an arm" <- doesnt make sense
* In Go composition: "A robot has an arm" <- realistic

So Car has an Engine.

| Aspect         | Inheritance (OOP)            | Composition in Go                        |
| -------------- | ---------------------------- | ---------------------------------------- |
| Relationship   | "Is-a"                       | "Has-a"                                  |
| Type Hierarchy | Builds a class tree          | Flat, no tree                            |
| Method Lookup  | Methods come from base class | Methods are delegated to embedded fields |
| Overriding     | Explicit override chains     | You hide methods by redefining them      |
| Polymorphism   | Classic subtype polymorphism | Interface-based polymorphism             |


## Key Principles of Composition

1. Implicit Method Promotion
Methods of embedded type are "promoted" to the containing type:
```Go
car := Car{
    Engine: Engine{Power: 150}
    Model: "Sedan"
}

car.Start() // Behind the scenes, Go calls Engine.Start() - the method is promoted to Car
```

2. Field Promotion and Name Collisions
Fields from embedded types are also promoted, but explicit fields take precedence:
```Go
type Engine struct {
    Weight int
}

type Car struct {
    Engine
    Weight int  
}
   
car := Car{
    Engine: Engine {Weight: 250}
    Weight: 1500
}
fmt.Println(car.Weight)         // 1500 (Car's field)
fmt.Println(car.Engine.Weight)  // 250 (Must access explicity)
```

3. Multiple Embedding
Go allows embedding multiple types:
```Go
type ElectricMotor struct {
    Voltage int
}

func (e *ElectricMotor) Charge() {
    fmt.Println("Charging at ", e.Voltage, "volts")
}

type HybridCar struct {
    Engine
    ElectricMotor
    Model string
}

// HybridCar has access to both Engine.Start() and ElectricMotor.Charge()
```

4. Interface Satisfaction Through Embedding
A type can satisfy interfaces through methods of embedded types:
```Go
type Starter interface {
    Start()
}

// Car implicitly satisfies Starter because Engine implements Start()
func StartVehicle(s Starter) {
    s.Start()
}

car := Car{
    Engine: Engine{Power: 150}
}
StarVehicle(&car) // Works because Car embeds Engine, which implements Start()
```

5. Method Overriding
You can override methods from embedded types:
```Go
func (c *Car) Start() {
    fmt.Println("Car starting sequence for model: ", c.Model)
    c.Engine.Start()    // Call the embeded type's method if needed
}
```

6. Embedding vs Named fields
Both named fields and embedding are forms of composition in Go.
But, there is an important distinction between Embedding and Named fields
```Go
type Car struct {
    Engine          // Embedded
    Model string
}

type BetterCar struct {
    engine Engine  // named field (not embedded)
    Model  string
}
```

| Feature                | Embedded Field (`Engine`)                               | Named Field (`engine Engine`)                  |
| ---------------------- | ------------------------------------------------------- | ---------------------------------------------- |
| Method/field promotion | ✅ Promoted to outer struct                              | ❌ Only accessed via `car.engine.Method()`      |
| Readability            | ❗ Can blur type boundaries                              | ✅ Makes relationships explicit                 |
| Encapsulation          | ❌ Breaks it (outer type exposes inner methods directly) | ✅ Keeps implementation detail hidden           |
| Interface satisfaction | ✅ Promotes methods to outer struct                      | ❌ Outer type must forward methods explicitly   |
| Composition use-case   | ✅ Useful when you want to mimic behavior sharing        | ✅ Better for delegation or has-a relationships |
| Idiomatic?             | ✅ For reuse and behavior exposure                       | ✅ For clear abstraction and separation         |
| Overriding             | ✅ Can override by defining method on outer type         | ❌ Needs explicit wrapping/forwarding           |

Use Embedding When:
* You want to promote behavior into the outer type
* The embedded type's methods should feel like they belong to the outer type
* You doing behavior composition, like `Logger`, `Sync.Mutex`, or `Context`.

But embedding blurs boundaries. If you embed something like Engine into Car, you're saying "Car behaves like Engine in some way".

Use a named field when:
* You want clear ownership or abstraction
* You want to hide internal implementation from users of the outer type
* You're modeling a real world-relationship (like Car has an Engine)

This is more idomatic when behavior separation matters, or the inner type shouldn't be exposed.

7. Compostion with interfaces
There is also composition with interfaces
```Go
// Approach 1: Embedding concrete type
type APIHandler struct {
    Logger
    // API-specific fields
}

// Approach 2: Implementing an interface (often preferred)
type LoggerInterface interface {
    Log(message string)
}

type CustomLogger struct{}

func (l *CustomLogger) Log(message string) {
    fmt.Println("[CUSTOM]:", message)
}

type APIHandler struct {
    logger LoggerInterface
    // API-specific fields
}

func (a *APIHandler) DoSomething() {
    a.logger.Log("Doing something")
}

// Then somewhere else in the main function
func main() {
    handler := APIHandler{
        logger: &CustomLogger{}, // <<< THIS is the wiring you're wondering about
    }

    handler.DoSomething() // Calls CustomLogger.Log
}
```

So composition with interfaces happen because it offers type flexibility. For example, APIHandler has a interface for its field.
Now, anything that can satisfy the `logger` field can be composed into APIHandler. In the case above, it is CustomLogger struct
which has the Log method to satisfy the interface. So then when APIHandler is constructed in the main, we can pass in CustomLogger
and CustomLogger now becomes available to APIHandler struct.

8. Composition with pointer types
A struct can embed a pointer to another type; promotion of its fields and methods works the same way. 

## Go Sortable Interface
`sort.Interface` is defines as

```Go
type Interface interface {
    // Len is the number of elements in the collection
    Len() int

    // Less reports whether the element with 
    // index i should sort before the element with index j
    Less(i, j int) bool

    // Swap swaps the elements with indexes i and j
    Swap(i, j int)
}

// and Sort.Sort as
func Sort(data Interface) {
    ...
}
```

A practical example:
```Go
package main

import (
	"fmt"
	"sort"
)

type Organ struct {
	Name string
	Weight int
}

type Organs []Organ

func (s Organs) Len() int {
	return len(s)
}

func (s Organs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ByName struct {
	Organs
}

type ByWeight struct {
	Organs
}

func (s ByName) Less(i, j int) bool {
	return s.Organs[i].Name < s.Organs[j].Name
}

func (s ByWeight) Less(i, j int) bool {
	return s.Organs[i].Weight < s.Organs[j].Weight
}

func main() {
	s := []Organ{
		{"brain", 1340},
		{"liver", 1494},
		{"spleen", 162},
		{"pacreas", 131},
		{"heart", 290},
	}

	sort.Sort(ByWeight{s})
	fmt.Println(s)
	sort.Sort((ByName{s}))
	fmt.Println(s)
}
```

## Making the nil value useful
Nothing in Go prevents calling a method with a nil receiver
```Go
type IntList struct {
    Value int
    Tail *IntList
}

// Sum returns the sum of the list elements
func (list *IntList) Sum() int {
    if list == nil {
        return 0
    }

    return list.Value + list.Tail.Sum()
}
```

When we reach the end of the linked list, list can be nil and we can recursively call Sum again 
with that nil node and no errors will be thrown.
