package notification

import (
	"context"
	"fmt"
)

type Sms struct {
}

func NewSms() *Sms {
	return &Sms{}
}

func (s *Sms) Send(ctx context.Context, param Param) error {
	msg := fmt.Sprintf("SMS was sent from %s to %s with title=%s description=%s", param.From, param.To, param.Title, param.Description)
	fmt.Println(msg)
	return nil
}
