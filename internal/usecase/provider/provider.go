package provider

import (
	"context"
	"fmt"

	"github.com/classydevv/fulfillment/internal/entity"
	"github.com/classydevv/fulfillment/internal/repo"
)

type UseCase struct {
	repo repo.ProviderRepo
}

func New(r repo.ProviderRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Save(ctx context.Context, provider entity.Provider) (entity.ProviderId, error) {
	err := uc.repo.Store(ctx, provider)
	if err != nil {
		return "", fmt.Errorf("UseCaseProvider - Save - uc.repo.Store: %w", err)
	}

	return provider.Id, nil
}

func (uc *UseCase) ListAll(ctx context.Context) ([]entity.Provider, error) {
	providers, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("UseCaseProvider - ListAll - uc.repo.GetAll: %w", err)
	}

	return providers, nil
}