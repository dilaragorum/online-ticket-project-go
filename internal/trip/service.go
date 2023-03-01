package trip

import (
	"context"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
)

type Service interface {
	FilterTrips(ctx context.Context, trip *Filter) ([]Trip, error)
}

type service struct {
	tripRepo Repository
}

func NewTripService(tripRepo Repository) *service {
	return &service{tripRepo: tripRepo}
}

func (s *service) FilterTrips(ctx context.Context, filter *Filter) ([]Trip, error) {
	trips, err := s.tripRepo.FindByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(trips) == 0 {
		return nil, user.ErrThereIsNoTrip
	}

	return trips, err
}
