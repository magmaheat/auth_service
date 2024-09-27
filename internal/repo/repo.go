package repo

import (
	"github.com/magmaheat/auth_service/internal/repo/pgdb"
	"github.com/magmaheat/auth_service/pkg/postgres"
)

type Token interface {
	CreateToken(userId, tokenId string) error
	GetStateToken(tokenId string) (string, error)
	DeactivateAllTokens(userId string) error
}

type Repositories struct {
	Token
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Token: pgdb.NewTokenRepo(pg),
	}
}
