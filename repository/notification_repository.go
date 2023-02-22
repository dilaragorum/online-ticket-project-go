package repository

import (
	"context"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

type NotificationRepository interface {
	Create(ctx context.Context, channel model.Channel, log string) error
}

type notificationRepository struct {
	database *gorm.DB
}

func NewNotificationRepository(database *gorm.DB) *notificationRepository {
	return &notificationRepository{database: database}
}

func (m *notificationRepository) Create(ctx context.Context, channel model.Channel, logMsg string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	logMessage := &model.NotificationLog{
		Channel: channel,
		Log:     logMsg,
	}

	if err := m.database.WithContext(timeoutCtx).Model(&model.NotificationLog{}).Create(logMessage).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}
