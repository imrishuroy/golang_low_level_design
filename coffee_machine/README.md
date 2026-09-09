# Coffee Machine, Low Level Design in Go

This project is a learning exercise for practicing Low Level Design (LLD)
in Go. It designs a coffee machine system that prepares different drinks,
tracks ingredient stock, supports customization, and manages its own
operating state, using four classic design patterns: Factory, Strategy,
State, and Observer.

The original problem statement is kept in [doc.md](doc.md). This README
explains what was built, why it was built that way, and the Go specific
lessons that came up while building it.

## What the system does

The coffee machine can:

- Make four drinks: Coffee, Tea, Cappuccino, Latte.
- Keep track of ingredient stock (water, coffee beans, milk, tea leaves, sugar).
- Check a drink's recipe against current stock before starting, and refuse
  to start if something is missing.
- Let a customer add customizations, like a sugar level or a milk type.
- Warn someone (through a simple alert message) when an ingredient runs low.
- Move through clear stages while working: idle, preparing, dispensing, and
  a separate maintenance mode.
- Support refilling stock.
- Do all of this safely even if many requests come in at the same time.

## How to run and test it

All commands below are run from the repository root (the folder that
contains `go.mod`).

### Run the demo

```
go run ./coffee_machine
```

This runs a full demo: it prepares several drinks (some with
customizations), shows a low stock alert firing, tests maintenance mode,
and refills ingredients.

### Run all tests

```
go test ./coffee_machine/...
```

The `...` means "this folder and every folder under it", so this runs the
tests in every package at once. Packages with no test files (like
`beverages`) just get reported as `[no test files]`, which is expected, not
an error.

### Run all tests with the race detector (the important one)

```
go test -race ./coffee_machine/...
```

