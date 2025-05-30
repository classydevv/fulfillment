package repo

import (
	"context"

	"github.com/classydevv/fulfillment/internal/providers/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks/mock_repo.go

type (
	ProviderRepo interface {
		Store(context.Context, *entity.Provider) error
		GetAll(context.Context) ([]*entity.Provider, error)
		Update(context.Context, entity.ProviderID, *entity.Provider) (*entity.Provider, error)
		Delete(context.Context, entity.ProviderID) error
	}
)
