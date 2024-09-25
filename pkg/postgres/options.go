package postgres

import "time"

type Options func(p *Postgres)

func MaxPoolSize(size int) Options {
	return func(p *Postgres) {
		p.maxPoolSize = size
	}
}

func ConnTimeout(timeout time.Duration) Options {
	return func(p *Postgres) {
		p.connTimeout = timeout
	}
}

func ConnAttempts(attempts int) Options {
	return func(p *Postgres) {
		p.connAttempts = attempts
	}
}
