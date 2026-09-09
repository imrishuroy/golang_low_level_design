package models

import (
	"sync"
	"testing"
)

// fakeObserver records every low-stock notification it receives, so tests
// can assert on how many times (and with what values) OnIngredientLow fired.
type fakeObserver struct {
	mu    sync.Mutex
	calls []int // currentLevel from each call
}

func (f *fakeObserver) OnIngredientLow(ingredient *Ingredient, currentLevel int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, currentLevel)
}

func (f *fakeObserver) callCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.calls)
}

func TestIngredient_ConsumeSuccess(t *testing.T) {
	ing := NewIngredient("Water", 100, 100)

	ok := ing.Consume(30)

	if !ok {
		t.Fatalf("Consume(30) = false, want true")
	}
	if got := ing.Quantity(); got != 70 {
		t.Errorf("Quantity() = %d, want 70", got)
	}
}

func TestIngredient_ConsumeInsufficientStock(t *testing.T) {
	ing := NewIngredient("Water", 10, 100)

	ok := ing.Consume(30)

	if ok {
		t.Fatalf("Consume(30) = true, want false (only 10 in stock)")
	}
	if got := ing.Quantity(); got != 10 {
		t.Errorf("Quantity() = %d, want 10 (unchanged after failed consume)", got)
	}
}

func TestIngredient_RefillCapsAtMaxCapacity(t *testing.T) {
	ing := NewIngredient("Water", 90, 100)

	ing.Refill(50) // would be 140, must cap at 100

	if got := ing.Quantity(); got != 100 {
		t.Errorf("Quantity() = %d, want 100 (capped at maxCapacity)", got)
	}
}

func TestIngredient_IsAvailable(t *testing.T) {
	ing := NewIngredient("Water", 50, 100)

	tests := []struct {
		name     string
		required int
		want     bool
	}{
		{"less than stock", 10, true},
		{"exactly equal to stock", 50, true},
		{"more than stock", 51, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ing.IsAvailable(tt.required); got != tt.want {
				t.Errorf("IsAvailable(%d) = %v, want %v", tt.required, got, tt.want)
			}
		})
	}
}

// TestIngredient_NotifiesExactlyOnceWhenCrossingThreshold verifies the
// Observer Pattern requirement: observers fire only at the moment stock
// crosses from at-or-above the 20% threshold to below it, not on every
// Consume call.
func TestIngredient_NotifiesExactlyOnceWhenCrossingThreshold(t *testing.T) {
	ing := NewIngredient("Water", 100, 100) // 100% full
	obs := &fakeObserver{}
	ing.AddObserver(obs)

	ing.Consume(30) // 100 -> 70, still above 20%, no notification expected
	if got := obs.callCount(); got != 0 {
		t.Fatalf("after dropping to 70%%: callCount = %d, want 0", got)
	}

	ing.Consume(55) // 70 -> 15, crosses below 20%, exactly one notification expected
	if got := obs.callCount(); got != 1 {
		t.Fatalf("after dropping to 15%%: callCount = %d, want 1", got)
	}

	ing.Consume(5) // 15 -> 10, still below 20%, no additional notification expected
	if got := obs.callCount(); got != 1 {
		t.Fatalf("after dropping to 10%%: callCount = %d, want still 1 (no re-notify)", got)
	}
}

// TestIngredient_RemoveObserverStopsNotifications verifies RemoveObserver
// actually detaches the observer.
func TestIngredient_RemoveObserverStopsNotifications(t *testing.T) {
	ing := NewIngredient("Water", 100, 100)
	obs := &fakeObserver{}
	ing.AddObserver(obs)
	ing.RemoveObserver(obs)

	ing.Consume(90) // would cross the threshold if obs were still registered

	if got := obs.callCount(); got != 0 {
		t.Errorf("callCount = %d, want 0 (observer was removed)", got)
	}
}

// TestIngredient_ConcurrentConsume is the key thread-safety test: many
// goroutines race to Consume from the same Ingredient at once. Run with
// `go test -race` to have Go's race detector verify there's no unsynchronized
// access, and the assertion below verifies no more was consumed in total
// than was actually in stock (the classic "two requests both see enough
// milk, together consume too much" race the assignment warns about).
func TestIngredient_ConcurrentConsume(t *testing.T) {
	const (
		initial    = 1000
		numWorkers = 100
		perWorker  = 20 // 100 * 20 = 2000 total attempted, only 1000 can succeed
	)
	ing := NewIngredient("Water", initial, initial)

	var wg sync.WaitGroup
	var succeeded int32
	var mu sync.Mutex

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ing.Consume(perWorker) {
				mu.Lock()
				succeeded++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	wantRemaining := initial - int(succeeded)*perWorker
	if got := ing.Quantity(); got != wantRemaining {
		t.Errorf("Quantity() = %d, want %d (initial %d minus %d successful consumes of %d each)",
			got, wantRemaining, initial, succeeded, perWorker)
	}
	if got := ing.Quantity(); got < 0 {
		t.Errorf("Quantity() went negative: %d, this means concurrent Consume calls over-consumed stock", got)
	}
}
