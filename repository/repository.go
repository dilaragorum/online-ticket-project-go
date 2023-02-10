package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDBDuplicatedEmail    = errors.New(`ERROR: duplicate key value violates unique constraint "idx_users_email" (SQLSTATE 23505)`)
	ErrDBDuplicatedUserName = errors.New(`ERROR: duplicate key value violates unique constraint "users_user_name_key" (SQLSTATE 23505)`)
	ErrDBNoRecord           = errors.New("there is no record in DB with that username")
)

type Repository interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	FindUser(ctx context.Context, username string) (*model.User, error)
}

type DefaultRepository struct {
	database *gorm.DB
}

func NewDefaultRepository(database *gorm.DB) *DefaultRepository {
	return &DefaultRepository{
		database: database,
	}
}

func (r *DefaultRepository) FindUser(ctx context.Context, username string) (*model.User, error) {
	user := model.User{}

	timeoutCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	if err := r.database.WithContext(timeoutCtx).Model(&user).First(&user, "user_name = ?", username).Error; err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDBNoRecord
		}
	}

	return &user, nil
}

func (r *DefaultRepository) Register(ctx context.Context, user *model.User) (*model.User, error) {
	userModel := model.User{
		UserName: user.UserName,
		Password: user.Password,
		UserType: user.UserType,
		Email:    user.Email,
		Model:    gorm.Model{},
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := r.database.WithContext(timeoutCtx).Model(&userModel).Create(&userModel).Error
	fmt.Println("---------------")
	fmt.Println(err.Error())
	fmt.Println(ErrDBDuplicatedEmail.Error())
	fmt.Println(ErrDBDuplicatedUserName.Error())

	if err != nil {
		switch err.Error() {
		case ErrDBDuplicatedEmail.Error():
			return nil, ErrDBDuplicatedEmail
		case ErrDBDuplicatedUserName.Error():
			return nil, ErrDBDuplicatedUserName
		default:
			return nil, err
		}
	}

	return &userModel, nil
}
