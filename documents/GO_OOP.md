# Go does OOP
This is a documentation on how Go does OOP (My personal notes from the goat: [link](https://www.youtube.com/watch?v=jexEpE7Yv2A&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=17))

## What is OOP?
For most people, the essential elements of OOP have been:
* abstraction
* encapsulation
* polymorphism
* inheritance

Sometimes those last two items are combined or confused

Go's approach to OOP is similar but different

### Abstraction
Abstaction: decoupling behavior from the implementation details

The Unix File system API is a great example of effective abstraction
Roughly five basic functions hide all the messy details:
* open
* close
* read
* write
* ioctl

All the complexities and inner workings of opening a file, closing it, reading it is hidden away and we only care about what it does.
Many different systems things can be treated like files

### Encapsulation
We can't really make abstraction work without encapsulation practically

Encapsulation: hiding implementation details from misuse

It's hard to maintain an abstraction if the details are exposed:
* the internals may be manipulated in ways contrary to the concept behind the abstraction
* users of the abstraction may come to depend on the internal details -- but those might change

Encapsulation usually means controlling the visibility of names ("private" variables)

Basically from my understanding so far, abstraction makes it easier to work with systems because it dumbs down the user interaction to something simple (e.g. open files). But to implement this abstaction, we need to use encapsulation where certain parts are hidden away so that it can't be used in a different way other than it was intended to. 

We want the abstraction to be protected in two different ways:
1. The user of the abstaction must not know the details behind it and must not depend on those details
2. The variables behind the abstraction must not change (as can't be reassigned to another variable or values) but only interacted through methods. This is because at beginning of programming, variables were globally declared and programmers just changed it wherever which led to things breaking. By doing this abstaction stuff, it somewhat gives a standardized structure of how things should be done. 

### Polymorphism
Polymorphims literally means "many shapes" -- multiple types behing a single interface

Three main types are recognized:
* ad-hoc ("ad hoc" refers to solutions that are developed specifically for a particular problem or task, without considering broader applications): typically found in function/operator overloading 
* parametric: commonly known as "generic programming"
* subtype: subclasses substituting for superclasses

"Protocol-oriented" programming uses explicit interface types, now supported in many popular language (an ad-hoc method)

In this case, behavior is completely separate from implementation, which is good for abstraction. 


### Inheritance
Inheritance has conflicting meanings:
* substitution (subtype) polymorphism
* structural sharing of implementation details

In theory, inheritance should always imply subtyping:
the subclass should be a "kind of" the superclass (circle subclass of shape superclass)

See the Liskov substitution principle

Theories about substituition can be pretty messy

#### Why would inheritance be bad?
It injects a dependence on the superclass into the subclass:
* What if the superclass changes behavior?
* What if teh abstract concept is leaky?

Not having inheritance means better encapsulation & isolation

"Interfaces will force you to think in term of communication between objects"
-- Nicolo Pignatelli in Inheritance is evil

See also Composition over inheritance and Inheritance tax (Pragmatic)

"Object-oriented programming to me means only messaging, local retention and protection and hiding of state-processes, and extreme late-bindings of all things." - Alan Kay

Alan Kay wrote this to:
* de-emphasize inheritance hierarchies as a key part of OOP
* emphasize the idea of self-contained objects sending messages to each other
* emphasize polymorphism in behavior

## OO in Go
Go offers four main supports for OO progamming:
* encapsulation using the package for visibility
* abstraction & polymorphism using interface types
* enhanced composition to provide structure sharing

Go does not offer inheritance or subsstitutability based on types

Substitutability is based only on interfaces: purely a function of abstract behavior

See also Go for Gophers

### Classes in Go
Not having classes can be liberating!

Go allows defining methods on any user-defined type, rather than only a "class"

Go allows any object to implement the method's of an interface, not just a subclass


## Interfaces

An interface specifies abstract behavior in terms of methods. 

Basically its a collection of methods to my knowledge so far. And it automatically maps to the right method based on passed in object.

```
type Stringer interface {
    String() string
}
```

Concrete types offer methods that satisfy the interface. In other words, interface lists a list of methods that an object with concrete type must provide so that an interface can map to that object.



