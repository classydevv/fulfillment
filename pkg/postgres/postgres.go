package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/classydevv/fulfillment/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts     = 10
	_defaultConnRetryTimeout = time.Second
)

type Postgres struct {
	maxPoolSize      int32
	connAttempts     int
	connRetryTimeout time.Duration

	Pool    *pgxpool.Pool
	Builder squirrel.StatementBuilderType
}

func New(url string, l logger.Interface, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		connAttempts:     _defaultConnAttempts,
		connRetryTimeout: _defaultConnRetryTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(pg)
	}

	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.ParseConfig: %w", err)
	}

	if pg.maxPoolSize != 0 {
		poolConfig.MaxConns = pg.maxPoolSize
	}

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		l.Info("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connRetryTimeout)
		pg.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - New - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
