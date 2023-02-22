package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/client"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
)

const (
	CompanyEmail = "onlineticket@wonderful.com"
)

var (
	MailCanNotSent = errors.New("mail could not be sent")
)

type MailService interface {
	SendWelcomeMail(ctx context.Context, email string) error
}

type mailService struct {
	mailClient             client.MailClient
	notificationRepository repository.NotificationRepository
}

func NewMailService(mailClient client.MailClient, notificationRepository repository.NotificationRepository) *mailService {
	return &mailService{mailClient: mailClient, notificationRepository: notificationRepository}
}

func (m *mailService) SendWelcomeMail(ctx context.Context, email string) error {
	if err := m.mailClient.Send(email, CompanyEmail, "Welcome", "Welcome to our Platform!"); err != nil {
		return MailCanNotSent
	}

	if err := m.notificationRepository.Create(ctx, model.ChannelEMAIL, fmt.Sprintf("A welcome e-mail has been sent to %s, who has just registered.", email)); err != nil {
		return MailCanNotSent
	}

	return nil
}
