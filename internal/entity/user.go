package entity

import (
	"github.com/google/uuid"
)

type User struct {
	UserId      uuid.UUID `db:"user_id"`
	Token       string    `db:"token"`
	StatusToken string    `db:"status_token"`
}
