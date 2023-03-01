package mail

import (
	"fmt"
)

type Client interface {
	Send(to, from, title, description string) error
}
type client struct {
}

func NewMail() *client {
	return &client{}
}

func (m *client) Send(to, from, title, description string) error {
	fmt.Printf("Email was sent from=%s to=%s title=%s description=%s", from, to, title, description)
	return nil
}
