# How to Program in Go
[youtube](https://www.youtube.com/watch?v=rXgUP_BNyaI&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=41)

## Go Build Tools
We've been using `go run` or maybe `go test` to run programs

Not its time to distribute
* `go build` makes binary
* `go install` makes one and copies it to `$GOPATH/bin`

## Pure Go
In Go you can build "pure" Go programs that literally don't depend on anything (including standard libraries)
This is good for docker images and containers because it makes it very small and eliminates security vulnerabilities
from the dependency packages

```bash
CGO_ENABLED=0 go build -a tags netgo,osusergo \
-ldflags "-extldflags '-static' -s -w" \
-o lister .
```
Here we must tell Go we're going to use pure Go networking

A "pure" Go program can put it into a "from-scratch" container
```bash
ldd lister
# not a dynamic executable
```

## Go build platforms
Go can cross-compile, too
* `$GOARCH` defines the architecture (e.g., amd64 or arm64)
* `$GOOS` defines the operating system (e.g., linux or darwin)
* `$GOARM` for the ARM chip version (v7, etc)

and can build for other platforms liek Raspberry Pi

## Go Project Layout
```
root/
├── README
├── Makefile
├── build/
│   └── Dockerfile
├── cmd/
│   └── programs
├── deploy/
│   └── K8s files
├── go.mod
├── go.sum
├── pkg/
│   └── libraries
├── scripts/
│   └── miscellany
├── test/
│   └── integration tests
└── vendor/
    └── modules
```
Here is a typical Go layout. The pkg directory is where you would have your Go code for your own packages.
NOTE! The Go project layout philosophy is to not have too many nested directories. Make pkg directory if you 
have lots of packages but if you have few packages, it just be in the top level of your module and you don't 
have to nest it inside a pkg folder.

There is also a separate test directory, but not all the tests have to go specifically inside this directory.
Typically, the tests are written alongside the package that is written to have tests against just that package.
Then maybe you can have a separate test directory to do the actual integration tests against the main program or
the program as a whole. 

## Documentation
Your README.md should talk about (among other things)
* overview - who and what is it for?
* developer setup
* project & directory structure
* dependency management
* how to build and/or install it (make targets, etc)
* how to run it (locally, in Docker, etc)
* database & schema
* credentials & security
* debugging monitoring (metrics, logs)
* CLI tools and their usage

## Makefiles
Reasons we might want a Makefile
* we need to calculate parameters
* we have other steps and/or dependencies
* because the options are way too long to type
* and we may have non-Go commands (Docker, cloud provider, etc)

Could be used for versioning the executable
(Everything that can be versioned, should be versioned)
For example:
## Versioning the executable
In the main program code, have an empty variable

```go
// MUST BE SET by go build -ldflags "-X main.version=999"
// like 0.6.14-0-g26fe727 or 0.6.14-2-g9118702-dirty
var version string // do not remove or modify
```
Then we can use makefile to have go build fill in the version string from
git version and branch.
From the makefile
```makefile
VERSION=<span class="math-inline">\(shell git describe \-\-tags \-\-long \-\-dirty 2\>/dev/null\)
BRANCH\=</span>(shell git rev-parse --abbrev-ref HEAD)

xyz: <span class="math-inline">\(SOURCES\)
go build \-mod\=vendor \-ldflags "\-X main\.version\=</span>(VERSION)" -o $@ ./cmd/xyz
```

## Building in Docker
We can use Docker to build as well as run
* multi-stage builds
* use a golang image to build it
* copy the results to a another language image

The result is a small Docker container build for Linux
And you can build it without even having Go installed
This is great for CI/CD environments

REWATCH THE VIDEO FOR MULTI STAGE DOCKER BUILDS!!! 