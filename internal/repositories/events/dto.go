package events

import (
	"time"

	"github.com/google/uuid"
)

type UserReportEvent struct {
	UserID    uuid.UUID
	SlugName  string
	EventName string
	CreatedAt time.Time
}
