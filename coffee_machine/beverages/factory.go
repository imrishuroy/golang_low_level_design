package beverages

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// NewBeverage is the Factory Pattern entry point: given a BeverageType, it
// returns a ready-to-use Beverage of the matching concrete type (its recipe
// already populated), or an error if the type isn't recognized.
func NewBeverage(beverageType models.BeverageType) (Beverage, error) {
	switch beverageType {
	case models.Coffee:
		return NewCoffee(), nil
	case models.Tea:
		return NewTea(), nil
	case models.Cappuccino:
		return NewCappuccino(), nil
	case models.Latte:
		return NewLatte(), nil
	default:
		return nil, fmt.Errorf("unsupported beverage type: %v", beverageType)
	}
}
