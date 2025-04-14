package usecase

import (
	"context"

	"github.com/classydevv/fulfillment/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	Provider interface {
		Save(context.Context, entity.Provider) (entity.ProviderId, error)
		ListAll(context.Context) ([]entity.Provider, error)
	}
)