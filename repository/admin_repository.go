package repository

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateIdx = errors.New(`ERROR: duplicate key value violates unique constraint "idx_trips_idx_member" (SQLSTATE 23505)`)
)

type AdminRepository interface {
	CreateTrip(ctx context.Context, trip *model.Trip) error
}

type adminRepository struct {
	database *gorm.DB
}

func NewAdminRepository(database *gorm.DB) *adminRepository {
	return &adminRepository{database: database}
}

func (ar *adminRepository) CreateTrip(ctx context.Context, trip *model.Trip) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := ar.database.WithContext(timeoutCtx).Model(&model.Trip{}).Create(trip).Error; err != nil {
		if err.Error() == ErrDuplicateIdx.Error() {
			return ErrDuplicateIdx
		}

		log.Error(err)
		return err
	}

	return nil
}
