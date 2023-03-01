package user

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/internal/aut"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicatedValue           = errors.New("username or email should be unique")
	ErrUsernameNotFound          = errors.New("there is no that username in record")
	ErrUsernameOrPasswordInvalid = errors.New("invalid username or password")

	ErrThereIsNoTrip = errors.New("there is no trip which meet these conditions")
)

type Service interface {
	Register(ctx context.Context, user *User) error
	Login(ctx context.Context, credentials aut.Credentials) (*User, error)
}

type service struct {
	userRepo Repository
}

func NewUserService(repository Repository) *service {
	return &service{userRepo: repository}
}

func (s *service) Register(ctx context.Context, user *User) error {
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, ErrUniqueViolation):
			return ErrDuplicatedValue
		default:
			return err
		}
	}
	return nil
}

func (s *service) Login(ctx context.Context, credentials aut.Credentials) (*User, error) {
	user, err := s.userRepo.GetByUserName(ctx, credentials.UserName)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			return nil, ErrUsernameNotFound
		}
		return nil, err
	}

	if s.isNotEqualHashAndPassword(user.Password, credentials.Password) {
		return nil, ErrUsernameOrPasswordInvalid
	}

	return user, nil
}

func (s *service) isNotEqualHashAndPassword(hashPassword string, password string) bool {
	return !s.isEqualHashAndPassword(hashPassword, password)
}

func (s *service) isEqualHashAndPassword(hashPassword string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
