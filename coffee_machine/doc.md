Design a coffee machine system that can prepare various beverages (coffee, tea, cappuccino, latte)
based on predefined recipes, manage ingredient inventory (coffee beans, milk, water, sugar),
allow user customization (sugar level, milk type), monitor ingredient levels with alerts,
and handle the dispensing process with proper state management.

In this problem, you’ll design a system that handles multiple beverage types,
user customizations, and ensures the machine remains in a valid state throughout the brewing process.


Coffee Machine is perfect for learning LLD because it involves:

State Pattern - Managing the machine’s lifecycle (Idle, Preparing, Dispensing, Maintenance).
Strategy Pattern - Handling various customization options like sugar levels and milk types.
Factory Pattern - Creating beverages based on predefined recipes.
Observer Pattern - Real-time monitoring and alerts for ingredient levels.

Core Requirements
Functional Requirements:

Beverage Variety: Support multiple types (Coffee, Tea, Cappuccino, Latte) using the Factory Pattern.
Inventory Management: Track ingredients (beans, milk, water, sugar) with real-time quantity tracking.
Recipe System: Each beverage has a predefined recipe; validate availability before starting.
Customization: Allow users to select sugar level (None to High) and milk type (Regular, Skim, Almond).
Alert System: Monitor ingredient levels and notify when they fall below a threshold (e.g., 20%).
State Management: Manage operational states (Idle, Preparing, Dispensing, Maintenance).
Preparation Flow: Handle the full lifecycle: validate → consume → prepare → dispense → Idle.
Error Handling: Gracefully reject requests if validation fails or ingredients run out mid-brew.
Maintenance: Support refilling ingredients to restore inventory levels.
Non-Functional Requirements:

Thread Safety: Ensure inventory operations are thread-safe for concurrent preparation requests.
Modular Design: Each class should have well-defined roles following the Single Responsibility Principle.
Extensibility: Easy to add new beverages, ingredients, or customizations without refactoring.
Reliability: The core preparation logic and state management must be easy to test and maintain.

What’s Expected?
1. System Architecture
The system coordinates between the User Interface, the Inventory Manager, and the Brewing Engine.

Diagram
2. Key Classes to Design
Diagram
System Flow
Beverage Preparation Flow
Diagram
Key Design Challenges
1. Machine State Transitions
A coffee machine shouldn’t allow a user to select a drink while it’s already brewing or when it’s in a maintenance/error state.

Solution: Use the State Pattern. Define states like IdleState, BrewingState, DispensingState, and OutOfOrderState. Each state object defines what actions are allowed, preventing invalid transitions.

2. Ingredient Monitoring
The system needs to alert someone when beans or milk are low, but the Inventory shouldn’t be tightly coupled to a notification system.

Solution: Use the Observer Pattern. The Inventory acts as a Subject, and various AlertSystems (Email, Console, LED) act as Observers that get notified when levels cross a threshold.

3. Concurrency in Inventory
If two requests come in simultaneously, they might both see enough milk available but together exceed the stock.

Solution: Implement Thread-Safe Inventory. Use locks or atomic operations (like ConcurrentHashMap and AtomicInteger) to ensure that checking and consuming ingredients is an atomic operation.

What You’ll Learn
By solving this problem, you’ll master:

✅ State Management - Controlling complex object lifecycles.
✅ Inventory Logic - Managing resource allocation and tracking.
✅ Decoupling - Using Observers to separate logic from notifications.
✅ Factory Method - Dynamically creating objects based on configuration.


Functional Requirements
The coffee machine should support multiple beverage types (coffee, tea, cappuccino, latte) using the Factory Pattern to create beverage instances.
The system should manage ingredient inventory (coffee beans, milk, water, sugar) with real-time quantity tracking.
Each beverage should have a predefined recipe specifying required ingredients and their quantities.
Before preparing a beverage, the system must validate that all required ingredients are available in sufficient quantities according to the recipe.
The system should monitor ingredient levels and alert when any ingredient falls below a low threshold (e.g., 20% of capacity) using the Observer Pattern.
Users should be able to customize beverages by selecting sugar level (none, low, medium, high) and milk type (regular, skim, almond) using the Strategy Pattern.
The machine should manage its operational states (Idle, Preparing, Dispensing, Maintenance) using the State Pattern with clear state transitions.
The machine should handle the complete beverage preparation flow: validate recipe → consume ingredients → prepare beverage → dispense beverage → return to Idle.
If ingredient validation fails, the system should reject the preparation request and provide appropriate feedback.
The system should support refilling ingredients to restore inventory levels.
Non-Functional Requirements
The design should be object-oriented, with each class having a well-defined role and responsibilities clearly separated.
The system should be built in a modular and flexible way, making it easy to add future enhancements such as additional beverage types, ingredients, or customization options.
Ingredient inventory operations should be thread-safe to handle concurrent beverage preparation requests without race conditions.
The core beverage preparation logic, recipe validation, and state management should be easy to test, understand, and maintain over time.
The system should use the Factory Pattern to create different beverage types, making it easy to add new beverages without modifying existing code.
The system should use the Strategy Pattern to handle customization options, making customization algorithms swappable.
The system should use the State Pattern to manage machine states, ensuring appropriate behavior during different phases of operation.
The system should use the Observer Pattern to monitor ingredient levels and notify when supplies are low.