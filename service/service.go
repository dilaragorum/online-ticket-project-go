package service

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
)

var (
	ErrDuplicatedUserName        = errors.New("username should be unique")
	ErrDuplicatedEmail           = errors.New("email should be unique")
	ErrUsernameNotFound          = errors.New("there is no that username in record")
	ErrUsernameOrPasswordInvalid = errors.New("invalid username or password")
)

type Service interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
	LogIn(ctx context.Context, credentials model.Credentials) (*model.User, error)
}

type DefaultService struct {
	repository repository.Repository
}

func NewDefaultService(repository repository.Repository) *DefaultService {
	return &DefaultService{repository: repository}
}

func (s *DefaultService) Register(ctx context.Context, user *model.User) (*model.User, error) {
	register, err := s.repository.Register(ctx, user)
	if err != nil {
		switch err.Error() {
		case repository.ErrDBDuplicatedEmail.Error():
			return nil, ErrDuplicatedEmail
		case repository.ErrDBDuplicatedUserName.Error():
			return nil, ErrDuplicatedUserName
		default:
			return nil, err
		}
	}

	return register, nil
}

func (s *DefaultService) LogIn(ctx context.Context, credentials model.Credentials) (*model.User, error) {
	user, err := s.repository.FindUser(ctx, credentials.UserName)
	if err != nil {
		if err.Error() == repository.ErrDBNoRecord.Error() {
			return nil, ErrUsernameNotFound
		}
		return nil, err
	}

	if credentials.UserName != user.UserName || credentials.Password != user.Password {
		return nil, ErrUsernameOrPasswordInvalid
	}

	return user, nil
}
