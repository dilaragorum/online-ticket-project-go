package notification

import (
	"context"
	"fmt"
)

type Mail struct {
}

func NewMail() *Mail {
	return &Mail{}
}

func (m *Mail) Send(ctx context.Context, param Param) error {
	msg := fmt.Sprintf("Email was sent from=%s to=%s title=%s description=%s", param.From, param.To, param.Title, param.Description)
	fmt.Println(msg)
	return nil
}
