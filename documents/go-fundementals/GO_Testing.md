# Testing in Go
[youtube](https://www.youtube.com/watch?v=PIPfNIWVbc8&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=38)

## Go test features
Go has standard tools and conventions for testing

Test files end with `_test.go` and have `TestXXX` functions
(they can be in the package directory, or in a separate directory)

You run tests with `go test`

Tests aren't run if the source wasn't changed since the last test

## Layers of Testing
(developer testing, fully automated and integrated with CI/CD)
1. Unit Testing
    * completely self contained
    * dont need the network, microservices, or even databases to test
2. Integration Testing
    * not self contained
    * actually test with/against network, microservices, and database
    * Typically CI/CD (ex, cicd tests the authnetication/authorization to another service)
3. End to End testing
(Tester testing, Exploratory, periodic or tied to release cycles)
4. Load / Performance Testing
5. System Testing
6. Chaos Testing

Goals: Things to test for 
* extreme values
* input validation
* race conditions
* error conditions
* boundary conditions
* pre and post conditions
* randomized data (fuzzing)
* configuration & deployment
* interfaces to other softwares

## Test Functions
Test functions have the same signature using `testing.T`
```Go
func TestCrypto(t *testing.T) {
    uuid := "some-uuid"
    key1 := "some-key1"
    key2 := "some-key2"

    ct, err := secrets.MakeAppKey(key1, uuid)
    if err != nil {
        t.Errorf("make failed: %s", err)
    }
    ...
}
```
Errors are reported through parameter `t` and fail the test

### Table-driven tests
```go
func TestValueFreeFloat(t *testing.T) {
    table := []struct {
        v float64
        s string
    }{
        {1, "1"},
        {1.1, "1.1"},
    }

    for _, tt := range table {
        v := Value{T: floater, V: tt.v, m: &Machine{}}

        if s := v.String(); s != tt.s {
            t.Errorf("%v: wanted %s, got %s", tt.v, tt.s, s)
        }
    }
}
```

### Table-driven subsets
We can run subtests under the parent using `t.Run()`

```go
func TestGraphqlResolver(t *testing.T) {
    table := subTest{
        name string
    }{
        name: "retrieve_offer",
    }

    for _, st := range table {
        // t.Run() takes a name and a closure and runs individual subtest
        t.Run(st.name, func(t *testing.T) { // closure
            ...
        })
    }
}
```

But these kinds of tests can get too big and cumbersome and complicated
and this is when refactoring is useful

### Refactoring tests
Can define some interface
```Go
type checker interface {
    check(*testing.T, string, string) bool
}

type subTest struct {
    name string
    shouldFail bool
    checker checker     // parameterize how we check results
    ...
}

// we can now define idfferent checker types
type check struct {...}

func (c checkGolden) check(t *testing.T, got, want string) bool {
    ...
}
```

### Mocking or Faking
You can mock some databases in tests without actually having
(good for small scale stuff or something)

```Go
type DB interface {
    GetThing(string) (thing, error)
}

type mockDB struct {
    shouldFail bool
}

var errShouldFail = errors.New("db should fail")

func (m mockDB) GetThing(key string) (thing, error) {
    if m.shouldFail {
        return thing{}, fmt.Errorf("%s: %w", key, errShouldFail)
    }
}
```

### Main test functions
You can define a root function for all testing; it will then
run all tests from this point
```Go
func TestMain(m *testing.M) {
    stop, err := startEmulator()

    if err != nil {
        log.Println("***Failed To Start Emulator***")
        os.Exit(-1)
    }

    result := m.Run()   // run all UTs

    stop()
    os.Exist(result)
}
```

### Special test-only packages
If you need to add test-only code as part of a package, you can place it in
a package that ends in _test

That package, like XXX_test.go files, will not be included in a regular build.

Unlike normal test files, it will only be allowed to access exported identifiers, so its useful for "opaque" or "black-box" tests
```go
// file myfunc_test.go
package myfunc_test
// this package is not part of package myfunc, so
// it has no internal access
```

# Philosophy of Testing

## Testing Culture
"Your tests are the contract about what your software does and does not do.
Unit tests at the package level should lock in the behaviour of the package's API.
They describe, in code, what the package promises to do. If there is a unit 
test for each input permutation, you have defined the contract for what the 
code will do in code, not documentation. 

"This is a contract you can assert as simply as typing `go test`. At any stage,
you can know with a high degree of confidence, that the behaviour people
relied on before your change continues to function after your change."
-- Dave Cheney

As a developer, you should assume your code doesn't work unless:
* you have tests (unit, integration, etc)
* they work correctly
* you run them
* they pass

Your work isn't done until you've added or updated the tests
This is the basic code hygiene: start clean, stay clean

"The hardest bugs are those where your mental model of the situation is just
wrong, so you can't see the problem at all." -- Brian Kernighan

This issue applies to testing also

In general, developers test to show that things are done & working
according to their understanding of the problem and solution

Most difficulties in software development are failures of imagination

## Program correctness
There are eight levels of correctness "in order of increasing difficulty of
achievement" (Gries & Conway)

1. it compiles (and passes static analysis)
2. it has no bugs that can be found just running the program
3. it works for some handpicked test data
4. it works for typical, reasonable input
5. it works with test data chosen to be difficult
6. it works for all inputs that follows the specifications
7. it works for valid inputs and likely error cases
8. in works for all input

"It works" means it produces the desired behaviour or fails safely

There are four distinct types of errors (Gries & Conway):
1. errors in understanding the problem requirements
2. errors in understanding the programming language
3. errors in understanding the underlying algorithm
4. errors where you knew better but simply slipped up (everybody makes mistakes)

"Type 1 errors tend to increase as problems become larger, more varied,
and less precisely stated."

Even worse, some requirements may just be missing

## Developer testing is necessary
You should aim for 75-85% code coverage
* unit tests
* integrations tests
* post-deployment sanity checks

Developers must be responsible for the quality of their code
They shouldn't just "throw crap over the wall"

Tests can be part of your documentation

and also, Testing is not "quality assurance"
Confusing "test" and "QA" is a basic mistake
* QA is different discipline in software development
* we're not dealing with a manufacturing process
* you can't "test in" or prove quality

Testing is not about running "acceptance" tests to show that things work
Its about surfacing defects by causing the system to fail (breaking it)
The wrong testing mindset leads to inadequate testing

## Developer testing isn't enough
You can have 100% code coverage and still be wrong
* the code may be bug free, but not match the requirements
* the requirements may not match expectations
* you can't test code thats missing

Testers test to show that things dont work

But they can't test your system well if the requirements aren't documented
(this is a major limitation of the agile method as practiced)

Code & unit tests are simply not enough documentation

## Reality Check
Pick any two
* good
* fast
* cheap

You can't have all three in the real world

Effective and thorough testing is hard & expensive
Software is annoying because most orgs pick fast and cheap over good

# Code coverage
[Youtube](https://www.youtube.com/watch?v=HfCsfuVqpcM&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=39)

Running `goo test -cover` finds what part of the code is exercised by the
tests.
```
$ go test -cover
PASS
coverage: 85.2% of statements
```

Using the `-coverfile` flag generates a file with coverage counts

This can be passed to another tool to didsplay coverage visually
```
$ go tool cover -html=coverage.out
```
Using the `-covermode=count` flag turns it into a heat map