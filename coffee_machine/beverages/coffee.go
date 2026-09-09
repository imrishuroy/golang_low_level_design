package beverages

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// Coffee is a basic black coffee.
type Coffee struct {
	name   string
	recipe *models.Recipe
}

// NewCoffee creates a Coffee beverage with its recipe already populated.
func NewCoffee() *Coffee {
	c := &Coffee{name: "Coffee", recipe: models.NewRecipe()}
	c.recipe.AddIngredient(models.Water, 200)
	c.recipe.AddIngredient(models.CoffeeBeans, 50)
	return c
}

// Name returns the beverage's display name.
func (c *Coffee) Name() string {
	return c.name
}

// Recipe returns the beverage's recipe.
func (c *Coffee) Recipe() *models.Recipe {
	return c.recipe
}

// Prepare narrates the physical brewing step. It never touches ingredient
// inventory; CoffeeMachine already validated and consumed ingredients
// before calling this.
func (c *Coffee) Prepare() {
	fmt.Println("Preparing Coffee...")
	fmt.Println("Brewing coffee with coffee beans...")
}

// Compile-time check that *Coffee satisfies Beverage. Go has no way to
// force a type to implement a method the way some languages enforce
// abstract methods, so this line does that job instead: if Coffee is ever
// missing a method Beverage requires, this line fails to compile with a
// clear error.
var _ Beverage = (*Coffee)(nil)
