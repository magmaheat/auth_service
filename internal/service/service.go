package service

import (
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/internal/service/sender"
	"time"
)

type ServicesDependencies struct {
	Repos           *repo.Repositories
	TokenManager    Manager
	SignKey         string
	TokenAccessTTL  time.Duration
	TokenRefreshTTL time.Duration
	Sender          *sender.Sender
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
		Auth: NewAuthService(deps.Repos.Token, deps.TokenManager, deps.SignKey, deps.TokenAccessTTL, deps.TokenRefreshTTL, deps.Sender),
	}
}
