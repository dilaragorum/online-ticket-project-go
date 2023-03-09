package ticket

import (
	"context"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/payment"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
)

type Service interface {
	Purchase(ctx context.Context, ticket *Ticket, claims auth.Claims) error
}

type defaultService struct {
	ticketRepo          Repository
	notificationService notification.Service
	tripRepo            trip.Repository
	payment             payment.Client
}

func NewService(ticketRepo Repository, notificationService notification.Service, tripRepo trip.Repository, payment payment.Client) Service {
	return &defaultService{ticketRepo: ticketRepo, notificationService: notificationService, tripRepo: tripRepo, payment: payment}
}

func (s *defaultService) Purchase(ctx context.Context, ticket *Ticket, claims auth.Claims) error {
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
		LogMsg:      fmt.Sprintf("The %s who has %d id purchase ticket/s", claims.Username, claims.UserID),
	}

	purchasedTicket := Ticket{
		TripID: trip.ID,
		UserID: claims.UserID,
		Passenger: Passenger{
			Gender:   ticket.Gender,
			FullName: ticket.FullName,
			Email:    ticket.Email,
			Phone:    ticket.Phone,
		},
	}

	err = s.ticketRepo.CreateTicketWithDetails(ctx, &purchasedTicket)
	if err != nil {
		return err
	}

	if err = s.notificationService.Send(ctx, param); err != nil {
		return err
	}

	return nil
}
