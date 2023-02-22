package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicatedEmail    = errors.New(`ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)`)
	ErrDuplicatedUserName = errors.New(`ERROR: duplicate key value violates unique constraint "users_user_name_key" (SQLSTATE 23505)`)
	ErrNoRecord           = errors.New("there is no record in DB with that username")
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByUserName(ctx context.Context, username string) (*model.User, error)
}

type defaultUserRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) *defaultUserRepository {
	return &defaultUserRepository{
		database: database,
	}
}

func (r *defaultUserRepository) GetByUserName(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&model.User{}).
		First(&user, "user_name = ?", username).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
	}

	return &user, nil
}

func (r *defaultUserRepository) Create(ctx context.Context, user *model.User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&model.User{}).Create(user).Error; err != nil {
		switch {
		case errors.Is(err, ErrDuplicatedEmail):
			return ErrDuplicatedEmail
		case errors.Is(err, ErrDuplicatedUserName):
			return ErrDuplicatedUserName
		default:
			return err
		}
	}

	return nil
}
