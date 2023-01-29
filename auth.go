package accessory

import (
	"context"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Auth represents authorzation storage
type Auth interface {
	// Auth implements the connection pool
	Pool

	// Create a new byte16 token with the given name, token duration and scopes
	CreateByte16(context.Context, string, time.Duration, ...string) (string, error)

	// List returns all tokens in access order, with the latest token accessed first
	List(context.Context, func(AuthToken)) error

	// Expire a token with the given name. When argument is set to
	// true, the token is deleted from the database, otherwise it is disabled
	Expire(context.Context, string, bool) error

	// Return no error if the token with the given name exists and not expired,
	// and has one of the given scopes. Returns ErrNotFound if the token
	// with the given name does not exist, ErrExpired if the token has
	// expired, or ErrUnauthorized if the token does not have any of the
	// given scopes.
	Valid(context.Context, string, ...string) error

	// ValidByValue returns the name of a token with the given value, and
	// if any of the scopes match. It will return ErrNotFound if no token
	// with the given value exists, ErrExpired if the token has expired or
	// ErrNotAuthorized if the token does not have any of the given scopes. It
	// updates the access_at field of the token if found
	ValidByValue(context.Context, string, ...string) (string, error)

	// UpdateExpiry updates the duration of a token's life with the given name,
	// or removes the expiry if duration is 0 and updates the access_at field
	// of the token
	UpdateExpiry(context.Context, string, time.Duration) error

	// UpdateScope sets the scopes of a token with the given name and updates
	// the access_at field of the token
	UpdateScope(context.Context, string, ...string) error
}

// AuthToken represents an authorization token
type AuthToken interface {
	// Return the name of the token
	Name() string

	// Return the token type (currently only 'byte16' is supported)
	Type() string

	// Return the scopes of the token
	Scope() []string

	// Return the last access time for the token
	AccessAt() time.Time

	// Return the expiry time for the token, or time.Zero if there
	// is no expiry
	ExpireAt() time.Time

	// Return ErrExpired if the token has expired, or ErrNotAuthorized if
	// any given scopes are not in the token's scopes. Otherwise return nil
	Valid(...string) error
}
