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

const (
	CorporatedLimit       = 20
	IndividualLimit       = 5
	LeastMaleTicketNumber = 2
)

var (
	ErrNoCapacity   = errors.New("capacity is full")
	ErrTripNotFound = errors.New("this trip does not exist")

	// TODO: handlerdaki gibi fonksiyon yapalım.
	ErrExceedAllowedTicketToPurchaseForTwenty = errors.New("exceed number of tickets allowed to be purchased(20)")
	ErrExceedAllowedTicketToPurchaseForFive   = errors.New("exceed number of tickets allowed to be purchased(5)")

	ErrExceedMaleTicketNumber = errors.New("exceed number of male ticket allowed to be purchased")
)

type Service interface {
	Purchase(ctx context.Context, tickets []Ticket, claims auth.Claims) error
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

func (s *defaultService) Purchase(ctx context.Context, tickets []Ticket, claims auth.Claims) error {
	if err := checkCorporatedLimit(claims, tickets); err != nil {
		return err
	}

	if err := checkIndividualLimit(claims, tickets); err != nil {
		return err
	}

	if err := checkMaleTicketLimit(tickets, claims); err != nil {
		return err
	}

	params := make([]notification.Param, 0, len(tickets))
	passengersNames := make([]string, 0, len(tickets))

	for i := range tickets {
		ticket := tickets[i]

		passengersNames = append(passengersNames, ticket.FullName)

		requestedTrip, err := s.tripRepo.FindByTripID(ctx, ticket.TripID)
		if err != nil {
			if errors.Is(err, trip.ErrTripNotFound) {
				return ErrTripNotFound
			}
			return err
		}

		params = append(params, notification.Param{
			Channel: notification.SMS,
			To:      ticket.Phone,
			From:    "X Ticket Company",
			Title:   "Purchase Detail",
			Description: fmt.Sprintf(`Congrats! Your transaction is successful. Here your ticket Details:
FromTo: %s-%s
Date: %s
Vehicle: %s
Passengers:`, requestedTrip.From, requestedTrip.To, requestedTrip.Date, requestedTrip.Vehicle),
			LogMsg: fmt.Sprintf("The %s who has %d id purchase ticket/s", claims.Username, claims.UserID),
		})

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

		if ok := requestedTrip.CheckAvailableSeat(len(tickets)); !ok {
			return ErrNoCapacity
		}

		if err = s.tripRepo.UpdateAvailableSeat(ctx, requestedTrip.ID, len(tickets)); err != nil {
			return ErrNoCapacity
		}

		if err = s.ticketRepo.CreateTicketWithDetails(ctx, &purchasedTicket); err != nil {
			return err
		}
	}

	if err := s.payment.Transfer(); err != nil {
		return err
	}

	for i := range params {
		for k := range passengersNames {
			params[i].Description += fmt.Sprintf("%s\n", passengersNames[k])
		}
	}

	for i := 0; i < len(params); i++ {
		if err := s.notificationService.Send(ctx, params[i]); err != nil {
			return err
		}
	}

	return nil
}

func checkIndividualLimit(claims auth.Claims, tickets []Ticket) error {
	if claims.IsIndividualUser() && len(tickets) > IndividualLimit {
		return ErrExceedAllowedTicketToPurchaseForFive
	}
	return nil
}

func checkCorporatedLimit(claims auth.Claims, tickets []Ticket) error {
	if claims.IsCorporatedUser() && len(tickets) > CorporatedLimit {
		return ErrExceedAllowedTicketToPurchaseForTwenty
	}
	return nil
}

func checkMaleTicketLimit(tickets []Ticket, claims auth.Claims) error {
	var maleNum int

	for i := range tickets {
		if tickets[i].Gender == Male {
			maleNum++
		}
	}

	if claims.IsIndividualUser() && maleNum > LeastMaleTicketNumber {
		return ErrExceedMaleTicketNumber
	}

	return nil
}
