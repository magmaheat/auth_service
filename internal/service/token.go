package service

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type Manager interface {
	GenerateRefresh(input GenerateRefreshInput) (string, error)
	GenerateAccess(input GenerateAccessInput) (string, error)
	GetUserIdAndTokenId(tkn, signKey string) (string, string, error)
	Validate(input ValidateInput) error
}

type Base64URL struct{}

type GenerateRefreshInput struct {
	UserIp  string
	UserId  string
	TokenId string
	SignKey string
	Expiry  time.Duration
}

type GenerateAccessInput struct {
	TokenId string
	SignKey string
	Expiry  time.Duration
}

type ValidateInput struct {
	AccessToken  string
	RefreshToken string
	UserIp       string
	SignKey      string
}

type PersonClaims struct {
	jwt.StandardClaims
	UserIp  string
	UserId  string
	TokenId string
}

func NewBase64URL() *Base64URL {
	return &Base64URL{}
}

func (b *Base64URL) GenerateRefresh(input GenerateRefreshInput) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS512, &PersonClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(input.Expiry).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserIp:  input.UserIp,
		UserId:  input.UserId,
		TokenId: input.TokenId,
	})

	tokenString, err := tkn.SignedString([]byte(input.SignKey))
	if err != nil {
		return "", fmt.Errorf("error signed refresh token: %v", err)
	}

	return tokenString, nil
}

func (b *Base64URL) GenerateAccess(input GenerateAccessInput) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS512, &PersonClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(input.Expiry).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenId: input.TokenId,
	})

	tokenString, err := tkn.SignedString([]byte(input.SignKey))
	if err != nil {
		return "", fmt.Errorf("error signed access token: %v", err)
	}

	return tokenString, nil
}

func (b *Base64URL) Validate(input ValidateInput) error {
	accessToken, err := decode(input.AccessToken, input.SignKey)
	if err != nil {
		return err
	}

	refreshToken, err := decode(input.RefreshToken, input.SignKey)

	if err != nil {
		return err
	}

	accessClaims, ok := accessToken.Claims.(*PersonClaims)
	if !ok || !accessToken.Valid {
		return fmt.Errorf("access token not valid")
	}

	refreshClaims, ok := refreshToken.Claims.(*PersonClaims)
	if !ok || !refreshToken.Valid {
		return fmt.Errorf("refresh token not valid")
	}

	if input.UserIp != refreshClaims.UserIp {
		return ErrMismatchIP
	}

	if accessClaims.TokenId != refreshClaims.TokenId {
		return fmt.Errorf("ID token mismatch")
	}

	return nil
}

func (b *Base64URL) GetUserIdAndTokenId(tkn, signKey string) (string, string, error) {
	token, err := decode(tkn, signKey)

	if err != nil {
		return "", "", err
	}

	if tokenClaims, ok := token.Claims.(*PersonClaims); ok && token.Valid {
		return tokenClaims.UserId, tokenClaims.TokenId, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

func decode(tkn, signKey string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tkn, &PersonClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parse claims: %v", err)
	}

	return token, nil
}
