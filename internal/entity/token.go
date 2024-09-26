package entity

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	UserId      uuid.UUID `db:"user_id"`
	Token       string    `db:"token"`
	StatusToken string    `db:"status_token"`
	CreatedAt   time.Time `db:"created_at"`
}
