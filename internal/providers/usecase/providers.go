package usecase

import (
	"context"
	"fmt"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/internal/providers/repo"
)

type UseCaseProviders struct {
	repo repo.ProviderRepo
}

func NewUseCaseProviders(r repo.ProviderRepo) *UseCaseProviders {
	return &UseCaseProviders{
		repo: r,
	}
}

func (uc *UseCaseProviders) Create(ctx context.Context, provider entity.Provider) (entity.ProviderId, error) {
	err := uc.repo.Store(ctx, provider)
	if err != nil {
		return "", fmt.Errorf("UseCaseProviders - Save - uc.repo.Store: %w", err)
	}

	return provider.Id, nil
}

func (uc *UseCaseProviders) ListAll(ctx context.Context) ([]entity.Provider, error) {
	providers, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("UseCaseProviders - ListAll - uc.repo.GetAll: %w", err)
	}

	return providers, nil
}
