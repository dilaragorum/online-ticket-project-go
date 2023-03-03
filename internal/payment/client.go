package payment

import "fmt"

type Client interface {
	Transfer() error
}

type defaultClient struct {
}

func NewClient() Client {
	return &defaultClient{}
}

func (p *defaultClient) Transfer() error {
	fmt.Println("Payment Received")
	return nil
}
