package notification

import (
	"context"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(ctx context.Context, channel Channel, log string) error
}

type repository struct {
	database *gorm.DB
}

func NewNotificationRepository(database *gorm.DB) Repository {
	return &repository{database: database}
}

func (m *repository) Create(ctx context.Context, channel Channel, logMsg string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	logMessage := &Log{
		Channel: channel,
		Log:     logMsg,
	}

	if err := m.database.WithContext(timeoutCtx).Model(&Log{}).Create(logMessage).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
