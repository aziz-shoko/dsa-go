# Go Modules
Go module support is intended to solve several problems:
* avoid the need for $GOPATH
* group packages versioned/released together
* support semantic versioning & backwards compatibility
* provide in-project dependency management
* offer strong dependency security & availability
* continue support of vendoring
* work transparently across the Go ecosystem

Go modules with proxying offers the value of vendoring without requiring your project 
to vendor to all 3rd-party code in your repo

Basically Go has a proxy server running 24/7 that handles all the module installs (check the checksums,
cache, etc). When it caches a module (copy), it permanently keeps it so that even if the original repo 
is deleted, it can still serve you the module from the copy that it had. Then in your machine, Go has copy
of this cache available to the whole system and each Go project then uses whatever module that it may need 
from it. 

## Why modules?
Go's dependency management protects against some risks:
* flaky repos 
* packages that disappear
* conflicting dependency versions
* surreptitious changes to public packages

But it cannot ensure the actual quality or security of the original code; see
* Reflections on Trusting Trust by Ken Thompson
* Our Software Dependency Problem by Russ Cox

"A little copying is better than a little dependency" - Go Proverb

## Import compatibility rule
"if an old package and a new package have the same import path, 
the new package must be backwards compatible with the old package"

An incompatible updated package should use a new URL (version)
```go
package hello

import (
    "github.com/x"
    x2 "github.com/x/v2"
)

func main() {
    x.SomeFunction()    // calls function from github.com/x
    x2.SomeFunction()   // calls function from github.com/x/v2
}
```
Note that you can import both versions if necessary

## Some control files
The `go.mod` file has your module name along with direct dependency requirements
(from Go 1.13, the version of Go)

```
module hello
require github.com/x v1.1
go 1.13
```

The `go.sum` file has checksums for all transitive dependencies
Always check them in to your repo (commit and push)

## Some environment variables
We typically use defaults for these (goes through the default proxy server)
```
GOPROXY=https://proxy.golang.org,direct
GOSUMDB=sum.golang.org
```

and set this for private repos (goes through the private proxy server set up by your company)
```
GOPRIVATE=github.com/xxx,github.com/yyy
GONOSUMDB=github.com/xxx,github.com/yyy
```

Remember also you must be set up for access to private github repos in order to
download private modules. 

## Maintaining dependencies
What do you do on daily bases?

Start a project with:
```bash
go mod init <module-name>   ## create the go.mod file
go build                    ## building updates go.mod
```

Once a version is set, Go will not update it automatically; you can
update every dependency with

```bash
go get -u ./...             ## update directly
go mod tidy                 ## remove unneeded modules
```

You must commit the `go.mod` and `go.sum` files in your repo

You can list avaialable versions of a dependency
```bash
go list -m -versions rsc.io/sampler
```
Then to update a single dependency, you can use go get

You can also vendor in Go:
Vendoring means copying all your dependencies' source code directly into your project 
repository, so everything needed to build your project is self-contained.

Use `go mod vendor` to create the vendor directory; it must be in the module's root directory (along with go.mod)

Go keeps a local cache in `$GOPATH/pkg`
* each package (using a directory tree)
* the hash of the root checksum DB tree

Use `go clean -modcache` to remove it all (cleanup)