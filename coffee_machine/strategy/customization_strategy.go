package strategy

import "golang_low_level_design/coffee_machine/beverages"

// CustomizationStrategy customizes a beverage in some way (e.g. adjusting
// milk type or sugar level) before it's prepared. Each concrete strategy
// implements this independently, so CoffeeMachine can apply any combination
// of customizations to a beverage without knowing the specifics of each one.
type CustomizationStrategy interface {
	Customize(beverage beverages.Beverage)
	Description() string
}
