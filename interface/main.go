package main

import "fmt"

type PaymentProcessor interface {
	ProcessPayment(amount float64) bool
	Refund(transactionId string) bool
}

func validateAmount(amount float64) bool {
	return amount > 0
}

type CreditCardProcessor struct {
	APIKey string
}

func (c *CreditCardProcessor) ProcessPayment(amount float64) bool {
	if !validateAmount(amount) {
		return false
	}

	fmt.Printf("Processing $%g via credit card\n", amount)
	return true
}

func (*CreditCardProcessor) Refund(transactionID string) bool {
	fmt.Printf("Refunding transaction %s via credit card\n", transactionID)
	return true
}

type PayPalProcessor struct{}

func (*PayPalProcessor) ProcessPayment(amount float64) bool {
	if !validateAmount(amount) {
		return false
	}
	fmt.Printf("Processing $%g via PayPal\n", amount)
	return true
}

func (*PayPalProcessor) Refund(transactionID string) bool {
	fmt.Printf("Refunding transaction %s via PayPal\n", transactionID)
	return true
}

func checkout(p PaymentProcessor, amount float64) {
	p.ProcessPayment(amount)
}

func main() {
	cc := &CreditCardProcessor{APIKey: "api_key_123"}
	pp := &PayPalProcessor{}

	checkout(cc, 100.0)
	checkout(pp, 100.0)
}
