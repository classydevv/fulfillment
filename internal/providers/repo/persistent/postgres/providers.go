package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepo struct {
	*postgres.Postgres
}

func NewPostgresRepo(pg *postgres.Postgres) *PostgresRepo {
	return &PostgresRepo{pg}
}

func (pg *PostgresRepo) Store(ctx context.Context, p *entity.Provider) error {
	query, args, err := pg.Builder.
		Insert("providers").
		Columns("provider_id, name").
		Values(p.ProviderID, p.Name).
		ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - Store - pg.Builder: %w", err)
	}

	_, err = pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("PostgresRepo - Store - pg.Pool.Exec: %w", entity.ErrAlreadyExists)
		}
		return fmt.Errorf("PostgresRepo - Store - pg.Pool.Exec: %w", err)
	}

	return nil
}

func (pg *PostgresRepo) GetAll(ctx context.Context) ([]*entity.Provider, error) {
	query, _, err := pg.Builder.
		Select("*").
		From("providers").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - pg.Builder: %w", err)
	}

	rows, err := pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - pg.Pool.Query: %w", err)
	}

	providers, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[entity.Provider])
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - pgx.CollectRows: %w", err)
	}

	return providers, nil
}

func (pg *PostgresRepo) Update(ctx context.Context, id entity.ProviderID, p *entity.Provider) (*entity.Provider, error) {
	query, args, err := pg.Builder.
		Update("providers").
		Set(
			"name", p.Name,
		).
		Where("provider_id = ?", id).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - Update - pg.Builder: %w", err)
	}

	rows, err := pg.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - Update - pg.Pool.Query: %w", err)
	}
	provider, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[entity.Provider])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("PostgresRepo - Update - pgx.CollectOneRow: %w", entity.ErrNotFound)
		}
		return nil, fmt.Errorf("PostgresRepo - Update - pgx.CollectOneRow: %w", err)
	}

	return provider, nil
}

func (pg *PostgresRepo) Delete(ctx context.Context, id entity.ProviderID) error {
	query, args, err := pg.Builder.
		Delete("providers").
		Where("provider_id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("PostgresRepo - Delete - pg.Builder: %w", err)
	}

	comm, err := pg.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("PostgresRepo - Delete - pg.Pool.Exec: %w", err)
	}

	if comm.RowsAffected() != 1 {
		return fmt.Errorf("PostgresRepo - Delete - pg.Pool.Exec: %w", entity.ErrNotFound)
	}

	return nil
}
