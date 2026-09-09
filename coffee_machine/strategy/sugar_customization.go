package strategy

import (
	"fmt"

	"golang_low_level_design/coffee_machine/beverages"
	"golang_low_level_design/coffee_machine/models"
)

// SugarCustomization applies a chosen sugar level to a beverage.
type SugarCustomization struct {
	sugarLevel models.SugarLevel
}

// NewSugarCustomization creates a strategy for the given sugar level.
func NewSugarCustomization(sugarLevel models.SugarLevel) *SugarCustomization {
	return &SugarCustomization{sugarLevel: sugarLevel}
}

// Customize reports the sugar level choice. The beverage parameter is
// unused; it exists only to satisfy CustomizationStrategy's signature.
func (s *SugarCustomization) Customize(beverage beverages.Beverage) {
	fmt.Printf("Adding %v sugar level\n", s.sugarLevel)
}

// Description returns a human-readable summary of this customization.
func (s *SugarCustomization) Description() string {
	return fmt.Sprintf("Sugar: %v", s.sugarLevel)
}

// Compile-time check that *SugarCustomization satisfies CustomizationStrategy.
var _ CustomizationStrategy = (*SugarCustomization)(nil)
