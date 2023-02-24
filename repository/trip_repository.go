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
	ErrTripNotFound = errors.New("there is no trip with that ID")
)

type TripRepository interface {
	Create(ctx context.Context, trip *model.Trip) error
	Delete(ctx context.Context, id int) error
}

type tripRepository struct {
	database *gorm.DB
}

func NewTripRepository(database *gorm.DB) *tripRepository {
	return &tripRepository{database: database}
}

func (tr *tripRepository) Create(ctx context.Context, trip *model.Trip) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := tr.database.WithContext(timeoutCtx).Model(&model.Trip{}).Create(trip).Error; err != nil {
		if err.Error() == ErrDuplicateIdx.Error() {
			return ErrDuplicateIdx
		}

		log.Error(err)
		return err
	}

	return nil
}

func (tr *tripRepository) Delete(ctx context.Context, id int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := tr.database.WithContext(timeoutCtx).Delete(&model.Trip{}, id).Error; err != nil {
		switch {
		case errors.Is(err, ErrNoRecord):
			return ErrTripNotFound
		default:
			log.Error(err)
			return err
		}
	}

	return nil
}
