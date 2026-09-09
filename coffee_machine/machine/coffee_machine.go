package machine

import (
	"fmt"
	"sync"

	"golang_low_level_design/coffee_machine/beverages"
	"golang_low_level_design/coffee_machine/models"
	"golang_low_level_design/coffee_machine/strategy"
)

// CoffeeMachine coordinates beverage preparation: it holds the current
// operational state, the ingredient inventory, any active customizations,
// and the currently-brewing beverage (if any).
//
// This is deliberately not a Singleton: callers construct exactly one with
// NewCoffeeMachine and pass it around explicitly. This keeps the type easy
// to test (each test gets its own isolated machine) and keeps its
// dependencies visible rather than hidden behind a global.
type CoffeeMachine struct {
	// mu guards state, currentBeverage, and customizations, all fields that
	// can be read and mutated from concurrent beverage-preparation requests.
	mu              sync.Mutex
	state           MachineState
	ingredients     map[string]*models.Ingredient
	currentBeverage beverages.Beverage
	customizations  []strategy.CustomizationStrategy
	alertService    models.IngredientObserver
}

// NewCoffeeMachine creates a machine stocked with the given ingredients,
// starting in the Idle state. alertService is registered as an observer on
// every ingredient, so it's notified whenever any of them run low.
//
// The map uses *models.Ingredient (a pointer), not models.Ingredient (a
// value), because Ingredient holds a sync.Mutex internally. Copying an
// Ingredient by value would copy its lock too, leaving two independent
// locks that no longer protect the same data (go vet's copylocks check
// exists specifically to catch this mistake). A pointer means every part
// of the system that touches an ingredient (CoffeeMachine, tests, the
// factory-created beverage recipes) is reading and writing the exact same
// shared instance, which is what makes the thread safety guarantees real.
func NewCoffeeMachine(ingredients map[string]*models.Ingredient, alertService models.IngredientObserver) *CoffeeMachine {
	m := &CoffeeMachine{
		state:        &IdleState{},
		ingredients:  ingredients,
		alertService: alertService,
	}

	for _, ingredient := range m.ingredients {
		ingredient.AddObserver(alertService)
	}

	return m
}

// DefaultIngredients returns the machine's standard starting inventory,
// matching the reference implementation's initializeIngredients(). Keeping
// this separate from NewCoffeeMachine means tests can supply their own
// (e.g. tiny quantities to easily trigger low-stock alerts) instead of
// being stuck with these hardcoded amounts.
func DefaultIngredients() map[string]*models.Ingredient {
	return map[string]*models.Ingredient{
		models.Water:       models.NewIngredient(models.Water, 1000, 2000),
		models.CoffeeBeans: models.NewIngredient(models.CoffeeBeans, 500, 1000),
		models.Milk:        models.NewIngredient(models.Milk, 800, 1500),
		models.TeaLeaves:   models.NewIngredient(models.TeaLeaves, 200, 500),
		models.Sugar:       models.NewIngredient(models.Sugar, 300, 500),
	}
}

// PrepareBeverage requests preparation of the given beverage type. Behavior
// depends entirely on the current state (e.g. Idle accepts it, Preparing
// rejects a second request), per the State Pattern.
func (m *CoffeeMachine) PrepareBeverage(beverageType models.BeverageType) error {
	return m.State().PrepareBeverage(m, beverageType)
}

// Dispense dispenses the currently prepared beverage, if the current state allows it.
func (m *CoffeeMachine) Dispense() error {
	return m.State().Dispense(m)
}

// EnterMaintenance transitions the machine into maintenance mode, if the current state allows it.
func (m *CoffeeMachine) EnterMaintenance() {
	m.State().EnterMaintenance(m)
}

// ExitMaintenance leaves maintenance mode, if the current state allows it.
func (m *CoffeeMachine) ExitMaintenance() {
	m.State().ExitMaintenance(m)
}

// AddCustomization registers a customization to be applied to the next beverage prepared.
func (m *CoffeeMachine) AddCustomization(customization strategy.CustomizationStrategy) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customizations = append(m.customizations, customization)
}

// ClearCustomizations removes all previously added customizations.
func (m *CoffeeMachine) ClearCustomizations() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.customizations = nil
}

// ApplyCustomizations runs every registered customization against beverage.
// It copies the customizations slice under the lock, then runs the
// (potentially slow, user-supplied) Customize calls after unlocking, same
// "don't call code you don't control while holding your own lock"
// discipline we used in Ingredient's notify path.
func (m *CoffeeMachine) ApplyCustomizations(beverage beverages.Beverage) {
	m.mu.Lock()
	customizations := make([]strategy.CustomizationStrategy, len(m.customizations))
	copy(customizations, m.customizations)
	m.mu.Unlock()

	for _, c := range customizations {
		c.Customize(beverage)
	}
}

// ValidateRecipe reports whether every ingredient required by beverageType's
// recipe is currently available in sufficient quantity.
func (m *CoffeeMachine) ValidateRecipe(beverageType models.BeverageType) (bool, error) {
	beverage, err := beverages.NewBeverage(beverageType)
	if err != nil {
		return false, err
	}

	for name, requiredQty := range beverage.Recipe().Ingredients() {
		ingredient, ok := m.ingredients[name]
		if !ok || !ingredient.IsAvailable(requiredQty) {
			return false, nil
		}
	}
	return true, nil
}

// ConsumeIngredients deducts every ingredient in recipe from inventory.
func (m *CoffeeMachine) ConsumeIngredients(recipe *models.Recipe) {
	for name, qty := range recipe.Ingredients() {
		if ingredient, ok := m.ingredients[name]; ok {
			ingredient.Consume(qty)
		}
	}
}

// RefillIngredient adds amount back into the named ingredient's stock.
func (m *CoffeeMachine) RefillIngredient(name string, amount int) {
	ingredient, ok := m.ingredients[name]
	if !ok {
		fmt.Printf("Unknown ingredient: %s\n", name)
		return
	}
	ingredient.Refill(amount)
	fmt.Printf("Refilled %s by %d units.\n", name, amount)
}

// ChangeState transitions the machine to a new state.
func (m *CoffeeMachine) ChangeState(state MachineState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.state = state
}

// State returns the machine's current state.
func (m *CoffeeMachine) State() MachineState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

// SetCurrentBeverage records the beverage currently being prepared/dispensed.
func (m *CoffeeMachine) SetCurrentBeverage(beverage beverages.Beverage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.currentBeverage = beverage
}

// CurrentBeverage returns the beverage currently being prepared/dispensed, if any.
func (m *CoffeeMachine) CurrentBeverage() beverages.Beverage {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.currentBeverage
}

// Ingredients returns a copy of the ingredient inventory (name -> Ingredient
// pointer). Copying the map, but not the *Ingredient values themselves,
// gives callers an independent map, while each *Ingredient is still the
// live, shared instance (copying Ingredient by value would copy its
// internal mutex, which is unsafe).
func (m *CoffeeMachine) Ingredients() map[string]*models.Ingredient {
	out := make(map[string]*models.Ingredient, len(m.ingredients))
	for name, ingredient := range m.ingredients {
		out[name] = ingredient
	}
	return out
}

// DisplayStatus prints the machine's current state and ingredient levels.
func (m *CoffeeMachine) DisplayStatus() {
	fmt.Println("\n=== Coffee Machine Status ===")
	fmt.Println("Current State:", m.State().StateName())
	fmt.Println("\nIngredient Levels:")
	for _, ingredient := range m.ingredients {
		fmt.Printf("  %s: %d\n", ingredient.Name(), ingredient.Quantity())
	}
	fmt.Println("=============================")
}
