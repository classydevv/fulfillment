package usecase

import (
	"context"

	"github.com/classydevv/fulfillment/internal/providers/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	Provider interface {
		Create(context.Context, entity.Provider) (entity.ProviderId, error)
		ListAll(context.Context) ([]entity.Provider, error)
		Update(context.Context, entity.ProviderId, entity.Provider) (*entity.Provider, error)
		Delete(context.Context, entity.ProviderId) error
	}
)
