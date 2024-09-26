package token

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type ServiceToken interface {
	Generate(input GenerateInput) (string, error)
	Decode(token string) (*Payload, error)
	Validate(token string) error
	Hash(token string) string
}

type Base64Token struct{}

func NewBase64Token () *Base64Token {
	return &Base64Token{}
}

type Payload struct {
	TokenId string `json:"token_id"`
	Id      string `json:"id"`
	Ip      string `json:"ip"`
	Expiry  int64  `json:"expiry"`
}

type GenerateInput struct {
	Ip      string
	Id      string
	SignKey string
	TokenId string
	Expiry  time.Duration
}

func (b *Base64Token) Generate(input GenerateInput) (string, error) {
	payload := Payload{
		Id:     input.Id,
		Ip:     input.Ip,
		Expiry: time.Now().Add(input.Expiry * time.Minute).Unix(),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("pkg - token - Generate: %v", err)
	}

	encodedToken := base64.StdEncoding.EncodeToString(jsonPayload)
	return encodedToken, nil
}

func (b *Base64Token) Decode(token string) (*Payload, error) {
	jsonPayload, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("pkg - token - Decode: %v", err)
	}

	var payload Payload
	err = json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		return nil, fmt.Errorf("pkg - token - Decode: %v", err)
	}

	return &payload, nil
}

func (b *Base64Token) Validate(token string) error {
	payload, err := b.Decode(token)
	if err != nil {
		return fmt.Errorf("token - Validate: %v", err)
	}

	if time.Now().Unix() > payload.Expiry {
		return fmt.Errorf("token has expired")
	}

	return nil
}

func (b *Base64Token) Hash(token string) string {

}

func (b *Base64Token)
