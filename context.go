package app

import "context"

type ctxKey string

const (
	keyAppPath = ctxKey("app-path")
)

func GetAppPath(ctx context.Context) string {
	return ctx.Value(keyAppPath).(string)
}
