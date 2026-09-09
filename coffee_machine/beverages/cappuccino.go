package beverages

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// Cappuccino is espresso topped with steamed, frothed milk.
type Cappuccino struct {
	name   string
	recipe *models.Recipe
}

// NewCappuccino creates a Cappuccino beverage with its recipe already populated.
func NewCappuccino() *Cappuccino {
	c := &Cappuccino{name: "Cappuccino", recipe: models.NewRecipe()}
	c.recipe.AddIngredient(models.Water, 150)
	c.recipe.AddIngredient(models.CoffeeBeans, 50)
	c.recipe.AddIngredient(models.Milk, 100)
	return c
}

// Name returns the beverage's display name.
func (c *Cappuccino) Name() string {
	return c.name
}

// Recipe returns the beverage's recipe.
func (c *Cappuccino) Recipe() *models.Recipe {
	return c.recipe
}

// Prepare narrates the physical brewing step.
func (c *Cappuccino) Prepare() {
	fmt.Println("Preparing Cappuccino...")
	fmt.Println("Brewing espresso and frothing milk...")
}

var _ Beverage = (*Cappuccino)(nil)
