package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/magmaheat/auth_service/internal/repo"
	"github.com/magmaheat/auth_service/pkg/token"
	log "github.com/sirupsen/logrus"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserGUID string
	UserIp   string
	TokenId  string
}

type Auth interface {
	GenerateTokens(id, ip string) (string, string, error)
}

type AuthService struct {
	UserRepo        repo.User
	Token           token.ServiceToken
	SignKey         string
	TokenAccessTTL  time.Duration
	TokenRefreshTTL time.Duration
}

func NewAuthService(userRepo repo.User, token token.ServiceToken, signKey string, tokenAccessTTL, tokenRefreshTTL time.Duration) *AuthService {
	return &AuthService{
		UserRepo:        userRepo,
		Token:           token,
		SignKey:         signKey,
		TokenAccessTTL:  tokenAccessTTL,
		TokenRefreshTTL: tokenRefreshTTL,
	}
}

func (a *AuthService) GenerateTokens(id, ip string) (string, string, error) {
	tokenId := uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(a.TokenAccessTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserIp:   ip,
		UserGUID: id,
		TokenId:  tokenId,
	})

	accessTokenString, err := accessToken.SignedString([]byte(a.SignKey))
	if err != nil {
		log.Errorf("AuthService.GenerateToken: cannot sing hasher: %v", err)
		return "", "", fmt.Errorf("error generate token: %v", err)
	}

	input := token.GenerateInput{
		Id:      id,
		Ip:      ip,
		SignKey: a.SignKey,
		TokenId: tokenId,
		Expiry:  a.TokenRefreshTTL,
	}
	refreshTokenString, err := a.Token.Generate(input)

	//TODO CreateToken(id)=

	return accessTokenString, refreshTokenString, nil
}

func (a *AuthService) UpdateTokens(ip, refreshToken, token string) (*AuthUpdateTokens, error) {
	//err := a.Token.Validate(refreshToken) // TODO middleware validate

}
