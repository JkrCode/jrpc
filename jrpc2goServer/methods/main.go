package methods

import (
	"context"
	"fmt"
	"strings"
)

func Join(ctx context.Context, args []string) string {
	return strings.Join(args, " ")
}

func LogPrint(ctx context.Context, level, message string) error {
	fmt.Println("[%s] %s\n", strings.ToUpper(level), message)
	return nil
}

func Count(ctx context.Context, args []string) int {
	return len(args[0])
}
