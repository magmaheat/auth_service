package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/pkg/token"
	log "github.com/sirupsen/logrus"
	"time"
)

type Auth interface {
	GenerateTokens(userId, userIp string) (string, string, error)
	UpdateTokens(accessToken, refreshToken, userIp string) (string, string, error)
}

type AuthService struct {
	TokenRepo       repo.Token
	TokenManager    token.Manager
	SignKey         string
	TokenAccessTTL  time.Duration
	TokenRefreshTTL time.Duration
}

func NewAuthService(userRepo repo.Token, tokenManager token.Manager, signKey string, tokenAccessTTL, tokenRefreshTTL time.Duration) *AuthService {
	return &AuthService{
		TokenRepo:       userRepo,
		TokenManager:    tokenManager,
		SignKey:         signKey,
		TokenAccessTTL:  tokenAccessTTL,
		TokenRefreshTTL: tokenRefreshTTL,
	}
}

func (a *AuthService) GenerateTokens(userId, userIp string) (string, string, error) {
	tokenId := uuid.New().String()

	accessToken, err := a.TokenManager.GenerateAccess(token.GenerateAccessInput{
		TokenId: tokenId,
		SignKey: a.SignKey,
		Expiry:  a.TokenAccessTTL,
	})

	if err != nil {
		log.Errorf("service - auth - GenerateTokens.GenerateAccess: %v", err)
		return "", "", fmt.Errorf("error generate tokens")
	}

	refreshToken, err := a.TokenManager.GenerateRefresh(token.GenerateRefreshInput{
		UserIp:  userIp,
		TokenId: tokenId,
		SignKey: a.SignKey,
		Expiry:  a.TokenRefreshTTL,
	})

	if err != nil {
		log.Errorf("service - auth - GenerateTokens.GenerateRefresh: %v", err)
		return "", "", fmt.Errorf("error generate tokens")
	}

	err = a.TokenRepo.CreateToken(userId, tokenId)
	if err != nil {
		log.Errorf("service - auth - GenerateTokens - CreateToken: %v", err)
		return "", "", fmt.Errorf("error create token")
	}

	return accessToken, refreshToken, nil
}

func (a *AuthService) UpdateTokens(accessToken, refreshToken, userIp string) (string, string, error) {
	err := a.TokenManager.Validate(token.ValidateInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserIp:       userIp,
		SignKey:      a.SignKey,
	})

	if err != nil {
		if errors.Is(err, token.ErrMismatchIP) {
			idUser, _, _ := a.TokenManager.GetUserIdAndTokenId(refreshToken, a.SignKey)
			log.Errorf("UpdateTokens - GetUserIdAndTokenId: %v", err)

			err = a.TokenRepo.DeactivateAllTokens(idUser)
			if err != nil {
				log.Errorf("UpdateTokens - DeactivateAllTokens: %v", err)
			}

			log.Info("success deactivate all tokens")
		}

		return "", "", fmt.Errorf("no valid token")
	}

	userId, tokenId, _ := a.TokenManager.GetUserIdAndTokenId(refreshToken, a.SignKey)
	state, err := a.TokenRepo.GetStateToken(tokenId)
	if err != nil {
		log.Errorf("UpdateTokens - GetToken: %v", err)
		return "", "", fmt.Errorf("error get token")
	}

	if state == "no" {
		err = a.TokenRepo.DeactivateAllTokens(userId)
		if err != nil {
			log.Errorf("UpdateTokens - DeactivateAllTokens: %v", err)
		}

		return "", "", fmt.Errorf("no valid token")
	}

	accessToken, refreshToken, err = a.GenerateTokens(userId, userIp)
	if err != nil {
		return "", "", fmt.Errorf("error update token")
	}

	return accessToken, refreshToken, nil
}
