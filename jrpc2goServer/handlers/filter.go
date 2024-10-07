package handlers

import (
	"context"
	"fmt"
	"jrpcServer/shared"
)

// New function to replace the anonymous function
func applyFilter(_ context.Context, req []shared.SecurityLog, filter shared.SecurityLog, messageCounter chan<- struct{}) (string, error) {
	if len(req) == 0 {
		return "Error Code 500: SecurityLog Slice empty", fmt.Errorf("SecurityLog Slice empty")
	}

	if filter.EventType != nil && *filter.EventType != *req[0].EventType {
		fmt.Println(*filter.EventType + " " + *req[0].EventType)
		return "thrown Eventtype missmatch", nil
	}
	if filter.Severity != nil && *filter.Severity > *req[0].Severity {
		return "thrown too low severity", nil
	}
	if filter.TimeStamp != nil && filter.TimeStamp.After(*req[0].TimeStamp) {
		return "thrown timestamp too old", nil
	}
	if filter.SourceIP != nil && *filter.SourceIP != *req[0].SourceIP {
		return "thrown missmatch in source IP", nil
	}
	if filter.UserID != nil && *filter.UserID != *req[0].UserID {
		return "thrown missmatch userID", nil
	}

	messageCounter <- struct{}{}
	return "200, message passed filter", nil
}

// Updated Filter function, now simplified
func Filter(filter shared.SecurityLog, messageCounter chan<- struct{}) func(ctx context.Context, req []shared.SecurityLog) (string, error) {
	return func(ctx context.Context, req []shared.SecurityLog) (string, error) {
		return applyFilter(ctx, req, filter, messageCounter)
	}
}
