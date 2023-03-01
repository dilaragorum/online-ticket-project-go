package mail

import (
	"context"
	"errors"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/pkg/mail"
)

const (
	CompanyEmail = "onlineticket@wonderful.com"
)

var (
	MailCanNotSent = errors.New("mail could not be sent")
)

type Service interface {
	SendWelcomeMail(ctx context.Context, email string) error
}

type service struct {
	mailClient             mail.Client
	notificationRepository notification.Repository
}

func NewService(mailClient mail.Client, notificationRepository notification.Repository) *service {
	return &service{mailClient: mailClient, notificationRepository: notificationRepository}
}

func (m *service) SendWelcomeMail(ctx context.Context, email string) error {
	if err := m.mailClient.Send(email, CompanyEmail, "Welcome", "Welcome to our Platform!"); err != nil {
		return MailCanNotSent
	}

	if err := m.notificationRepository.Create(ctx, notification.ChannelEMAIL, fmt.Sprintf("A welcome e-mail has been sent to %s, who has just registered.", email)); err != nil {
		return MailCanNotSent
	}

	return nil
}
