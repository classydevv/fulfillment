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

func (uc *UseCaseProviders) Create(ctx context.Context, provider *entity.Provider) (entity.ProviderID, error) {
	err := uc.repo.Store(ctx, provider)
	if err != nil {
		return "", fmt.Errorf("UseCaseProviders - Save - uc.repo.Store: %w", err)
	}
	return provider.ProviderID, nil
}

func (uc *UseCaseProviders) ListAll(ctx context.Context) ([]*entity.Provider, error) {
	providers, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("UseCaseProviders - ListAll - uc.repo.GetAll: %w", err)
	}

	return providers, nil
}

func (uc *UseCaseProviders) Update(ctx context.Context, providerID entity.ProviderID, provider *entity.Provider) (*entity.Provider, error) {
	providerUpdated, err := uc.repo.Update(ctx, providerID, provider)
	if err != nil {
		return nil, fmt.Errorf("UseCaseProviders - Update - uc.repo.Update: %w", err)
	}

	return providerUpdated, nil
}

func (uc *UseCaseProviders) Delete(ctx context.Context, providerID entity.ProviderID) error {
	err := uc.repo.Delete(ctx, providerID)
	if err != nil {
		return fmt.Errorf("UseCaseProviders - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
