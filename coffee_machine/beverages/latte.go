package beverages

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// Latte is espresso with a larger portion of steamed milk than a cappuccino.
type Latte struct {
	name   string
	recipe *models.Recipe
}

// NewLatte creates a Latte beverage with its recipe already populated.
func NewLatte() *Latte {
	l := &Latte{name: "Latte", recipe: models.NewRecipe()}
	l.recipe.AddIngredient(models.Water, 150)
	l.recipe.AddIngredient(models.CoffeeBeans, 50)
	l.recipe.AddIngredient(models.Milk, 150)
	return l
}

// Name returns the beverage's display name.
func (l *Latte) Name() string {
	return l.name
}

// Recipe returns the beverage's recipe.
func (l *Latte) Recipe() *models.Recipe {
	return l.recipe
}

// Prepare narrates the physical brewing step.
func (l *Latte) Prepare() {
	fmt.Println("Preparing Latte...")
	fmt.Println("Brewing espresso and steaming milk...")
}

var _ Beverage = (*Latte)(nil)
