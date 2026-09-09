package models

import "sync"

// IngredientObserver is notified when an ingredient's level drops below its
// low-stock threshold. It's declared here, next to Ingredient, because
// Ingredient is the one that calls it (the consumer of the interface).
// If this lived in a separate package instead, that package would need to
// import models (for *Ingredient), while models would need to import it
// back for this interface, an import cycle Go refuses to compile. Declaring
// the interface where it's consumed avoids that entirely.
type IngredientObserver interface {
	OnIngredientLow(ingredient *Ingredient, currentLevel int)
}

// lowStockThresholdPercent is the level, as a percentage of maxCapacity,
// below which an ingredient is considered low and observers are notified.
const lowStockThresholdPercent = 20.0

// Ingredient tracks the live quantity of a single stocked item (beans, milk,
// water, sugar). All access goes through mu, since multiple beverage
// preparations can run concurrently and must not race on quantity.
type Ingredient struct {
	mu          sync.Mutex
	name        string
	quantity    int
	maxCapacity int
	observers   []IngredientObserver
}

// NewIngredient creates an ingredient with the given starting quantity and
// maximum capacity.
func NewIngredient(name string, initialQuantity, maxCapacity int) *Ingredient {
	return &Ingredient{
		name:        name,
		quantity:    initialQuantity,
		maxCapacity: maxCapacity,
	}
}

// Name returns the ingredient's name. It's immutable after construction, so
// it's safe to read without locking.
func (i *Ingredient) Name() string {
	return i.name
}

// Quantity returns the current quantity in stock.
func (i *Ingredient) Quantity() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.quantity
}

// Consume attempts to deduct amount from stock, reporting whether there was
// enough available. The check and the deduction happen atomically under the
// same lock, so two concurrent Consume calls can't both succeed against
// stock that's only sufficient for one of them (the race the assignment's
// thread-safety requirement is about).
func (i *Ingredient) Consume(amount int) bool {
	i.mu.Lock()
	if i.quantity < amount {
		i.mu.Unlock()
		return false
	}

	oldQuantity := i.quantity
	i.quantity -= amount
	crossed, level, obs := i.crossedThresholdLocked(oldQuantity)
	i.mu.Unlock()

	// Notify only after unlocking. sync.Mutex is not reentrant, so calling
	// observer code (code this package doesn't control) while still holding
	// the lock risks a deadlock if an observer ever calls back into this
	// Ingredient. It's also just safer practice: never run a callback you
	// don't control while holding your own lock.
	if crossed {
		notify(obs, i, level)
	}
	return true
}

// Refill adds amount back to stock, capped at maxCapacity.
func (i *Ingredient) Refill(amount int) {
	i.mu.Lock()
	oldQuantity := i.quantity
	i.quantity = min(i.quantity+amount, i.maxCapacity)
	crossed, level, obs := i.crossedThresholdLocked(oldQuantity)
	i.mu.Unlock()

	if crossed {
		notify(obs, i, level)
	}
}

// crossedThresholdLocked reports whether quantity just dropped from at-or-
// above the low-stock threshold to below it, given the quantity before this
// change. Callers must hold mu when calling this. It returns a snapshot
// copy of the observer list so notify can run after mu is released.
func (i *Ingredient) crossedThresholdLocked(oldQuantity int) (crossed bool, level int, obs []IngredientObserver) {
	oldPercentage := float64(oldQuantity) * 100.0 / float64(i.maxCapacity)
	newPercentage := float64(i.quantity) * 100.0 / float64(i.maxCapacity)

	crossed = oldPercentage > lowStockThresholdPercent && newPercentage <= lowStockThresholdPercent
	if !crossed {
		return false, 0, nil
	}

	obsCopy := make([]IngredientObserver, len(i.observers))
	copy(obsCopy, i.observers)
	return true, i.quantity, obsCopy
}

// notify calls OnIngredientLow on each observer. It's a free function, not a
// method, precisely so it never has access to i.mu and can't accidentally
// be called while the lock is held.
func notify(observers []IngredientObserver, ingredient *Ingredient, level int) {
	for _, o := range observers {
		o.OnIngredientLow(ingredient, level)
	}
}

// AddObserver registers an observer to be notified on low stock.
func (i *Ingredient) AddObserver(o IngredientObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.observers = append(i.observers, o)
}

// RemoveObserver unregisters a previously added observer.
func (i *Ingredient) RemoveObserver(o IngredientObserver) {
	i.mu.Lock()
	defer i.mu.Unlock()
	for idx, existing := range i.observers {
		if existing == o {
			i.observers = append(i.observers[:idx], i.observers[idx+1:]...)
			return
		}
	}
}

// IsAvailable reports whether at least requiredAmount is currently in stock.
func (i *Ingredient) IsAvailable(requiredAmount int) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.quantity >= requiredAmount
}
