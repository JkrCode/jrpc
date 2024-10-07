package handlers

import (
	"context"
	"strconv"
)

// Count function accepts a channel for counting messages
func Count(messageCounter chan<- struct{}) func(ctx context.Context, req []string) (string, error) {
	return func(ctx context.Context, req []string) (string, error) {
		count := strconv.Itoa(len(req[0]))
		messageCounter <- struct{}{}
		return count, nil
	}
}
