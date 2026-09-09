package machine_test

import (
	"testing"

	"golang_low_level_design/coffee_machine/beverages"
	"golang_low_level_design/coffee_machine/machine"
	"golang_low_level_design/coffee_machine/models"
)

// noopObserver satisfies models.IngredientObserver without doing anything;
// tests don't care about alerts, only about state/inventory behavior.
type noopObserver struct{}

func (noopObserver) OnIngredientLow(ingredient *models.Ingredient, currentLevel int) {}

// fakeStrategy is a test double for strategy.CustomizationStrategy that
// records whether (and with what beverage) it was invoked, so tests can
// assert ApplyCustomizations actually calls through to it.
type fakeStrategy struct {
	called bool
	got    beverages.Beverage
}

func (f *fakeStrategy) Customize(b beverages.Beverage) {
	f.called = true
	f.got = b
}

func (f *fakeStrategy) Description() string { return "fake customization" }

func TestCoffeeMachine_PrepareBeverage_Success(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})

	err := m.PrepareBeverage(models.Coffee)

	if err != nil {
		t.Fatalf("PrepareBeverage(Coffee) error = %v, want nil", err)
	}
	if got := m.State().StateName(); got != "PREPARING" {
		t.Errorf("State().StateName() = %q, want %q", got, "PREPARING")
	}
}

func TestCoffeeMachine_PrepareBeverage_InsufficientIngredients(t *testing.T) {
	tinyIngredients := map[string]*models.Ingredient{
		models.Water:       models.NewIngredient(models.Water, 1, 100),
		models.CoffeeBeans: models.NewIngredient(models.CoffeeBeans, 1, 100),
	}
	m := machine.NewCoffeeMachine(tinyIngredients, noopObserver{})

	err := m.PrepareBeverage(models.Coffee)

	if err == nil {
		t.Fatal("PrepareBeverage(Coffee) error = nil, want an error (insufficient ingredients)")
	}
	if got := m.State().StateName(); got != "IDLE" {
		t.Errorf("State().StateName() = %q, want %q (should stay Idle after rejection)", got, "IDLE")
	}
}

func TestCoffeeMachine_PrepareBeverage_RejectedWhileAlreadyPreparing(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})

	if err := m.PrepareBeverage(models.Coffee); err != nil {
		t.Fatalf("first PrepareBeverage(Coffee) error = %v, want nil", err)
	}

	err := m.PrepareBeverage(models.Tea)

	if err == nil {
		t.Fatal("second PrepareBeverage(Tea) error = nil, want an error (already preparing)")
	}
	if got := m.State().StateName(); got != "PREPARING" {
		t.Errorf("State().StateName() = %q, want %q (unchanged)", got, "PREPARING")
	}
}

func TestCoffeeMachine_FullDispenseCycleConsumesIngredients(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})

	if err := m.PrepareBeverage(models.Coffee); err != nil {
		t.Fatalf("PrepareBeverage(Coffee) error = %v, want nil", err)
	}

	// First Dispense(): PreparingState actually runs Prepare() and consumes
	// ingredients, then transitions to DispensingState.
	if err := m.Dispense(); err != nil {
		t.Fatalf("first Dispense() error = %v, want nil", err)
	}
	if got := m.State().StateName(); got != "DISPENSING" {
		t.Errorf("State().StateName() = %q, want %q", got, "DISPENSING")
	}

	water := m.Ingredients()[models.Water]
	if got, want := water.Quantity(), 1000-200; got != want {
		t.Errorf("Water.Quantity() after first Dispense = %d, want %d", got, want)
	}

	// Second Dispense(): DispensingState hands over the beverage and
	// returns to IdleState.
	if err := m.Dispense(); err != nil {
		t.Fatalf("second Dispense() error = %v, want nil", err)
	}
	if got := m.State().StateName(); got != "IDLE" {
		t.Errorf("State().StateName() = %q, want %q", got, "IDLE")
	}
	if m.CurrentBeverage() != nil {
		t.Errorf("CurrentBeverage() = %v, want nil after dispensing completes", m.CurrentBeverage())
	}
}

func TestCoffeeMachine_MaintenanceBlocksPreparationUntilExited(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})

	m.EnterMaintenance()
	if got := m.State().StateName(); got != "MAINTENANCE" {
		t.Fatalf("State().StateName() = %q, want %q", got, "MAINTENANCE")
	}

	if err := m.PrepareBeverage(models.Coffee); err == nil {
		t.Error("PrepareBeverage(Coffee) during maintenance error = nil, want an error")
	}

	m.ExitMaintenance()
	if got := m.State().StateName(); got != "IDLE" {
		t.Fatalf("State().StateName() = %q, want %q", got, "IDLE")
	}

	if err := m.PrepareBeverage(models.Coffee); err != nil {
		t.Errorf("PrepareBeverage(Coffee) after exiting maintenance error = %v, want nil", err)
	}
}

func TestCoffeeMachine_ApplyCustomizationsCallsRegisteredStrategies(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})
	fake := &fakeStrategy{}
	m.AddCustomization(fake)

	if err := m.PrepareBeverage(models.Tea); err != nil {
		t.Fatalf("PrepareBeverage(Tea) error = %v, want nil", err)
	}
	beverage := m.CurrentBeverage()
	if beverage == nil {
		t.Fatal("CurrentBeverage() = nil, want the beverage just prepared")
	}

	m.ApplyCustomizations(beverage)

	if !fake.called {
		t.Error("fakeStrategy.Customize was never called; ApplyCustomizations should invoke every registered strategy")
	}
	if fake.got != beverage {
		t.Errorf("fakeStrategy received %v, want the same beverage instance %v", fake.got, beverage)
	}
}

func TestCoffeeMachine_ClearCustomizationsRemovesStrategies(t *testing.T) {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), noopObserver{})
	fake := &fakeStrategy{}
	m.AddCustomization(fake)
	m.ClearCustomizations()

	if err := m.PrepareBeverage(models.Coffee); err != nil {
		t.Fatalf("PrepareBeverage(Coffee) error = %v, want nil", err)
	}
	m.ApplyCustomizations(m.CurrentBeverage())

	if fake.called {
		t.Error("fakeStrategy.Customize was called, but it should have been cleared before ApplyCustomizations")
	}
}
