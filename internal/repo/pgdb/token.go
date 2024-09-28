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

	_, err := t.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		log.Errorf("pgdb - CreateToeken.QueryRow: %v", err)
		return fmt.Errorf("error create token")
	}

	return nil
}

func (t *TokenRepo) GetStateToken(tokenId string) (string, error) {
	ctx := context.Background()

	tx, err := t.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("pgdb - Begin transaction: %v", err)
		return "", fmt.Errorf("error beginning transaction")
	}
	defer tx.Rollback(ctx) // Откат транзакции в случае ошибки

	sql, args, _ := t.Builder.
		Select("valid").
		From("tokens").
		Where(squirrel.Eq{"token_id": tokenId}).
		ToSql()

	var valid string
	err = tx.QueryRow(ctx, sql, args...).Scan(&valid)
	if err != nil {
		log.Errorf("pgdb - GetAndInvalidateToken.QueryRow: %v", err)
		return "", fmt.Errorf("error getting token state")
	}

	sql, args, _ = t.Builder.
		Update("tokens").
		Set("valid", "no").
		Where(squirrel.Eq{"token_id": tokenId}).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("pgdb - GetAndInvalidateToken.Update: %v", err)
		return "", fmt.Errorf("error updating token state")
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("pgdb - Commit transaction: %v", err)
		return "", fmt.Errorf("error committing transaction")
	}

	return valid, nil
}

func (t *TokenRepo) DeactivateAllTokens(userId string) error {
	sql, args, _ := t.Builder.
		Update("tokens").
		Set("valid", "no").
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()

	_, err := t.Pool.Exec(context.Background(), sql, args...)
	if err != nil {
		log.Errorf("pgdb - DeactivateAllTokens.QueryRow: %v", err)
		return fmt.Errorf("error deactivate all tokens")
	}

	return nil
}
