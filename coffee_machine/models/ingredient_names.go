package models

// Ingredient name constants, used as the shared vocabulary between Recipe
// entries and the CoffeeMachine's Ingredient inventory. Centralizing them
// here means a typo becomes a compile error (undefined identifier) instead
// of a silent runtime mismatch between two different string literals.
const (
	Water       = "Water"
	CoffeeBeans = "CoffeeBeans"
	Milk        = "Milk"
	TeaLeaves   = "TeaLeaves"
	Sugar       = "Sugar"
)
