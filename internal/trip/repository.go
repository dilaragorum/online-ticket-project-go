package trip

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/user"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateIdx = errors.New(`ERROR: duplicate key value violates unique constraint "idx_trips_idx_member" (SQLSTATE 23505)`)
	ErrTripNotFound = errors.New("there is no trip with that ID")
)

type Repository interface {
	Create(ctx context.Context, trip *Trip) error
	Delete(ctx context.Context, id int) error
	FindByFilter(ctx context.Context, trip *Filter) ([]Trip, error)
	FindByTripID(ctx context.Context, tripID int) (*Trip, error)
}

type repository struct {
	database *gorm.DB
}

func NewTripRepository(database *gorm.DB) Repository {
	return &repository{database: database}
}

func (t *repository) Create(ctx context.Context, trip *Trip) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := t.database.WithContext(timeoutCtx).Model(&Trip{}).Create(trip).Error; err != nil {
		if err.Error() == ErrDuplicateIdx.Error() {
			return ErrDuplicateIdx
		}

		log.Error(err)
		return err
	}

	return nil
}

func (t *repository) Delete(ctx context.Context, id int) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := t.database.WithContext(timeoutCtx).Delete(&Trip{}, id).Error; err != nil {
		switch {
		case errors.Is(err, user.ErrNoRecord):
			return ErrTripNotFound
		default:
			log.Error(err)
			return err
		}
	}

	return nil
}

func (t *repository) FindByFilter(ctx context.Context, filter *Filter) ([]Trip, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var trips []Trip

	if err := t.database.WithContext(timeoutCtx).Where(&Trip{
		ID:      filter.TripID,
		From:    filter.From,
		To:      filter.To,
		Vehicle: filter.Vehicle,
		Date:    filter.Date,
	}).Find(&trips).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return trips, nil
}

func (t *repository) FindByTripID(ctx context.Context, tripID int) (*Trip, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var trip Trip

	if err := t.database.WithContext(timeoutCtx).Where(&Trip{
		ID: tripID,
	}).Find(&trip).Error; err != nil {
		log.Error(err)
		return nil, err
	}

	return &trip, nil
}
