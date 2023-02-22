package service

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedUserName        = errors.New("username should be unique")
	ErrDuplicatedEmail           = errors.New("email should be unique")
	ErrUsernameNotFound          = errors.New("there is no that username in record")
	ErrUsernameOrPasswordInvalid = errors.New("invalid username or password")
)

type UserService interface {
	Register(ctx context.Context, user *model.User) error
	Login(ctx context.Context, credentials model.Credentials) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{userRepo: repository}
}

func (s *userService) Register(ctx context.Context, user *model.User) error {
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicatedEmail):
			return ErrDuplicatedEmail
		case errors.Is(err, repository.ErrDuplicatedUserName):
			return ErrDuplicatedUserName
		default:
			return err
		}
	}
	return nil
}

func (s *userService) Login(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	user, err := s.userRepo.GetByUserName(ctx, credentials.UserName)
	if err != nil {
		if errors.Is(err, repository.ErrNoRecord) {
			return nil, ErrUsernameNotFound
		}
		return nil, err
	}

	if s.isNotEqualHashAndPassword(user.Password, credentials.Password) {
		return nil, ErrUsernameOrPasswordInvalid
	}

	return user, nil
}

func (s *userService) isNotEqualHashAndPassword(hashPassword string, password string) bool {
	return !s.isEqualHashAndPassword(hashPassword, password)
}

func (s *userService) isEqualHashAndPassword(hashPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
