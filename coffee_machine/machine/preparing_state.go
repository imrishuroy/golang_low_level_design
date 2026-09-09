package machine

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// PreparingState is active while a beverage is being brewed. It rejects new
// preparation requests and maintenance transitions until preparation finishes.
type PreparingState struct{}

// PrepareBeverage is invalid while already preparing; the request is rejected.
func (s *PreparingState) PrepareBeverage(machine *CoffeeMachine, beverageType models.BeverageType) error {
	fmt.Println("Already preparing a beverage. Please wait...")
	return fmt.Errorf("already preparing a beverage")
}

// Dispense runs the beverage's actual preparation, consumes its ingredients,
// and transitions to DispensingState.
func (s *PreparingState) Dispense(machine *CoffeeMachine) error {
	beverage := machine.CurrentBeverage()
	if beverage == nil {
		return fmt.Errorf("no beverage to prepare")
	}
	beverage.Prepare()
	machine.ConsumeIngredients(beverage.Recipe())
	machine.ChangeState(&DispensingState{})
	fmt.Println("Beverage prepared. Ready to dispense.")
	return nil
}

// EnterMaintenance is rejected while a beverage is being prepared.
func (s *PreparingState) EnterMaintenance(machine *CoffeeMachine) {
	fmt.Println("Cannot enter maintenance while preparing beverage.")
}

// ExitMaintenance is a no-op; the machine isn't in maintenance mode.
func (s *PreparingState) ExitMaintenance(machine *CoffeeMachine) {
	fmt.Println("Not in maintenance mode.")
}

// StateName returns this state's name.
func (s *PreparingState) StateName() string {
	return "PREPARING"
}

// Compile-time check that *PreparingState satisfies MachineState.
var _ MachineState = (*PreparingState)(nil)
