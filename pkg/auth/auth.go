package auth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	// Packages
	multierror "github.com/hashicorp/go-multierror"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type auth struct {
	Pool
}

var _ Auth = (*auth)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new authorization service with the given pool
func New(pool Pool, opts ...Option) Auth {
	auth := new(auth)
	if pool == nil {
		return nil
	} else {
		auth.Pool = pool
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(auth); err != nil {
			return nil
		}
	}

	return auth
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (auth *auth) String() string {
	str := "<auth"
	if auth.Pool != nil {
		str += fmt.Sprint(" pool=", auth.Pool)
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a new authorization token with the given name, expiry and scopes and
// return the token value
func (auth *auth) CreateByte16(ctx context.Context, name string, duration time.Duration, scopes ...string) (string, error) {
	// Check name
	if !isIdentifier(name) {
		return "", ErrBadParameter.Withf("%q", name)
	}

	// Create the 16-byte token
	token := NewByte16(duration, scopes...)
	if token == nil {
		return "", ErrInternalAppError.With("CreateByte16")
	} else {
		token.Name = name
	}

	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return "", ErrChannelBlocked.With("CreateByte16")
	}
	defer auth.Pool.Put(conn)

	// Check for duplicate name, and insert token if not found
	err := conn.Do(ctx, func(ctx context.Context) error {
		// Fetch token with the given name
		_, _, err := tokenByName(ctx, conn, name)
		if errors.Is(err, ErrNotFound) {
			// continue
		} else if err != nil {
			return err
		} else {
			return ErrDuplicateEntry.With(name)
		}

		// Insert token into database
		if err := conn.Insert(ctx, token); err != nil {
			return err
		}

		// Return success
		return nil
	})

	// Return token value and any errors
	return token.Value, err
}

// List returns all tokens in access order, with the latest token accessed first
func (auth *auth) List(ctx context.Context, fn func(AuthToken)) error {
	if fn == nil {
		return ErrBadParameter
	}

	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return ErrChannelBlocked.With("List")
	}
	defer auth.Pool.Put(conn)

	// Sort by access time, descending
	sort := conn.S()
	sort.Desc("access_at")

	// Iterate over tokens
	cursor, err := conn.Collection(Token{}).FindMany(ctx, sort, nil)
	if err != nil {
		return err
	}
	defer cursor.Close()
	token := new(authtoken)
	for {
		t, err := cursor.Next(ctx)
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		} else {
			fn(token.Set(t.(*Token)))
		}
	}
}

// Expire a token with the given name
func (auth *auth) Expire(ctx context.Context, name string, delete bool) error {
	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return ErrChannelBlocked.With("Expire")
	}
	defer auth.Pool.Put(conn)

	// Check for duplicate name, and insert token if not found
	return conn.Do(ctx, func(ctx context.Context) error {
		// Fetch token with the given name
		_, filter, err := tokenByName(ctx, conn, name)
		if err != nil {
			return err
		}

		// Either update or delete token with the given name
		if delete {
			// Delete token from database
			affected, err := conn.Collection(Token{}).Delete(ctx, filter)
			if err != nil {
				return err
			} else if affected == 0 {
				return ErrNotFound.With("Expire")
			} else {
				return nil
			}
		} else {
			// Update a token, setting the expiry to now
			matched, _, err := conn.Collection(Token{}).Update(ctx, TokenUpdateExpiry{Time: time.Now(), Expire: time.Now()}, filter)
			if err != nil {
				return err
			} else if matched == 0 {
				return ErrNotFound.With("Expire")
			} else {
				return nil
			}
		}
	})
}

// Valid returns true if the token with the given name exists and not
// expired, and has one of the given scopes. Returns ErrNotFound if the token
// with the given name does not exist, ErrExpired if the token has
// expired, or ErrUnauthorized if the token does not have any of the
// given scopes.
func (auth *auth) Valid(ctx context.Context, name string, scope ...string) error {
	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return ErrChannelBlocked.With("Valid")
	}
	defer auth.Pool.Put(conn)

	// Fetch token with the given name
	if token, _, err := tokenByName(ctx, conn, name); err != nil {
		return err
	} else if !token.IsValid() {
		return ErrExpired.With(name)
	} else if !token.IsScope(scope...) {
		return ErrNotAuthorized.With(name)
	}

	// Return success
	return nil
}

