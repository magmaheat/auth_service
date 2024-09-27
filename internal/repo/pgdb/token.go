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

func (t *TokenRepo) CreateToken(idUser, idToken string) error {
	return nil
}

func (t *TokenRepo) GetToken(idToken string) (entity.Token, error) {
	return entity.Token{}, nil
}

func (t *TokenRepo) DeactivateToken(idToken string) error {
	return nil
}
