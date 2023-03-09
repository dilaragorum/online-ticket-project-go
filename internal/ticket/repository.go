package ticket

import (
	"context"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	CreateTicketWithDetails(ctx context.Context, ticket *Ticket) error
}

type repository struct {
	database *gorm.DB
}

func NewTicketRepository(database *gorm.DB) Repository {
	return &repository{database: database}
}

func (r *repository) CreateTicketWithDetails(ctx context.Context, ticket *Ticket) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&Ticket{}).Create(ticket).Error; err != nil {
		log.Error(err)
		return err
	}

	return nil
}
