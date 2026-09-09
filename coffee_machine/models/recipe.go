package models

// Recipe defines the ingredients and quantities required to prepare a
// beverage. It's a plain data holder; validation against actual stock
// happens elsewhere (in the Ingredient/inventory layer), keeping this type
// single purpose (SRP).
type Recipe struct {
	// ingredients maps ingredient name to required quantity.
	// Lowercase field means unexported: only this package (models) can touch
	// it directly. Everyone else must go through the methods below. This is
	// how Go does encapsulation, no "private" keyword, just lowercase names.
	ingredients map[string]int
}

// NewRecipe creates an empty recipe with no ingredients.
// Go has no constructors, so "New<Type>" returning a pointer is the standard
// convention (see e.g. bytes.NewBuffer, time.NewTimer in the stdlib).
func NewRecipe() *Recipe {
	// Struct literal, must initialize the map: a nil map can be read from
	// (returns zero value) but panics if you try to write to it. Since
	// AddIngredient below writes to this map, it must be non-nil here.
	return &Recipe{ingredients: make(map[string]int)}
}

// NewRecipeFromIngredients creates a recipe pre-populated from the given
// ingredients. The map is copied so later changes by the caller don't affect
// the recipe.
func NewRecipeFromIngredients(ingredients map[string]int) *Recipe {
	r := NewRecipe()

	// Copy entry-by-entry rather than doing r.ingredients = ingredients.
	// Maps are reference types in Go, so a plain assignment would make r
	// share the exact same backing map as the caller. If the caller later
	// mutated their map, this recipe would silently change too. Copying
	// breaks that aliasing and keeps the recipe's data independent.
	for name, quantity := range ingredients {
		r.ingredients[name] = quantity
	}

	return r
}

// AddIngredient sets the required quantity for an ingredient in this recipe.
func (r *Recipe) AddIngredient(name string, quantity int) {
	r.ingredients[name] = quantity
}

// Ingredients returns a copy of the recipe's required ingredients so callers
// can't mutate the recipe's internal state through the returned map.
func (r *Recipe) Ingredients() map[string]int {
	out := make(map[string]int, len(r.ingredients))
	for name, quantity := range r.ingredients {
		out[name] = quantity
	}
	return out
}

// RequiredQuantity returns the quantity required for the given ingredient,
// or 0 if the recipe doesn't use it.
//
// Note: map lookups in Go return the zero value for a missing key instead of
// throwing, so r.ingredients[name] is safe even if name was never added.
func (r *Recipe) RequiredQuantity(name string) int {
	return r.ingredients[name]
}
