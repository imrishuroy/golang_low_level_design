package models

// SugarLevel is the amount of sugar a customer can choose for a beverage.
type SugarLevel int

const (
	None SugarLevel = iota
	Low
	Medium
	High
)

// String returns SugarLevel's human-readable name, satisfying fmt.Stringer.
func (s SugarLevel) String() string {
	switch s {
	case None:
		return "None"
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Unknown"
	}
}
