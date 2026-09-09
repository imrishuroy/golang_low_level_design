package machine

import "golang_low_level_design/coffee_machine/models"

// MachineState is the heart of the State Pattern used in this project.
// The coffee machine is only allowed to do certain things depending on what
// it is currently doing. For example, it should not start a new drink while
// it is still brewing the last one. Instead of writing a big pile of
// if/else checks everywhere ("if state is X, do this, if state is Y, do
// that"), each state is its own small type that knows exactly what is
// allowed while that state is active. CoffeeMachine simply asks its current
// state to handle the request, and the state decides what happens next
// (including whether to switch to a different state).
type MachineState interface {
	// PrepareBeverage handles a request to start making a drink. Whether
	// this succeeds depends entirely on which state is currently active.
	PrepareBeverage(machine *CoffeeMachine, beverageType models.BeverageType) error

	// Dispense handles a request to move the drink along (for example,
	// actually brewing it, or handing it over once it is ready).
	Dispense(machine *CoffeeMachine) error

	// EnterMaintenance handles a request to take the machine out of service.
	EnterMaintenance(machine *CoffeeMachine)

	// ExitMaintenance handles a request to bring the machine back into service.
	ExitMaintenance(machine *CoffeeMachine)

	// StateName returns a short, readable name for this state, such as
	// "IDLE", used for logging and status messages.
	StateName() string
}
