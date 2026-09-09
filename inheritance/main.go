package main

import "fmt"

type Employee struct {
	Name       string
	EmployeeID int
}

func (e *Employee) GetInfo() string {
	return fmt.Sprintf("Employee: %s (ID: %d)", e.Name, e.EmployeeID)
}

type Manager struct {
	Employee
	Department string
}

func (m *Manager) GetInfo() string {
	return fmt.Sprintf("%s, Department: %s", m.Employee.GetInfo(), m.Department)
}

type Director struct {
	Manager
	Budget float64
}

func (d *Director) GetInfo() string {
	return fmt.Sprintf("%s, Budget: $%.2f", d.Manager.GetInfo(), d.Budget)
}

func main() {
	manager := &Manager{Employee: Employee{Name: "Alice", EmployeeID: 101}, Department: "Engineering"}
	fmt.Println(manager.GetInfo())

	director := &Director{Manager: Manager{Employee: Employee{Name: "Bob", EmployeeID: 102}, Department: "Engineering"}, Budget: 1000000.0}
	fmt.Println(director.GetInfo())
}
