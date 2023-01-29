package auth

import (
	"time"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// authtoken is a wrapper around a *Token which implements the
// AuthToken interface
type authtoken struct {
	t *Token
}

var _ AuthToken = (*authtoken)(nil)

/////////////////////////////////////////////////////////////////////
// LIEFCYCLE

func NewAuthToken(t *Token) *authtoken {
	return &authtoken{t}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (authtoken *authtoken) String() string {
	return authtoken.t.String()
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the name of the token
func (authtoken *authtoken) Name() string {
	if authtoken.t == nil {
		return ""
	} else {
		return authtoken.t.Name
	}
}

// Return the token type (currently only 'byte16' is supported)
func (authtoken *authtoken) Type() string {
	if authtoken.t == nil {
		return ""
	} else {
		return string(authtoken.t.Type)
	}
}

// Return the scopes of the token
func (authtoken *authtoken) Scope() []string {
	if authtoken.t == nil {
		return nil
	} else {
		return authtoken.t.Scope
	}
}

// Return the last access time for the token
func (authtoken *authtoken) AccessAt() time.Time {
	if authtoken.t == nil {
		return time.Time{}
	} else {
		return authtoken.t.Time
	}
}

// Return the expiry time for the token, or time.Zero if there
// is no expiry
func (authtoken *authtoken) ExpireAt() time.Time {
	if authtoken.t == nil {
		return time.Time{}
	} else {
		return authtoken.t.Expire
	}
}

// Return ErrExpired if the token has expired, or ErrNotAuthorized if
// any given scopes are not in the token's scopes. Otherwise return nil
func (authtoken *authtoken) Valid(scope ...string) error {
	if authtoken.t == nil {
		return ErrNotFound
	} else if !authtoken.t.IsValid() {
		return ErrExpired
	} else if !authtoken.t.IsScope(scope...) {
		return ErrNotAuthorized
	} else {
		return nil
	}
}

/////////////////////////////////////////////////////////////////////
// JSON MARSHAL

func (authtoken *authtoken) MarshalJSON() ([]byte, error) {
	if authtoken == nil || authtoken.t == nil {
		return []byte("null"), nil
	} else {
		return authtoken.t.MarshalJSON()
	}
}
