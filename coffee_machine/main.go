package main

import (
	"fmt"

	"golang_low_level_design/coffee_machine/machine"
	"golang_low_level_design/coffee_machine/models"
	"golang_low_level_design/coffee_machine/observer"
	"golang_low_level_design/coffee_machine/strategy"
)

func main() {
	m := machine.NewCoffeeMachine(machine.DefaultIngredients(), observer.NewAlertService())

	fmt.Println("=== Coffee Machine Demo ===")
	m.DisplayStatus()

	fmt.Println("--- Preparing Coffee ---")
	if err := m.PrepareBeverage(models.Coffee); err != nil {
		fmt.Println("Error:", err)
	}
	m.Dispense()
	m.Dispense()

	fmt.Println("\n--- Preparing Tea with Customization ---")
	m.AddCustomization(strategy.NewSugarCustomization(models.Medium))
	if err := m.PrepareBeverage(models.Tea); err != nil {
		fmt.Println("Error:", err)
	}
	if tea := m.CurrentBeverage(); tea != nil {
		m.ApplyCustomizations(tea)
	}
	m.Dispense()
	m.Dispense()
	m.ClearCustomizations()

	fmt.Println("\n--- Preparing Cappuccino ---")
	if err := m.PrepareBeverage(models.Cappuccino); err != nil {
		fmt.Println("Error:", err)
	}
	m.Dispense()
	m.Dispense()

	fmt.Println("\n--- Preparing Latte with Customization ---")
	m.AddCustomization(strategy.NewSugarCustomization(models.High))
	m.AddCustomization(strategy.NewMilkCustomization(models.Almond))
	if err := m.PrepareBeverage(models.Latte); err != nil {
		fmt.Println("Error:", err)
	}
	if latte := m.CurrentBeverage(); latte != nil {
		m.ApplyCustomizations(latte)
	}
	m.Dispense()
	m.Dispense()
	m.ClearCustomizations()

	m.DisplayStatus()

	fmt.Println("--- Testing Maintenance Mode ---")
	m.EnterMaintenance()
	if err := m.PrepareBeverage(models.Coffee); err != nil {
		fmt.Println("Error:", err)
	}
	m.ExitMaintenance()

	fmt.Println("\n--- Refilling Ingredients ---")
	m.RefillIngredient(models.CoffeeBeans, 200)
	m.RefillIngredient(models.Milk, 300)

	m.DisplayStatus()
}
