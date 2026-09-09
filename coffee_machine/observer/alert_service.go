package observer

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// AlertService is a concrete IngredientObserver that logs a warning
// whenever an ingredient's stock drops below the low-stock threshold.
type AlertService struct{}

// NewAlertService creates an AlertService. It holds no state today, but a
// constructor is kept for consistency with the rest of the codebase, and in
// case a future alert channel (e.g. email, LED) needs configuration here.
func NewAlertService() *AlertService {
	return &AlertService{}
}

// OnIngredientLow logs a low-stock warning for ingredient. It takes
// *models.Ingredient, not models.Ingredient, because Ingredient contains a
// sync.Mutex; copying it by value would copy the lock too (go vet's
// copylocks check would flag exactly this), and it also wouldn't match
// models.IngredientObserver's required signature.
func (a *AlertService) OnIngredientLow(ingredient *models.Ingredient, currentLevel int) {
	fmt.Printf("ALERT: %v is running low, current level is %v\n", ingredient.Name(), currentLevel)
}

// Compile-time check that *AlertService satisfies models.IngredientObserver.
var _ models.IngredientObserver = (*AlertService)(nil)
