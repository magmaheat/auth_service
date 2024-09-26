package repo

import (
	"github.com/magmaheat/auth_service/internal/entity"
	"github.com/magmaheat/auth_service/internal/repo/pgdb"
	"github.com/magmaheat/auth_service/pkg/postgres"
)

type Token interface {
	CreateToken(id, token string) error
	GetAllTokens(id string) ([]entity.Token, error)
	DeactivateToken(hashToken string) error
	DeactivateAllTokens(id string) error
}

type Repositories struct {
	Token
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		Token: pgdb.NewTokenRepo(pg),
	}
}
