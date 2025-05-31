# Generics in Go
[youtube](https://www.youtube.com/watch?v=Si0rAE8yT9g&list=PLoILbKo9rG3skRCj37Kn5Zj803hhiuRK6&index=42)

"Generics" is shorthand for parametric polymorphism

That means we have a type parameter on a type or function
```go
type MyType[T any] struct {
    v T     // can be any valid Go type
    n int
}
```

Generics are a powerful feature for abstraction
And a possible source of unnecessary abstraction and complexity

When should Generics be used? to replcae dynamic typing
Use type parameters to replace dynamic typing with static typing
By making dynamic typing into static typing using generics, we can make
the program safer by checking its type before it runs as opposed to figuring
out the type when the program is running (dynamic typing, type known during
runtime). This is the real value of generics.
```Go
interface {} + v.(T)
// change that to this
type MyType[T any] struct{...}
```
If it runs faster, consider that a bonus
Continue to use (non-empty) interfaces wherever possible
Performance should not be your principle reason for generics (in most cases)

## Generic type & function
```go
type Vector[T any] []T

func (v *Vector[T]) Push(x T) {
    *v = append(*v, x)      // may reallocate
}

// note: F and T are both used in the parameter list
func Map[F, T any](s []F, f func(F) T) []T {
    r := make([]T, len(s))

    for i, v := range s {
        r[i] = f(v)
    }
    return r
}
```