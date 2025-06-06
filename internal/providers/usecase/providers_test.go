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
	t.Parallel()
	
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
		want    entity.ProviderID
		wantErr error
	}{
		{
			name: "provider created successfully",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"}).Return(nil)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"}},
			want:    entity.ProviderID("id"),
			wantErr: nil,
		},
		{
			name: "error - missing required field",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderID: entity.ProviderID("id")}).Return(entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderID: entity.ProviderID("id")}},
			want:    entity.ProviderID(""),
			wantErr: entity.ErrInternalServerError,
		},
		{
			name: "error - entity already exists",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderID: entity.ProviderID("id")}).Return(entity.ErrAlreadyExists)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderID: entity.ProviderID("id")}},
			want:    entity.ProviderID(""),
			wantErr: entity.ErrAlreadyExists,
		},
		{
			name: "error - database not available",
			prepare: func(f *fields) {
				f.repo.EXPECT().Store(context.Background(), &entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"}).Return(entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background(), provider: &entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"}},
			want:    entity.ProviderID(""),
			wantErr: entity.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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

func TestUseCaseProviders_ListAll(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo *mock_repo.MockProviderRepo
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    []*entity.Provider
		wantErr error
	}{
		{
			name: "providers listed successfully",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAll(context.Background()).Return([]*entity.Provider{}, nil)
			},
			args:    args{ctx: context.Background()},
			want:    []*entity.Provider{},
			wantErr: nil,
		},
		{
			name: "error - database not available",
			prepare: func(f *fields) {
				f.repo.EXPECT().GetAll(context.Background()).Return(nil, entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: entity.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: mock_repo.NewMockProviderRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.NewUseCaseProviders(f.repo)

			res, err := uc.ListAll(tt.args.ctx)

			require.Equal(t, tt.want, res)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUseCaseProviders_Update(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo *mock_repo.MockProviderRepo
	}

	type args struct {
		ctx      context.Context
		id       entity.ProviderID
		provider *entity.Provider
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		want    *entity.Provider
		wantErr error
	}{
		{
			name: "provider updated successfully",
			prepare: func(f *fields) {
				f.repo.EXPECT().Update(context.Background(), entity.ProviderID("id"), &entity.Provider{Name: "name"}).Return(&entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"}, nil)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id"), provider: &entity.Provider{Name: "name"}},
			want:    &entity.Provider{ProviderID: entity.ProviderID("id"), Name: "name"},
			wantErr: nil,
		},
		{
			name: "error - provider not found",
			prepare: func(f *fields) {
				f.repo.EXPECT().Update(context.Background(), entity.ProviderID("id"), &entity.Provider{Name: "name"}).Return(nil, entity.ErrNotFound)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id"), provider: &entity.Provider{Name: "name"}},
			want:    nil,
			wantErr: entity.ErrNotFound,
		},
		{
			name: "error - database not available",
			prepare: func(f *fields) {
				f.repo.EXPECT().Update(context.Background(), entity.ProviderID("id"), &entity.Provider{Name: "name"}).Return(nil, entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id"), provider: &entity.Provider{Name: "name"}},
			want:    nil,
			wantErr: entity.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: mock_repo.NewMockProviderRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.NewUseCaseProviders(f.repo)

			res, err := uc.Update(tt.args.ctx, tt.args.id, tt.args.provider)

			require.Equal(t, tt.want, res)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUseCaseProviders_Delete(t *testing.T) {
	t.Parallel()

	type fields struct {
		repo *mock_repo.MockProviderRepo
	}

	type args struct {
		ctx context.Context
		id  entity.ProviderID
	}

	tests := []struct {
		name    string
		prepare func(f *fields)
		args    args
		wantErr error
	}{
		{
			name: "provider deleted successfully",
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(context.Background(), entity.ProviderID("id")).Return(nil)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id")},
			wantErr: nil,
		},
		{
			name: "error - provider not found",
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(context.Background(), entity.ProviderID("id")).Return(entity.ErrNotFound)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id")},
			wantErr: entity.ErrNotFound,
		},
		{
			name: "error - database not available",
			prepare: func(f *fields) {
				f.repo.EXPECT().Delete(context.Background(), entity.ProviderID("id")).Return(entity.ErrInternalServerError)
			},
			args:    args{ctx: context.Background(), id: entity.ProviderID("id")},
			wantErr: entity.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				repo: mock_repo.NewMockProviderRepo(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(&f)
			}

			uc := usecase.NewUseCaseProviders(f.repo)

			err := uc.Delete(tt.args.ctx, tt.args.id)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
