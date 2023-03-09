package trip

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
)

var (
	ErrAlreadyCreatedTrip = errors.New("this trip is already created")
	ErrTripNotExist       = errors.New("this trip does not exist")
)

type Service interface {
	FilterTrips(ctx context.Context, trip *Filter) ([]Trip, error)
	CreateTrip(ctx context.Context, trip *Trip) error
	CancelTrip(ctx context.Context, id int) error
	GetSoldTicketNumber(ctx context.Context, tripID int) (int, error)
	GetTotalRevenueForSpecificTrip(ctx context.Context, tripID int) (float64, error)
}

type defaultService struct {
	tripRepo Repository
}

func NewTripService(tripRepo Repository) Service {
	return &defaultService{tripRepo: tripRepo}
}

func (s *defaultService) FilterTrips(ctx context.Context, filter *Filter) ([]Trip, error) {
	trips, err := s.tripRepo.FindByFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	if len(trips) == 0 {
		return nil, user.ErrThereIsNoTrip
	}

	return trips, err
}

func (s *defaultService) CreateTrip(ctx context.Context, t *Trip) error {
	if err := s.tripRepo.Create(ctx, t); err != nil {
		if errors.Is(err, ErrDuplicateIdx) {
			return ErrAlreadyCreatedTrip
		}
		return err
	}

	return nil
}

func (s *defaultService) CancelTrip(ctx context.Context, id int) error {
	if err := s.tripRepo.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, ErrTripNotFound):
			return ErrTripNotExist
		}
		return err
	}

	return nil
}

func (s *defaultService) GetSoldTicketNumber(ctx context.Context, tripID int) (int, error) {
	number, err := s.tripRepo.GetSoldTicketNumber(ctx, tripID)
	if err != nil {
		return -1, err
	}

	return number, nil
}

func (s *defaultService) GetTotalRevenueForSpecificTrip(ctx context.Context, tripID int) (float64, error) {
	trip, err := s.tripRepo.FindByTripID(ctx, tripID)
	if err != nil {
		return -1, err
	}

	soldTicketNumber, err := s.tripRepo.GetSoldTicketNumber(ctx, tripID)
	if err != nil {
		return -1, err
	}

	total := trip.Price * float64(soldTicketNumber)

	return total, nil
}
