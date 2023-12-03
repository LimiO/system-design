package web

import "context"

type UserIDKey struct{}

func GetLogin(ctx context.Context) string {
	if val, ok := ctx.Value(UserIDKey{}).(string); ok {
		return val
	}
	return ""
}
