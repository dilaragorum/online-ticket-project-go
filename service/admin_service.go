package service

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
)

var (
	ErrAlreadyCreatedTrip = errors.New("this trip is already created")
	ErrTripNotExist       = errors.New("this trip does not exist")
)

type AdminService interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
	CancelTrip(ctx context.Context, id int) error
}

type adminService struct {
	tripRepo repository.TripRepository
}

func NewAdminService(tripRepo repository.TripRepository) *adminService {
	return &adminService{tripRepo: tripRepo}
}

func (as *adminService) CreateTrip(ctx context.Context, trip *model.Trip) error {
	if err := as.tripRepo.Create(ctx, trip); err != nil {
		if errors.Is(err, repository.ErrDuplicateIdx) {
			return ErrAlreadyCreatedTrip
		}
		return err
	}

	return nil
}

func (as *adminService) CancelTrip(ctx context.Context, id int) error {
	if err := as.tripRepo.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, repository.ErrTripNotFound):
			return ErrTripNotExist
		}
		return err
	}

	return nil
}
