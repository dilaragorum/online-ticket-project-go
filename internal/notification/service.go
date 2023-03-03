package notification

import "context"

type Service interface {
	Send(ctx context.Context, param Param) error
}

type Notifier interface {
	Send(ctx context.Context, param Param) error
}

type defaultService struct {
	repo        Repository
	strategyMap map[Channel]Notifier
}

func NewService(r Repository) Service {
	return &defaultService{
		repo: r,
		strategyMap: map[Channel]Notifier{
			SMS:   NewSms(),
			Email: NewMail(),
		},
	}
}

func (d *defaultService) Send(ctx context.Context, param Param) error {
	if err := d.strategyMap[param.Channel].Send(ctx, param); err != nil {
		return err
	}

	return d.repo.Create(ctx, param.Channel, param.LogMsg)
}
