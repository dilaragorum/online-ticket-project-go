package ticket

import (
	"context"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
)

type Service interface {
	Purchase(ctx context.Context, ticket *Ticket) error
}

type service struct {
	notificationService notification.Service
	tripRepo            trip.Repository
}

func NewService(notificationService notification.Service, tripRepo trip.Repository) *service {
	return &service{notificationService: notificationService, tripRepo: tripRepo}
}

func (s *service) Purchase(ctx context.Context, ticket *Ticket) error {
	//Ödeme İşlemi

	//Mesaj atma
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
		LogMsg:      "şu kullanıcı şöyle bir satın alma yaptı",
	}

	if err = s.notificationService.Send(ctx, param); err != nil {
		return err
	}

	return nil
}
