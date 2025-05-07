package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/classydevv/fulfillment/internal/providers/entity"
	"github.com/classydevv/fulfillment/pkg/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepo struct {
	*postgres.Postgres
}

func NewPostgresRepo(pg *postgres.Postgres) *PostgresRepo {
	return &PostgresRepo{pg}
}

func (pg *PostgresRepo) Store(ctx context.Context, provider entity.Provider) error {
	query, args, err := pg.Builder.
		Insert("providers").
		Columns("id, name").
		Values(provider.Id, provider.Name).
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

func (pg *PostgresRepo) GetAll(ctx context.Context) ([]entity.Provider, error) {
	query, _, err := pg.Builder.
		Select("id, name").
		From("providers").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - pg.Builder: %w", err)
	}

	rows, err := pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	providers := make([]entity.Provider, 0)

	for rows.Next() {
		provider := entity.Provider{}
		err := rows.Scan(&provider.Id, &provider.Name)
		if err != nil {
			return nil, fmt.Errorf("PostgresRepo - GetAll - rows.Scan: %w", err)
		}
		providers = append(providers, provider)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("PostgresRepo - GetAll - rows.Err: %w", err)
	}

	return providers, nil
}
