package beverages

import "golang_low_level_design/coffee_machine/models"

// Beverage is implemented by every drink the machine can prepare.
type Beverage interface {
	Name() string
	Recipe() *models.Recipe
	Prepare()
}
