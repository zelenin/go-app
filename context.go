package app

import "context"

const (
	keyAppPath = "app-path"
)

func GetAppPath(ctx context.Context) string {
	return ctx.Value(keyAppPath).(string)
}
