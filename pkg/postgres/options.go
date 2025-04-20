package postgres

type Option func(*Postgres)

func MaxPoolSize(size int32) Option {
	return func(c *Postgres) {
		c.maxPoolSize = size
	}
}
