package machine

import (
	"fmt"

	"golang_low_level_design/coffee_machine/models"
)

// MaintenanceState takes the machine out of service; no beverage operations
// are allowed until maintenance is exited.
type MaintenanceState struct{}

// PrepareBeverage is rejected while in maintenance mode.
func (s *MaintenanceState) PrepareBeverage(machine *CoffeeMachine, beverageType models.BeverageType) error {
	fmt.Println("Machine is in maintenance mode. Cannot prepare beverages.")
	return fmt.Errorf("machine is in maintenance mode")
}

// Dispense is rejected while in maintenance mode.
func (s *MaintenanceState) Dispense(machine *CoffeeMachine) error {
	fmt.Println("Machine is in maintenance mode. Cannot dispense beverages.")
	return fmt.Errorf("machine is in maintenance mode")
}

// EnterMaintenance is a no-op; the machine is already in maintenance mode.
func (s *MaintenanceState) EnterMaintenance(machine *CoffeeMachine) {
	fmt.Println("Already in maintenance mode.")
}

// ExitMaintenance returns the machine to IdleState.
func (s *MaintenanceState) ExitMaintenance(machine *CoffeeMachine) {
	machine.ChangeState(&IdleState{})
	fmt.Println("Exiting maintenance mode. Machine ready for operation.")
}

// StateName returns this state's name.
func (s *MaintenanceState) StateName() string {
	return "MAINTENANCE"
}

// Compile-time check that *MaintenanceState satisfies MachineState.
var _ MachineState = (*MaintenanceState)(nil)
