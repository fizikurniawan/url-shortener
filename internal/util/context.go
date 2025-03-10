package util

import (
	"context"
	"url-shortener/internal/middleware"
	"url-shortener/internal/model"
)

// GetUserFromContext is a helper function to get user from context
func GetUserFromContext(ctx context.Context) *model.User {
	if ctx == nil {
		return nil
	}

	if user, ok := ctx.Value(middleware.UserContextKey).(*model.User); ok {
		return user
	}

	return nil
}