// ValidByValue returns the name of a token with the given value, and
// if any of the scopes match. It will return ErrNotFound if no token
// with the given value exists, ErrExpired if the token has expired or
// ErrNotAuthorized if the token does not have any of the given scopes. It
// updates the access_at field of the token if found
func (auth *auth) ValidByValue(ctx context.Context, value string, scope ...string) (string, error) {
	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return "", ErrChannelBlocked.With("NameByValue")
	}
	defer auth.Pool.Put(conn)

	token, filter, err := tokenByValue(ctx, conn, value)
	if err != nil {
		return "", err
	}

	// Check for expired or unauthorized token
	var result error
	if !token.IsValid() {
		result = ErrExpired.With(token.Name)
	} else if !token.IsScope(scope...) {
		result = ErrNotAuthorized.With(token.Name)
	}

	// Update access_at field
	if matched, modifed, err := conn.Collection(Token{}).Update(ctx, TokenUpdateTime{Time: time.Now()}, filter); err != nil {
		result = multierror.Append(result, err)
	} else if matched != 1 || modifed != 1 {
		result = multierror.Append(result, ErrInternalAppError.Withf("Expected 1 match and 1 modified, got %d and %d", matched, modifed))
	}

	// Return token name and any errors
	return token.Name, result
}

// UpdateExpiry updates the duration of a token's life with the given name, or removes
// the expiry if duration is 0 and updates the access_at field of the token
func (auth *auth) UpdateExpiry(ctx context.Context, name string, duration time.Duration) error {
	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return ErrChannelBlocked.With("Valid")
	}
	defer auth.Pool.Put(conn)

	var update TokenUpdateExpiry
	if duration == 0 {
		update = TokenUpdateExpiry{Time: time.Now()}
	} else {
		update = TokenUpdateExpiry{Time: time.Now(), Expire: time.Now().Add(duration)}
	}

	return conn.Do(ctx, func(ctx context.Context) error {
		// Fetch token and filter
		_, filter, err := tokenByName(ctx, conn, name)
		if err != nil {
			return err
		}
		// Update token
		if matched, modified, err := conn.Collection(Token{}).Update(ctx, update, filter); err != nil {
			return err
		} else if matched != 1 || modified != 1 {
			return ErrInternalAppError.Withf("Expected 1 match and 1 modified, got %d and %d", matched, modified)
		}
		// Return success
		return nil
	})
}

// UpdateScopes sets the scopes of a token with the given name and updates
// the access_at field of the token
func (auth *auth) UpdateScope(ctx context.Context, name string, scope ...string) error {
	// Get database connection from pool
	conn := auth.Pool.Get()
	if conn == nil {
		return ErrChannelBlocked.With("Valid")
	}
	defer auth.Pool.Put(conn)

	// Set the update
	update := TokenUpdateScope{
		Time:  time.Now(),
		Scope: scope,
	}

	return conn.Do(ctx, func(ctx context.Context) error {
		// Fetch token and filter
		_, filter, err := tokenByName(ctx, conn, name)
		if err != nil {
			return err
		}
		// Update token
		if matched, modified, err := conn.Collection(Token{}).Update(ctx, update, filter); err != nil {
			return err
		} else if matched != 1 || modified != 1 {
			return ErrInternalAppError.Withf("Expected 1 match and 1 modified, got %d and %d", matched, modified)
		}
		// Return success
		return nil
	})
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func tokenByFilter(ctx context.Context, conn Conn, filter Filter) (*Token, Filter, error) {
	token, err := conn.Collection(Token{}).Find(ctx, nil, filter)
	if err != nil {
		return nil, nil, err
	} else if token, ok := token.(*Token); !ok {
		return nil, nil, ErrInternalAppError
	} else {
		return token, filter, nil
	}
}

func tokenByName(ctx context.Context, conn Conn, name string) (*Token, Filter, error) {
	filter := conn.F()
	if err := filter.Eq("name", name); err != nil {
		return nil, nil, err
	} else {
		return tokenByFilter(ctx, conn, filter)
	}
}

func tokenByValue(ctx context.Context, conn Conn, value string) (*Token, Filter, error) {
	filter := conn.F()
	if err := filter.Eq("value", value); err != nil {
		return nil, nil, err
	} else {
		return tokenByFilter(ctx, conn, filter)
	}
}
