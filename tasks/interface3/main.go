package main

import (
	"errors"
	"fmt"
)

type PaymentProcessor interface {
	Process(amount float64) error
	Verify(amount float64) bool
}

type CreditCardProcessor struct {
	limit float64
}

func (c *CreditCardProcessor) Verify(amount float64) bool {
	return amount <= c.limit
}

func (c *CreditCardProcessor) Process(amount float64) error {
	if amount > c.limit {
		return errors.New("amount exceeds limit")
	}

	c.limit -= amount

	fmt.Println("Payment processor received ", amount)

	return nil
}

type PayPalProcessor struct {
	balance float64
}

func (p *PayPalProcessor) Process(amount float64) error {
	if amount > p.balance {
		return errors.New("amount exceeds limit")
	}

	p.balance -= amount
	fmt.Println("Payment processor received ", amount)
	return nil
}

func (p *PayPalProcessor) Verify(amount float64) bool {
	return amount <= p.balance
}

func ExecutePayment(processor PaymentProcessor, amount float64) {
	if processor.Verify(amount) {
		err := processor.Process(amount)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Payment processing failed due to invalid amount.")
	}
}

func main() {
	creditCardProcessor := &CreditCardProcessor{
		limit: 100,
	}

	payPalProcessor := &PayPalProcessor{
		balance: 200,
	}

	ExecutePayment(creditCardProcessor, 50)
	ExecutePayment(creditCardProcessor, 50)
	ExecutePayment(payPalProcessor, 150)
	ExecutePayment(payPalProcessor, 150)
}
