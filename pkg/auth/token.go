package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	// Package imports
	slices "golang.org/x/exp/slices"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type Token struct {
	Key    string    `json:"key,omitempty" bson:"_id,omitempty"`               // Key
	Name   string    `json:"name" bson:"name"`                                 // Name of token
	Type   TokenType `json:"type" bson:"type"`                                 // Type of token
	Value  string    `json:"-" bson:"value"`                                   // Token value
	Expire time.Time `json:"expires_at,omitempty" bson:"expires_at,omitempty"` // Time of expiration for the token
	Time   time.Time `json:"access_at" bson:"access_at"`                       // Time of last access
	Scope  []string  `json:"scope,omitempty" bson:"scope,omitempty"`           // Authorization scopes
}

type TokenType string

type TokenUpdateExpiry struct {
	Time   time.Time `bson:"access_at"`            // Time of last access
	Expire time.Time `bson:"expires_at,omitempty"` // Time of expiration for the token
}

type TokenUpdateTime struct {
	Time time.Time `bson:"access_at"` // Time of last access
}

type TokenUpdateScope struct {
	Time  time.Time `bson:"access_at"`       // Time of last access
	Scope []string  `bson:"scope,omitempty"` // Authentication scopes
}

/////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TokenByte16 TokenType = "byte16" // 16 byte value
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

// NewByte16 creates a new 16-byte token without a name
func NewByte16(duration time.Duration, scope ...string) *Token {
	var expire time.Time
	if duration != 0 {
		expire = time.Now().Add(duration)
	}
	return &Token{
		Value:  generateToken(TokenByte16),
		Type:   TokenByte16,
		Time:   time.Now(),
		Scope:  scope,
		Expire: expire,
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t *Token) String() string {
	str := "<token"
	str += fmt.Sprintf(" value=%q", t.Value)
	if t.Name != "" {
		str += fmt.Sprintf(" name=%q", t.Name)
	}
	if t.Key != "" {
		str += fmt.Sprintf(" key=%q", t.Key)
	}
	if t.Type != "" {
		str += fmt.Sprintf(" type=%q", t.Type)
	}
	if !t.Time.IsZero() {
		str += fmt.Sprintf(" access_time=%q", t.Time.Format(time.RFC3339))
	}
	if !t.Expire.IsZero() {
		str += fmt.Sprintf(" expire_time=%q", t.Expire.Format(time.RFC3339))
	}
	if len(t.Scope) > 0 {
		str += fmt.Sprintf(" scopes=%q", t.Scope)
	}
	if t.IsValid() {
		str += " valid"
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return true if the token is valid (not expired)
func (t *Token) IsValid() bool {
	if t.Expire.IsZero() || t.Expire.After(time.Now()) {
		return true
	}
	return false
}

// Return true if the token has the specified scope, and is valid
func (t *Token) IsScope(scopes ...string) bool {
	if !t.IsValid() {
		return false
	}
	if len(scopes) == 0 {
		return true
	}
	for _, scope := range scopes {
		if slices.Contains(t.Scope, scope) {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////
// JSON MARSHAL

func (t *Token) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if t == nil {
		return []byte("null"), nil
	}
	buf.WriteRune('{')

	// Write the fields
	if t.Key != "" {
		buf.WriteString(`"key":`)
		buf.WriteString(strconv.Quote(t.Key))
		buf.WriteRune(',')
	}
	if t.Name != "" {
		buf.WriteString(`"name":`)
		buf.WriteString(strconv.Quote(t.Name))
		buf.WriteRune(',')
	}
	if t.Type != "" {
		buf.WriteString(`"type":`)
		buf.WriteString(strconv.Quote(string(t.Type)))
		buf.WriteRune(',')
	}
	/*
		if t.Value != "" {
			buf.WriteString(`"token":`)
			buf.WriteString(strconv.Quote(t.Value))
			buf.WriteRune(',')
		}
	*/
	if !t.Expire.IsZero() {
		buf.WriteString(`"expires_at":`)
		buf.WriteString(strconv.Quote(t.Expire.Format(time.RFC3339)))
		buf.WriteRune(',')
	}
	if !t.Time.IsZero() {
		buf.WriteString(`"access_at":`)
		buf.WriteString(strconv.Quote(t.Time.Format(time.RFC3339)))
		buf.WriteRune(',')
	}
	if len(t.Scope) > 0 {
		buf.WriteString(`"scopes":`)
		if data, err := json.Marshal(t.Scope); err != nil {
			return nil, err
		} else {
			buf.Write(data)
		}
		buf.WriteRune(',')
	}

	// Include the valid flag
	buf.WriteString(`"valid":`)
	buf.WriteString(strconv.FormatBool(t.IsValid()))

	// Return success
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func generateToken(t TokenType) string {
	switch t {
	case TokenByte16:
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return ""
		}
		return hex.EncodeToString(b)
	default:
		return ""
	}
}
