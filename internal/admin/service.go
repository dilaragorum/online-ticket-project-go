package admin

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/trip"
)

var (
	ErrAlreadyCreatedTrip = errors.New("this trip is already created")
	ErrTripNotExist       = errors.New("this trip does not exist")
)

type Service interface {
	CreateTrip(ctx context.Context, trip *trip.Trip) error
	CancelTrip(ctx context.Context, id int) error
}

type service struct {
	tripRepo trip.Repository
}

func NewAdminService(tripRepo trip.Repository) *service {
	return &service{tripRepo: tripRepo}
}

func (as *service) CreateTrip(ctx context.Context, t *trip.Trip) error {
	if err := as.tripRepo.Create(ctx, t); err != nil {
		if errors.Is(err, trip.ErrDuplicateIdx) {
			return ErrAlreadyCreatedTrip
		}
		return err
	}

	return nil
}

func (as *service) CancelTrip(ctx context.Context, id int) error {
	if err := as.tripRepo.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, trip.ErrTripNotFound):
			return ErrTripNotExist
		}
		return err
	}

	return nil
}
