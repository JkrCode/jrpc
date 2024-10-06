package handlers

import (
	"context"
)

// Count function accepts a channel for counting messages
func Count(messageCounter chan<- struct{}) func(ctx context.Context, req []string) (int, error) {
	return func(ctx context.Context, req []string) (int, error) {
		count := len(req[0])

		messageCounter <- struct{}{}

		return count, nil
	}
}