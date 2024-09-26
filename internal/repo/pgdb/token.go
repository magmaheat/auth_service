package pgdb

import (
	"github.com/magmaheat/auth_service/internal/entity"
	"github.com/magmaheat/auth_service/pkg/postgres"
)

type TokenRepo struct {
	*postgres.Postgres
}

func NewTokenRepo(pg *postgres.Postgres) *TokenRepo {
	return &TokenRepo{Postgres: pg}
}

func (t *TokenRepo) CreateToken(id, token string) error {

}

func (t *TokenRepo) GetAllTokens(id string) ([]entity.Token, error) {

}

func (t *TokenRepo) DeactivateAllTokens(id string) error {

}
