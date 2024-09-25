package repo

import (
	"context"
	"github.com/magmaheat/auth_service/internal/entity"
	"github.com/magmaheat/auth_service/internal/repo/pgdb"
	"github.com/magmaheat/auth_service/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
}

type Repositories struct {
	User
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(pg),
	}
}
