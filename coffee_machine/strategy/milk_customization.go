package strategy

import (
	"fmt"

	"golang_low_level_design/coffee_machine/beverages"
	"golang_low_level_design/coffee_machine/models"
)

// MilkCustomization applies a chosen milk type to a beverage.
type MilkCustomization struct {
	milkType models.MilkType
}

// NewMilkCustomization creates a strategy for the given milk type.
func NewMilkCustomization(milkType models.MilkType) *MilkCustomization {
	return &MilkCustomization{milkType: milkType}
}

// Customize reports the milk choice. The beverage parameter is unused
// (fine in Go, unused parameters aren't an error like unused local
// variables are); it exists only to satisfy CustomizationStrategy's
// signature, since customization here doesn't need to modify the recipe.
func (m *MilkCustomization) Customize(beverage beverages.Beverage) {
	fmt.Printf("Using %v milk\n", m.milkType)
}

// Description returns a human-readable summary of this customization.
func (m *MilkCustomization) Description() string {
	return fmt.Sprintf("Milk: %v", m.milkType)
}

// Compile-time check that *MilkCustomization satisfies CustomizationStrategy.
var _ CustomizationStrategy = (*MilkCustomization)(nil)
