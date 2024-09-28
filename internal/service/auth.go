package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/internal/service/sender"
	log "github.com/sirupsen/logrus"
	"time"
)

type Auth interface {
	GenerateTokens(userId, userIp string) (string, string, error)
	UpdateTokens(accessToken, refreshToken, userIp string) (string, string, error)
}

type AuthService struct {
	TokenRepo       repo.Token
	Sender          *sender.Sender
	TokenManager    Manager
	SignKey         string
	TokenAccessTTL  time.Duration
	TokenRefreshTTL time.Duration
}

func NewAuthService(tokenRepo repo.Token, tokenManager Manager, signKey string, tokenAccessTTL, tokenRefreshTTL time.Duration, snd *sender.Sender) *AuthService {
	return &AuthService{
		TokenRepo:       tokenRepo,
		Sender:          snd,
		TokenManager:    tokenManager,
		SignKey:         signKey,
		TokenAccessTTL:  tokenAccessTTL,
		TokenRefreshTTL: tokenRefreshTTL,
	}
}

func (a *AuthService) GenerateTokens(userId, userIp string) (string, string, error) {
	tokenId := uuid.New().String()

	accessToken, err := a.TokenManager.GenerateAccess(GenerateAccessInput{
		TokenId: tokenId,
		SignKey: a.SignKey,
		Expiry:  a.TokenAccessTTL,
	})

	if err != nil {
		log.Errorf("service - auth - GenerateTokens.GenerateAccess: %v", err)
		return "", "", fmt.Errorf("error generate tokens")
	}

	refreshToken, err := a.TokenManager.GenerateRefresh(GenerateRefreshInput{
		UserId:  userId,
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
	err := a.TokenManager.Validate(ValidateInput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserIp:       userIp,
		SignKey:      a.SignKey,
	})

	if err != nil {
		if errors.Is(err, ErrMismatchIP) {
			_ = a.Sender.SendEmail(context.Background(), "")

			userId, _, _ := a.TokenManager.GetUserIdAndTokenId(refreshToken, a.SignKey)
			log.Errorf("UpdateTokens - GetUserIdAndTokenId: %v", err)

			err = a.TokenRepo.DeactivateAllTokens(userId)
			if err != nil {
				log.Errorf("UpdateTokens - DeactivateAllTokens: %v", err)
			}

			log.Info("success deactivate all tokens")
		}

		log.Errorf("service - Auth - UpdateTokens.Validate: %v", err)
		return "", "", fmt.Errorf("no valid token")
	}

	userId, tokenId, _ := a.TokenManager.GetUserIdAndTokenId(refreshToken, a.SignKey)
	log.Infof("userId: %s", userId)
	state, err := a.TokenRepo.GetStateToken(tokenId)
	if err != nil {
		log.Errorf("UpdateTokens - GetToken: %v", err)
		return "", "", fmt.Errorf("error get token")
	}

	if state == "no" {
		_ = a.Sender.SendEmail(context.Background(), "")

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
