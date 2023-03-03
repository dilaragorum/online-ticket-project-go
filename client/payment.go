package client

import "fmt"

type Payment interface {
	Transfer() error
}

type payment struct {
}

func NewPayment() Payment {
	return &payment{}
}

func (p *payment) Transfer() error {
	fmt.Println("Payment Received")
	return nil
}
