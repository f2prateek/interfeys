package coffeemaker


var _ CoffeeMakerInterface = (*CoffeeMaker)(nil)

type CoffeeMakerInterface interface {
	Brew() error
}
