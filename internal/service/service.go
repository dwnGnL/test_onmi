package service

import (
	"context"

	"github.com/dwnGnL/testWork/internal/config"
	externalLib "github.com/dwnGnL/testWork/lib/external"

	"github.com/dwnGnL/testWork/internal/service/external"
)

type ServiceImpl struct {
	conf     *config.Config
	external ExternalService
}

type Auther interface {
	Login(ctx context.Context, username string, password string) (string, error)
	CheckToken(tokenStr string) (int64, error)
}

type ExternalService interface {
	ProcessBatch(objects externalLib.Item)
}

type Option func(*ServiceImpl)

func New(ctx context.Context, conf *config.Config, opts ...Option) *ServiceImpl {

	s := ServiceImpl{
		conf:     conf,
		external: external.New(ctx, conf),
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

func (s ServiceImpl) GetExternal() ExternalService {
	return s.external
}
