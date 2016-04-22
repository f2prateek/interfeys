# interfeys

A code generator for Go that generates an interface for a given struct. Useful for generating interfaces for 3rd party packages.

# Installing

`go get github.com/f2prateek/interfeys`

# Example

Given a file.

```go
package coffeemaker

type CoffeeMaker struct {
}

func (c *CoffeeMaker) Brew() {
	...
}
```

Running `interfeys -type CoffeeMaker` will generate:

```go
package coffeemaker

var _ CoffeeMakerInterface = (*CoffeeMaker)(nil)

type CoffeeMakerInterface interface {
	Brew() error
}
```
