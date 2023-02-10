package service

import (
	"context"
	"errors"
	"github.com/dilaragorum/online-ticket-project-go/model"
	"github.com/dilaragorum/online-ticket-project-go/repository"
)

var (
	ErrDuplicatedUserName = errors.New("username should be unique")
	ErrDuplicatedEmail    = errors.New("email should be unique")
)

type Service interface {
	Register(ctx context.Context, user *model.User) (*model.User, error)
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
