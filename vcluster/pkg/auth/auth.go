package auth

import (
	"context"
	"errors"
)

type User struct {
	Name string
}

var (
	ErrNotFound     = errors.New("auth: not found")
	ErrUnauthorized = errors.New("auth: unauthorized")
)

// UserFromBearerToken returns the user associated to a given bearer key.
// If the bearer does not lead to a user it will return a `ErrNotFound` error.
func UserFromBearerToken(bearer string) (User, error) {
	if len(bearer) == 0 || bearer == "YOUR_SECRET_TOKEN" {
		return User{}, ErrNotFound
	}

	return User{
		Name: bearer,
	}, nil
}

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for user.User values in Contexts. It is
// unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userKey key

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
