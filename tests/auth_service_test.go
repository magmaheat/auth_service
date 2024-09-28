package tests

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
	"net/url"
	"testing"
)

const (
	host = "localhost:8089"
)

func TestAuthServiceAccessAuth(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	id := uuid.New().String()

	response := e.POST("/auth").
		WithQuery("user_id", id).
		Expect().
		Status(201).
		JSON().
		Object()

	response.Value("access_token").String().NotEmpty()
	response.Value("refresh_token").String().NotEmpty()
}

func TestAuthServiceFailAuth(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	_ = e.POST("/auth").
		Expect().
		Status(400).
		JSON().
		Object()
}

func TestAuthServiceAccessUpdate(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	id := uuid.New().String()

	response := e.POST("/auth").
		WithQuery("user_id", id).
		Expect().
		Status(201).
		JSON().
		Object()

	accessToken := response.Value("access_token").String().Raw()
	refreshToken := response.Value("refresh_token").String().Raw()

	response = e.POST("/update").
		WithHeader("Authorization", "Bearer "+accessToken).
		WithJSON(map[string]string{
			"refresh_token": refreshToken,
		}).
		Expect().
		Status(200).
		JSON().
		Object()

	response.Value("access_token").String().NotEmpty()
	response.Value("refresh_token").String().NotEmpty()
}

func TestAuthServiceFailUpdate(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	id := uuid.New().String()

	response := e.POST("/auth").
		WithQuery("user_id", id).
		Expect().
		Status(201).
		JSON().
		Object()

	accessToken := response.Value("access_token").String().Raw()
	refreshToken := response.Value("refresh_token").String().Raw()

	response = e.POST("/update").
		WithHeader("Authorization", "Bearer "+accessToken).
		WithJSON(map[string]string{
			"refresh_token": refreshToken,
		}).
		Expect().
		Status(200).
		JSON().
		Object()

	response.Value("access_token").String().NotEmpty()
	response.Value("refresh_token").String().NotEmpty()

	response = e.POST("/update").
		WithHeader("Authorization", "Bearer "+accessToken).
		WithJSON(map[string]string{
			"refresh_token": refreshToken,
		}).
		Expect().
		Status(401).
		JSON().
		Object()
}
