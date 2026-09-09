package machine

import (
	"fmt"

	"golang_low_level_design/coffee_machine/beverages"
	"golang_low_level_design/coffee_machine/models"
)

// IdleState is the machine's default, ready-to-accept-requests state.
type IdleState struct{}

// PrepareBeverage validates the recipe and, if enough ingredients are
// available, creates the beverage and transitions to PreparingState.
// Otherwise the request is rejected and the machine stays Idle.
func (s *IdleState) PrepareBeverage(machine *CoffeeMachine, beverageType models.BeverageType) error {
	ok, err := machine.ValidateRecipe(beverageType)
	if err != nil {
		return err
	}

	if !ok {
		fmt.Printf("Cannot prepare %v: Insufficient ingredients.\n", beverageType)
		return fmt.Errorf("cannot prepare %v: insufficient ingredients", beverageType)
	}
	beverage, err := beverages.NewBeverage(beverageType)
	if err != nil {
		return err
	}
	machine.SetCurrentBeverage(beverage)
	machine.ChangeState(&PreparingState{})
	fmt.Println("Starting beverage preparation...")
	return nil
}

// Dispense is invalid while Idle; there's nothing prepared yet.
func (s *IdleState) Dispense(machine *CoffeeMachine) error {
	fmt.Println("No beverage ready to dispense.")
	return fmt.Errorf("no beverage ready to dispense")
}

// EnterMaintenance transitions the machine into MaintenanceState.
func (s *IdleState) EnterMaintenance(machine *CoffeeMachine) {
	machine.ChangeState(&MaintenanceState{})
	fmt.Println("Entering maintenance mode...")
}

// ExitMaintenance is a no-op; the machine is already in normal operation.
func (s *IdleState) ExitMaintenance(machine *CoffeeMachine) {
	fmt.Println("Already in normal operation mode.")
}

// StateName returns this state's name.
func (s *IdleState) StateName() string {
	return "IDLE"
}

// Compile-time check that *IdleState satisfies MachineState.
var _ MachineState = (*IdleState)(nil)
