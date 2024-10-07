package shared

import (
	"time"
)

type SecurityLog struct {
	TimeStamp *time.Time
	EventType *string
	SourceIP  *string
	UserID    *string
	Severity  *int
}
