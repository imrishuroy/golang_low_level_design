package models

// BeverageType lists every drink the coffee machine knows how to make.
// It is used instead of a plain string so that typos are caught by the
// compiler: passing an invalid value like models.BeverageType(99) still
// works technically, but passing a random string like "coffe" would not
// even compile, since Go checks the type at build time.
type BeverageType int

const (
	Coffee BeverageType = iota
	Tea
	Cappuccino
	Latte
)

// String returns a readable name for the beverage type, such as "Coffee".
// It exists because a plain BeverageType value only stores a number (0, 1,
// 2, 3) inside the computer; this method turns that number into text that
// makes sense to a person reading logs or messages. Go's fmt package
// automatically calls this method whenever a BeverageType is printed.
func (b BeverageType) String() string {
	switch b {
	case Coffee:
		return "Coffee"
	case Tea:
		return "Tea"
	case Cappuccino:
		return "Cappuccino"
	case Latte:
		return "Latte"
	default:
		return "Unknown"
	}
}