`-race` turns on Go's race detector, a tool that checks, while the program
actually runs, whether two goroutines (Go's lightweight threads) touched the
same piece of memory at the same time without protection. This is the main
way this project proves its thread safety claims are real, not just
"looks right on paper". Always use `-race` when testing anything that
touches `Ingredient` or `CoffeeMachine`.

### Run tests with detailed output

```
go test -race -v ./coffee_machine/...
```

`-v` prints every individual test name as it runs (`=== RUN` / `--- PASS`),
instead of just one summary line per package. Useful when a test fails and
you need to see exactly which one.

### Run just one package's tests

```
go test -race ./coffee_machine/models/...
go test -race ./coffee_machine/machine/...
```

### Run just one specific test by name

```
go test -race -v ./coffee_machine/models/... -run TestIngredient_ConcurrentConsume
```

`-run` takes a regular expression matched against test function names, so
this only runs (and shows output for) that one test instead of the whole
package.

### Check the code compiles, without running it

```
go build ./coffee_machine/...
```

### Check for common mistakes (including the ones this project ran into)

```
go vet ./coffee_machine/...
```

`go vet` looks for things that compile fine but are almost certainly bugs,
like passing a struct containing a `sync.Mutex` by value instead of by
pointer (the `copylocks` check), which is exactly the kind of mistake this
project's own `models.IngredientObserver` interface was designed to avoid.

### Check formatting

```
gofmt -l coffee_machine/
```

This lists any file that is not formatted the standard Go way. No output
means everything is clean. To fix formatting automatically instead of just
checking it:

```
gofmt -w coffee_machine/
```

### One command that does most of the above at once

```
gofmt -l coffee_machine/ && go vet ./coffee_machine/... && go test -race ./coffee_machine/...
```

Because of `&&`, each command only runs if the one before it succeeded, so
this stops early and shows you the first real problem instead of a wall of
unrelated output.

## Project structure

```
coffee_machine/
  models/     Plain data: Recipe, Ingredient, and simple enum-like types
  beverages/  The drinks themselves, and the factory that creates them
  strategy/   Customer customizations (sugar level, milk type)
  observer/   The alert that fires when stock runs low
  machine/    The machine itself, and the states it can be in
  main.go     A demo program that exercises the whole system
```

This layout is organized by what each part of the system *does*
(`beverages`, `machine`, `strategy`), not by which design pattern it
happens to use. There is a practical reason for this, explained in the
State Pattern section below: Go does not allow two packages to import each
other in a circle, so types that depend on each other closely (like
`CoffeeMachine` and `MachineState`) have to live in the same package.

## The four design patterns used

### 1. Factory Pattern (in `beverages/factory.go`)

**In plain words:** a factory is just a function whose only job is to
create the right object for a given input, so nobody else in the program
has to know how to build it themselves.

**Why this problem needs it:** the system supports four drinks today, and
the requirements say it should be easy to add more later. Without a
factory, every part of the code that needs to create a drink would need its
own copy of "if type is Coffee, make a Coffee, if type is Tea, make a
Tea...". That logic would need to be repeated and kept in sync everywhere.

**Where it lives:** `NewBeverage(beverageType)` in `beverages/factory.go`.
It takes a `BeverageType` value and returns something that satisfies the
`Beverage` interface:

```go
func NewBeverage(beverageType models.BeverageType) (Beverage, error) {
	switch beverageType {
	case models.Coffee:
		return NewCoffee(), nil
	case models.Cappuccino:
		return NewCappuccino(), nil
	// ...
	default:
		return nil, fmt.Errorf("unsupported beverage type: %v", beverageType)
	}
}
```

Adding a fifth drink later means writing one new file (like `mocha.go`) and
adding one line to this switch statement. Nothing else in the system needs
to change.

### 2. Strategy Pattern (in `strategy/`)

**In plain words:** instead of writing one big function full of "if the
customer wants this option, do that; if they want this other option, do
this other thing", each option becomes its own small, separate type. All
of these types share one common interface, so the rest of the program can
treat them the same way without caring which specific one it is dealing
with.

**Why this problem needs it:** customers can pick a sugar level and a milk
type, and the requirements ask for this to be easy to extend later (more
sugar levels, more milk types, maybe entirely new kinds of customization).

**Where it lives:** `CustomizationStrategy` in
`strategy/customization_strategy.go` is the shared interface:

```go
type CustomizationStrategy interface {
	Customize(beverage beverages.Beverage)
	Description() string
}
```

`MilkCustomization` and `SugarCustomization` each implement this interface
in their own way. `CoffeeMachine` just keeps a list of whichever
customizations were requested, and calls `Customize` on each one, without
needing to know or care what is inside any of them.

### 3. State Pattern (in `machine/`)

**In plain words:** some objects behave completely differently depending on
what "mode" they are currently in. Instead of writing a giant pile of
if/else checks everywhere ("if the machine is idle, allow this; if it is
busy, refuse it"), each mode becomes its own small type. The object simply
asks whichever mode is currently active to handle the request, and that
mode decides what happens, including whether to switch to a different mode.

**Why this problem needs it:** the coffee machine should not accept a new
drink request while it is still preparing the last one, and it should
refuse almost everything while it is in maintenance mode. There are four
modes in total: Idle, Preparing, Dispensing, Maintenance.

**Where it lives:** `MachineState` in `machine/machine_state.go` is the
shared interface. `IdleState`, `PreparingState`, `DispensingState`, and
`MaintenanceState` each implement it with their own rules. `CoffeeMachine`
holds whichever state is currently active and forwards requests to it:

```go
func (m *CoffeeMachine) PrepareBeverage(beverageType models.BeverageType) error {
	return m.State().PrepareBeverage(m, beverageType)
}
```

**An important Go specific detail:** `MachineState` and `CoffeeMachine` live
in the same package (`machine`), instead of being split into separate
packages. This is not just a style choice, it is required. States need to
call methods on `CoffeeMachine` (to switch to a new state, to read or
update ingredients), and `CoffeeMachine` needs to refer to `MachineState`
(to remember its current state). If these lived in two separate packages
that each imported the other, Go's compiler would refuse to build the
project at all, with an "import cycle not allowed" error. Types that depend
on each other this closely need to live together in one package.

### 4. Observer Pattern (`models.IngredientObserver` and `observer/alert_service.go`)

**In plain words:** one part of the system needs to announce "something
happened", without needing to know who is listening or what they will do
about it. Anyone interested registers themselves beforehand. When the event
happens, everyone who registered gets notified, one by one.

**Why this problem needs it:** the requirements ask for a low stock alert
when an ingredient drops below 20 percent of its capacity, but the
`Ingredient` type itself should not need to know anything about how alerts
are shown (printed to a console, emailed, sent to a dashboard light, and so
on).

**Where it lives:** the interface, `IngredientObserver`, is declared in
`models/ingredient.go`, right next to `Ingredient`, even though its one
implementation, `AlertService`, lives in a separate `observer` package.
This placement is deliberate, and it teaches a genuinely useful Go rule:

> The consumer of an interface, not the implementer, should define it.

`Ingredient` is the type that actually *calls* `OnIngredientLow(...)`, so it
is the "consumer" of this interface, and that is where the interface
belongs. `AlertService` is the "implementer": in Go, a type does not need
to say "implements IngredientObserver" anywhere, or even import the package
the interface lives in for any special reason. As long as it has a method
with a matching name and signature, it automatically counts as satisfying
that interface. This is sometimes called "structural typing" or "duck
typing": if it has the right shape, it fits.

If the interface had instead been declared in its own separate package,
that package would need to import `models` (for the `*Ingredient`
parameter type), while `models` would need to import it right back (for
the interface type on `Ingredient.observers`). That is an import cycle, and
Go would refuse to compile it. Declaring the interface next to its
consumer avoids the problem entirely.

## Other Go and LLD lessons that came up

These are not separate design patterns, but they came up constantly while
building this project, and each one taught something specific about how Go
works.

**Interfaces instead of abstract classes.** Some object oriented languages
let you declare an abstract class with an unimplemented method, and the
compiler refuses to let any subclass skip implementing it. Go has no such
feature. Instead, Go defines a small interface (like `Beverage`, with
`Name()`, `Recipe()`, `Prepare()`), and any type that has all of those
methods automatically satisfies it, with no keyword needed. To catch
mistakes early (for example, forgetting to write `Prepare()` on a new drink
type), this project uses a small, zero-cost trick everywhere a type is
meant to satisfy an interface:

```go
var _ Beverage = (*Coffee)(nil)
```

This line does nothing at runtime. It exists purely so the compiler checks,
right there, that `*Coffee` really does have every method `Beverage`
requires. If a method is missing or has the wrong signature, the build
fails immediately, with a clear error, instead of failing much later when
some other part of the program tries to actually use the type.

One real trap found while building this: `var _CustomizationStrategy = (*MilkCustomization)(nil)`
(no space after the underscore) looks almost identical to the check above,
but it is not the same thing at all. Without the space, `_CustomizationStrategy`
is parsed as one ordinary variable name, and its type is inferred from
whatever is assigned to it, so no interface check happens. This one
silently let a real bug through for several turns: a method with the wrong
signature that nothing ever caught, until the day some other code actually
tried to use it and the build failed. The lesson: that one space matters a
lot, and it is worth double checking any time this pattern is used.

**Enums, the Go way.** Go has no `enum` keyword. The standard way to build
one is three parts together: a named number type (`type SugarLevel int`), a
block of constants built with `iota` (which auto-numbers each line starting
from 0), and a `String()` method so the value prints as readable text
instead of a bare number:

```go
type SugarLevel int

const (
	None SugarLevel = iota
	Low
	Medium
	High
)

func (s SugarLevel) String() string { /* ... */ }
```

One easy mistake here: the type must be written on the *first* constant
only (`None SugarLevel = iota`). Every line after that repeats both the
type and the `iota` expression automatically. Leaving the type off
entirely still compiles, but produces a plain `int` with none of the type
safety an enum is supposed to give.

**Encapsulation without a `private` keyword.** Go does not have `public`
and `private`. Instead, whether a name can be seen from outside its package
depends entirely on its first letter: capitalized names (`Recipe`,
`AddIngredient`) are exported and usable from anywhere, lowercase names
(`ingredients`, `crossedThresholdLocked`) are only visible inside the same
package. Several types in this project (like `Recipe` and `CoffeeMachine`)
also return *copies* of internal maps from their getters, instead of the
live internal map itself, so that outside code cannot quietly reach in and
corrupt internal state (maps in Go are reference types, so handing one out
directly would let a caller mutate it without the owner ever knowing).

**Constructors are just functions.** Go has no `new ClassName(...)` syntax
and no constructor overloading. The convention is a plain function named
`New<Type>`, returning a pointer:

```go
func NewRecipe() *Recipe { /* ... */ }
func NewRecipeFromIngredients(ingredients map[string]int) *Recipe { /* ... */ }
```

**Errors are values, not exceptions.** Many languages signal failure (an
unrecognized beverage type, insufficient ingredients) by throwing an
exception. Go has no exceptions for this kind of everyday failure.
Functions that can fail simply return an extra `error` value alongside
their normal result, and the caller is expected to check it:

```go
beverage, err := beverages.NewBeverage(beverageType)
if err != nil {
	return err
}
```

**Concurrency and thread safety.** The requirements specifically call for
thread safe ingredient handling, since two requests could come in at the
same time and both see "enough" stock, even though together they would use
more than what is actually available. `Ingredient` protects its quantity
with a `sync.Mutex`, and `Consume` checks and deducts stock as a single,
uninterrupted operation while holding that lock, so two concurrent
`Consume` calls can never both succeed against stock that is only enough
for one of them.

One important rule followed throughout this project: **never call code you
do not control while holding your own lock.** When an ingredient's stock
crosses below the low stock threshold, it needs to notify its observers,
but it always unlocks first, and only then calls
`observer.OnIngredientLow(...)`. Go's `sync.Mutex` is not reentrant (the
same goroutine cannot lock it a second time while already holding it), so
calling back into code that might try to lock the same mutex again, while
still holding it, would freeze the program forever. This
also protects against a different mistake Go's own tools actively check
for: passing a struct that contains a mutex by value (a copy) instead of by
pointer, which silently produces two independent locks that no longer
protect anything together. `go vet` has a dedicated check for exactly this
(`copylocks`), and it is part of what `go test -race` and `go vet` verify
in this project.

**Testing.** Go's testing tools are built in, no external framework needed.
Test files end in `_test.go` and live next to the code they test. This
project uses two styles:

- `models/ingredient_test.go` uses `package models` (same package as the
  code under test), since some of its tests need to closely inspect
  behavior like exact notification timing.
- `machine/coffee_machine_test.go` uses `package machine_test` (an
  external, "black box" package), since everything it needs to test is
  already exposed through `CoffeeMachine`'s public methods. This is
  generally the preferred style when it is possible: it tests the same
  public API that any other real caller would use, so it cannot
  accidentally depend on some internal detail that later changes.

The most important test in the project is
`TestIngredient_ConcurrentConsume`, which launches 100 goroutines that all
try to consume stock from the same `Ingredient` at once, then checks that
the final quantity is exactly what it should be, no more, no less. Running
this under `go test -race` is what actually proves the thread safety
requirement is met, rather than just assumed.

## Walking through one full request

Here is what happens, step by step, when the demo asks for a Latte with
extra sugar and almond milk:

1. `main.go` calls `m.AddCustomization(...)` twice, once for sugar, once
   for milk. These are stored on the machine (**Strategy Pattern**).
2. `main.go` calls `m.PrepareBeverage(models.Latte)`. `CoffeeMachine`
   forwards this straight to whatever state is currently active
   (**State Pattern**). Assuming the machine is Idle, `IdleState` handles it.
3. `IdleState` asks `CoffeeMachine.ValidateRecipe` whether there is enough
   water, coffee beans, and milk in stock for a Latte. This internally asks
   `beverages.NewBeverage(models.Latte)` for a fresh `*Latte` just to read
   its recipe (**Factory Pattern**).
4. If there is enough stock, `IdleState` creates the real `*Latte` to use
   (**Factory Pattern**, again), stores it as the machine's current
   beverage, and switches the machine into `PreparingState`.
5. `main.go` calls `m.Dispense()`. Since the machine is now in
   `PreparingState`, this actually runs the Latte's `Prepare()` method
   (which prints the brewing steps), consumes the recipe's ingredients from
   stock, and switches to `DispensingState`.
6. If consuming milk happens to push it below 20 percent of capacity,
   `Ingredient` notifies every registered observer, which in this project
   means `AlertService` prints a warning (**Observer Pattern**).
7. `main.go` reads the current beverage and calls
   `m.ApplyCustomizations(latte)`, which loops through the two stored
   customizations and calls `Customize` on each (**Strategy Pattern**,
   completing the loop from step 1).
8. `main.go` calls `m.Dispense()` again. Since the machine is now in
   `DispensingState`, this hands the drink over, clears the current
   beverage, and switches the machine back to `IdleState`, ready for the
   next request.

## Real bugs found and fixed while building this

Keeping this list because each one was a genuine, useful lesson, not just a
typo.

1. **A recipe built with a bare `Recipe{}` instead of `NewRecipe()`.** The
   zero value of a Go map is `nil`. Reading from a nil map is safe, but
   writing to one panics. Any code path that skipped the real constructor
   would have panicked the first time it tried to add an ingredient.

2. **A getter that returned a live internal map directly.** Since maps in
   Go are reference types, handing one out directly would let any caller
   quietly mutate a `Recipe`'s or `CoffeeMachine`'s internal state from the
   outside. Fixed by returning a fresh copy instead, so the caller gets an
   independent map it can't use to corrupt the original.

3. **An import cycle between `models` and a separate observer package.**
   Go refuses to compile two packages that import each other in a circle.
   Fixed by moving the `IngredientObserver` interface into `models`, next
   to the `Ingredient` type that actually uses it.

4. **A struct containing a `sync.Mutex` passed by value.** An early draft
   of `IngredientObserver` accepted `models.Ingredient` instead of
   `*models.Ingredient`. Since `Ingredient` contains a mutex, passing it by
   value would copy the lock, and the copy would protect nothing. Go's own
   `go vet` tool caught this immediately with its `copylocks` check.

5. **A factory left half finished.** `beverages.NewBeverage` was missing
   the `Cappuccino` and `Latte` cases for a long stretch, along with a
   placeholder `"???"` error message. Nothing caught this until `main.go`
   actually tried to prepare those drinks and got a mysterious error,
   which is exactly why exercising the full flow end to end, not just each
   piece in isolation, matters.

6. **A silently broken compile-time safety check**, `var _CustomizationStrategy = (*MilkCustomization)(nil)`
   missing a space, described in more detail above. It looked like a real
   check, compiled fine, and checked nothing.

Each of these was caught either by `go build`, `go vet`, `go test -race`,
or by actually running the full program end to end, which is a good
reminder that all four of those are worth doing, not just one of them.
