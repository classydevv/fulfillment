package usecase_test

import (
	"context"
	"testing"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	mock_repo "github.com/classydevv/fulfillment/internal/providers/repo/mocks"
	"github.com/classydevv/fulfillment/internal/providers/usecase"
	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

func TestUseCaseProviders_Create(t *testing.T) {
	type fields struct {
		repo *mock_repo.MockProviderRepo
	}
	type args struct {
		ctx      context.Context
		provider *entity.Provider
	}
	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    entity.ProviderId
		wantErr error
	}{
		{
			name: "provider created successfully",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderId: entity.ProviderId("id"), Name: "name"}).Return(nil)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderId: entity.ProviderId("id"), Name: "name"}},
			want:    entity.ProviderId("id"),
			wantErr: nil,
		},
		{
			name: "error - missing required field",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderId: entity.ProviderId("id")}).Return(entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderId: entity.ProviderId("id")}},
			want:    entity.ProviderId(""),
			wantErr: entity.ErrInternalServerError,
		},
		{
			name: "error - entity already exists",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderId: entity.ProviderId("id")}).Return(entity.ErrAlreadyExists)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderId: entity.ProviderId("id")}},
			want:    entity.ProviderId(""),
			wantErr: entity.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := fields{
				repo: mock_repo.NewMockProviderRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.NewUseCaseProviders(f.repo)

			res, err := uc.Create(tt.args.ctx, tt.args.provider)

			require.Equal(t, res, tt.want)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
