package users

import (
	"github.com/google/uuid"
)

type User struct {
	ID     int64
	UserID uuid.UUID
}
