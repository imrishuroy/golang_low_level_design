package models

// MilkType is the milk variant a customer can choose for a beverage.
type MilkType int

const (
	Regular MilkType = iota
	Skim
	Almond
)

// String returns MilkType's human-readable name, satisfying fmt.Stringer.
func (m MilkType) String() string {
	switch m {
	case Regular:
		return "Regular"
	case Skim:
		return "Skim"
	case Almond:
		return "Almond"
	default:
		return "Unknown"
	}
}
