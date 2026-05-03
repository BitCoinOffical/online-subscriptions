package models

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id            uuid.UUID
	Email         string
	Password_hash string
	Created_at    time.Time
	Updated_at    time.Time
}
