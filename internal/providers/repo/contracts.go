package repo

import (
	"context"

	"github.com/classydevv/fulfillment/internal/providers/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	ProviderRepo interface {
		Store(context.Context, entity.Provider) error
		GetAll(context.Context) ([]entity.Provider, error)
		Update(context.Context, entity.ProviderId, entity.Provider) (*entity.Provider, error)
		Delete(context.Context, entity.ProviderId) error
	}
)
