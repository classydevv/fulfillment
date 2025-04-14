package persistent

import (
	"context"

	"github.com/classydevv/fulfillment/internal/entity"
)

type ProviderPostgres struct {
	//
}

func New() *ProviderPostgres {
	return &ProviderPostgres{}
}

func (pg *ProviderPostgres) Store(ctx context.Context, provider entity.Provider) error {
	return nil
}

func (pg *ProviderPostgres) GetAll(ctx context.Context) ([]entity.Provider, error) {
	return []entity.Provider{}, nil
}