package beverages

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// Tea is steeped tea leaves in hot water.
type Tea struct {
	name   string
	recipe *models.Recipe
}

// NewTea creates a Tea beverage with its recipe already populated.
func NewTea() *Tea {
	t := &Tea{name: "Tea", recipe: models.NewRecipe()}
	t.recipe.AddIngredient(models.Water, 200)
	t.recipe.AddIngredient(models.TeaLeaves, 30)
	return t
}

// Name returns the beverage's display name.
func (t *Tea) Name() string {
	return t.name
}

// Recipe returns the beverage's recipe.
func (t *Tea) Recipe() *models.Recipe {
	return t.recipe
}

// Prepare narrates the physical brewing step.
func (t *Tea) Prepare() {
	fmt.Println("Preparing Tea...")
	fmt.Println("Steeping tea leaves...")
}

var _ Beverage = (*Tea)(nil)
