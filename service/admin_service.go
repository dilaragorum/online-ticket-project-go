package service

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
)

var (
	ErrAlreadyCreatedTrip = errors.New("this trip is already created")
)

type AdminService interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
}

type adminService struct {
	adminRepo repository.AdminRepository
}

func NewAdminService(adminRepo repository.AdminRepository) *adminService {
	return &adminService{adminRepo: adminRepo}
}

func (as *adminService) CreateTrip(ctx context.Context, trip *model.Trip) error {
	if err := as.adminRepo.CreateTrip(ctx, trip); err != nil {
		if errors.Is(err, repository.ErrDuplicateIdx) {
			return ErrAlreadyCreatedTrip
		}
		return err
	}

	return nil
}
