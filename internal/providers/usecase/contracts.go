package usecase

import (
	"context"

	"github.com/classydevv/fulfillment/internal/providers/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks/mock_usecase.go

type (
	Provider interface {
		Create(context.Context, *entity.Provider) (entity.ProviderID, error)
		ListAll(context.Context) ([]*entity.Provider, error)
		Update(context.Context, entity.ProviderID, *entity.Provider) (*entity.Provider, error)
		Delete(context.Context, entity.ProviderID) error
	}
)
