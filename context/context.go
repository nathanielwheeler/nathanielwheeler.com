package context

import (
  "context"

  "nathanielwheeler.com/models"
)

// NOTE these should never be exported.  This prevents outside code from changing these values.
type privateKey string

const userKey privateKey = "user"

// WithUser accepts an existing context and a user, then returns a new context with that user set as a value.
func WithUser(ctx context.Context, user *models.User) context.Context {
  return context.WithValue(ctx, userKey, user)
}

// User will look up a user from a given context.
func User(ctx context.Context) *models.User {
  if temp := ctx.Value(userKey); temp != nil {
    if user, ok := temp.(*models.User); ok {
      return user
    }
  }
  return nil
}
