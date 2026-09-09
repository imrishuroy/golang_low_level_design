package machine

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// DispensingState is active once a beverage has been prepared and is ready
// to be handed to the customer. It rejects new preparation requests until
// the current beverage is dispensed.
type DispensingState struct{}

// PrepareBeverage is rejected; the current beverage must be dispensed first.
func (s *DispensingState) PrepareBeverage(machine *CoffeeMachine, beverageType models.BeverageType) error {
	fmt.Println("Please dispense current beverage first.")
	return fmt.Errorf("please dispense current beverage first")
}

// Dispense hands over the prepared beverage, clears it, and returns the
// machine to IdleState.
func (s *DispensingState) Dispense(machine *CoffeeMachine) error {
	beverage := machine.CurrentBeverage()
	if beverage == nil {
		return fmt.Errorf("no beverage to dispense")
	}

	fmt.Printf("Dispensing %v...\n", beverage.Name())
	fmt.Printf("%v ready!\n", beverage.Name())
	machine.SetCurrentBeverage(nil)
	machine.ChangeState(&IdleState{})
	return nil
}

// EnterMaintenance is rejected while dispensing.
func (s *DispensingState) EnterMaintenance(machine *CoffeeMachine) {
	fmt.Println("Cannot enter maintenance while dispensing.")
}

// ExitMaintenance is a no-op; the machine isn't in maintenance mode.
func (s *DispensingState) ExitMaintenance(machine *CoffeeMachine) {
	fmt.Println("Not in maintenance mode.")
}

// StateName returns this state's name.
func (s *DispensingState) StateName() string {
	return "DISPENSING"
}

// Compile-time check that *DispensingState satisfies MachineState.
var _ MachineState = (*DispensingState)(nil)
