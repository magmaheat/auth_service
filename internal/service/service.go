package service

import (
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/pkg/token"
	"time"
)

type ServicesDependencies struct {
	Repos           *repo.Repositories
	Token           token.ServiceToken
	SignKey         string
	TokenAccessTTL  time.Duration
	TokenRefreshTTL time.Duration
}

type AuthUpdateTokens struct {
	AccessToken  string
	RefreshToken string
}

type Services struct {
	Auth
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repos.Token, deps.Token, deps.SignKey, deps.TokenAccessTTL, deps.TokenRefreshTTL),
	}
}
