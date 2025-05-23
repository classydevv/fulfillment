package entity

import "time"

type Provider struct {
	ProviderID ProviderID `db:"provider_id"`
	Name       string     `db:"name"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type ProviderID string
