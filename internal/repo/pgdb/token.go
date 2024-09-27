package pgdb

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/magmaheat/auth_service/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

type TokenRepo struct {
	*postgres.Postgres
}

func NewTokenRepo(pg *postgres.Postgres) *TokenRepo {
	return &TokenRepo{Postgres: pg}
}

func (t *TokenRepo) CreateToken(userId, tokenId string) error {
	sql, args, _ := t.Builder.
		Insert("tokens").
		Columns("user_id", "token_id").
		Values(userId, tokenId).
		ToSql()

	err := t.Pool.QueryRow(context.Background(), sql, args...)
	if err != nil {
		log.Errorf("pgdb - CreateToeken.QueryRow: %v", err)
		return fmt.Errorf("error create token")
	}

	return nil
}

func (t *TokenRepo) GetStateToken(tokenId string) (string, error) {
	sql, args, _ := t.Builder.
		Select("valid").
		From("get_and_invalidate_token(?)").
		Where(squirrel.Eq{"token_id": tokenId}).
		ToSql()

	var valid string
	err := t.Pool.QueryRow(context.Background(), sql, args...).Scan(&valid)
	if err != nil {
		log.Errorf("pgdb - GetStateToken.QueryRow: %v", err)
		return "", fmt.Errorf("error get state token")
	}

	return valid, nil
}

func (t *TokenRepo) DeactivateAllTokens(userId string) error {
	return nil
}
