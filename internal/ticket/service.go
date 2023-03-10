package ticket

import (
	"context"
	"errors"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/internal/auth"
	"github.com/dilaragorum/online-ticket-project-go/internal/notification"
	"github.com/dilaragorum/online-ticket-project-go/internal/payment"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
)

var (
	ErrNoCapacity = errors.New("capacity is full")
	ErrTripNotFound = errors.New("this trip does not exist")
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

	requestedTrip, err := s.tripRepo.FindByTripID(ctx, ticket.TripID)
	if err != nil {
		if errors.Is(err, trip.ErrTripNotFound) {
			return ErrTripNotFound
		}
		return err
	}

	param := notification.Param{
		Channel:     notification.SMS,
		To:          ticket.Phone,
		From:        "company ticket",
		Title:       "Purchase Detail",
		Description: fmt.Sprintf("Traveler Name: %s FromTo: %s-%s Date: %s Vehicle: %s", ticket.FullName, requestedTrip.From, requestedTrip.To, requestedTrip.Date, requestedTrip.Vehicle),
		LogMsg:      fmt.Sprintf("The %s who has %d id purchase ticket/s", claims.Username, claims.UserID),
	}

	purchasedTicket := Ticket{
		TripID: requestedTrip.ID,
		UserID: claims.UserID,
		Passenger: Passenger{
			Gender:   ticket.Gender,
			FullName: ticket.FullName,
			Email:    ticket.Email,
			Phone:    ticket.Phone,
		},
	}

	if requestedTrip.AvailableSeat == 0 {
		return ErrNoCapacity
	}

	//Trips tablosundan capacity'den alınan bilet sayısı kadar kişi düşeceğiz.
	if err = s.tripRepo.UpdateAvailableSeat(ctx, requestedTrip.ID, 1); err != nil {
		return err
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
