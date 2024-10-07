package shared

import (
	"time"
)

type Filter struct {
	TimeStamp time.Time
	EventType string
	SourceIP  string
	UserID    string
	Severity  int
}
