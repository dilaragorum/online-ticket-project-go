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
	ErrNoCapacity   = errors.New("capacity is full")
	ErrTripNotFound = errors.New("this trip does not exist")

	ErrExceedAllowedTicketToPurchaseForTwenty = errors.New("exceed number of tickets allowed to be purchased(20)")
	ErrExceedAllowedTicketToPurchaseForFive   = errors.New("exceed number of tickets allowed to be purchased(5)")
	ErrExceedMaleTicketNumber                 = errors.New("exceed number of male ticket allowed to be purchased")
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
	if claims.UserType == auth.CorporateUser && len(tickets) > 20 {
		return ErrExceedAllowedTicketToPurchaseForTwenty
	}

	if claims.UserType == auth.IndividualUser && len(tickets) > 5 {
		return ErrExceedAllowedTicketToPurchaseForFive
	}

	var maleNum int

	for i := range tickets {
		if tickets[i].Gender == Male {
			maleNum++
		}
	}

	if claims.UserType == auth.IndividualUser && maleNum > 2 {
		return ErrExceedMaleTicketNumber
	}

	if err := s.payment.Transfer(); err != nil {
		return err
	}

	params := make([]notification.Param, 0)
	passengersNames := make([]string, 0)

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
			Channel:     notification.SMS,
			To:          ticket.Phone,
			From:        "X Ticket Company",
			Title:       "Purchase Detail",
			Description: fmt.Sprintf("Congrats! Your transaction is successfull. Here your ticket Details:\n FromTo: %s-%s\n Date: %s\n Vehicle: %s\n Passengers:\n", requestedTrip.From, requestedTrip.To, requestedTrip.Date, requestedTrip.Vehicle),
			LogMsg:      fmt.Sprintf("The %s who has %d id purchase ticket/s", claims.Username, claims.UserID),
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

		if requestedTrip.AvailableSeat == 0 || requestedTrip.AvailableSeat < uint(len(tickets)) {
			return ErrNoCapacity
		}

		if err = s.tripRepo.UpdateAvailableSeat(ctx, requestedTrip.ID, len(tickets)); err != nil {
			return ErrNoCapacity
		}

		err = s.ticketRepo.CreateTicketWithDetails(ctx, &purchasedTicket)
		if err != nil {
			return err
		}

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
