package ticket

import (
	"context"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/payment"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
)

type Service interface {
	Purchase(ctx context.Context, ticket *Ticket) error
}

type defaultService struct {
	notificationService notification.Service
	tripRepo            trip.Repository
	payment             payment.Client
}

func NewService(notificationService notification.Service, tripRepo trip.Repository, payment payment.Client) Service {
	return &defaultService{notificationService: notificationService, tripRepo: tripRepo, payment: payment}
}

func (s *defaultService) Purchase(ctx context.Context, ticket *Ticket) error {
	if err := s.payment.Transfer(); err != nil {
		return err
	}

	trip, err := s.tripRepo.FindByTripID(ctx, ticket.TripID)
	if err != nil {
		return err
	}

	param := notification.Param{
		Channel:     notification.SMS,
		To:          ticket.Phone,
		From:        "company ticket",
		Title:       "Purchase Detail",
		Description: fmt.Sprintf("Traveler Name: %s FromTo: %s-%s Date: %s Vehicle: %s", ticket.FullName, trip.From, trip.To, trip.Date, trip.Vehicle),
		LogMsg:      fmt.Sprintf("The user called as %s purchase ticket/s from which has %d id", ticket.UserID, ticket.ID),
	}

	if err = s.notificationService.Send(ctx, param); err != nil {
		return err
	}

	return nil
}
