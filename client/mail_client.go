package client

import (
	"fmt"
)

type MailClient interface {
	Send(to, from, title, description string) error
}
type mailClient struct {
}

func NewMail() *mailClient {
	return &mailClient{}
}

func (m *mailClient) Send(to, from, title, description string) error {
	fmt.Printf("Email was sent from=%s to=%s title=%s description=%s", from, to, title, description)
	return nil
}
