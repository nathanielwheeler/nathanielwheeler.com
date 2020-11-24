package util

import (
	"context"

	"nathanielwheeler.com/server/services"
)

// NOTE these should never be exported.  This prevents outside code from changing these values.
type privateKey string

const userKey privateKey = "user"

// WithUser accepts an existing context and a user, then returns a new context with that user set as a value.
func WithUser(ctx context.Context, user *services.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// User will look up a user from a given context.
func User(ctx context.Context) *services.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*services.User); ok {
			return user
		}
	}
	return nil
}
