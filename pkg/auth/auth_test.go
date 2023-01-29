package auth_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	// Packages
	auth "github.com/mutablelogic/go-accessory/pkg/auth"
	pool "github.com/mutablelogic/go-accessory/pkg/pool"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	assert "github.com/stretchr/testify/assert"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

const (
	MONGO_URL = "${MONGO_URL}"
)

func Test_Auth_001(t *testing.T) {
	assert := assert.New(t)
	pool := pool.New(context.TODO(), uri(t))
	auth := auth.New(pool)
	assert.NotNil(auth)
	assert.NoError(pool.Close())
}

func Test_Auth_002(t *testing.T) {
	assert := assert.New(t)
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create a new key
	value, err := auth.CreateByte16(context.Background(), t.Name(), 0)
	assert.NoError(err)
	assert.NotNil(value)

	// Delete the key
	assert.NoError(auth.Expire(context.Background(), t.Name(), true))

	assert.NoError(auth.Close())
}

func Test_Auth_003(t *testing.T) {
	assert := assert.New(t)
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create a new key
	value, err := auth.CreateByte16(context.Background(), t.Name(), 0)
	assert.NoError(err)
	assert.NotNil(value)
	t.Log("Token name=", t.Name(), " value=", value)

	// Try and create the same key again
	_, err = auth.CreateByte16(context.Background(), t.Name(), 0)
	assert.ErrorIs(err, ErrDuplicateEntry)

	// Expire the key
	assert.NoError(auth.Expire(context.Background(), t.Name(), false))

	// Check for expired key
	assert.ErrorIs(auth.Valid(context.Background(), t.Name()), ErrExpired)

	// Delete the key
	assert.NoError(auth.Expire(context.Background(), t.Name(), true))

	// Check for deleted key
	assert.ErrorIs(auth.Valid(context.Background(), t.Name()), ErrNotFound)

	// Close pool
	assert.NoError(auth.Close())
}

func Test_Auth_004(t *testing.T) {
	assert := assert.New(t)
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create a new key
	value, err := auth.CreateByte16(context.Background(), t.Name(), 0)
	assert.NoError(err)
	assert.NotNil(value)
	t.Log("Token name=", t.Name(), " value=", value)

	// Fetch the key
	key, err := auth.ValidByValue(context.Background(), value)
	assert.NoError(err)
	assert.Equal(t.Name(), key)

	// Expire the key
	assert.NoError(auth.Expire(context.Background(), t.Name(), false))

	// Fetch the key, check for expiry
	key, err = auth.ValidByValue(context.Background(), value)
	assert.ErrorIs(err, ErrExpired)
	assert.Equal(t.Name(), key)

	// Delete the token
	assert.NoError(auth.Expire(context.Background(), t.Name(), true))

	// Close pool
	assert.NoError(auth.Close())
}

func Test_Auth_005(t *testing.T) {
	assert := assert.New(t)
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create a new key
	value, err := auth.CreateByte16(context.Background(), t.Name(), 0)
	assert.NoError(err)
	assert.NotNil(value)
	t.Log("Token name=", t.Name(), " value=", value)

	// Update the expiry
	assert.NoError(auth.UpdateExpiry(context.Background(), t.Name(), -1))

	// Fetch the key, check for expiry
	key, err := auth.ValidByValue(context.Background(), value)
	assert.ErrorIs(err, ErrExpired)
	assert.Equal(t.Name(), key)

	// Delete the token
	assert.NoError(auth.Expire(context.Background(), t.Name(), true))

	// Close pool
	assert.NoError(auth.Close())
}

func Test_Auth_006(t *testing.T) {
	assert := assert.New(t)
	r, w := auth.ScopeRead, auth.ScopeWrite
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create a new key with write scope
	value, err := auth.CreateByte16(context.Background(), t.Name(), 0, r)
	assert.NoError(err)
	assert.NotNil(value)
	t.Log("Token name=", t.Name(), " value=", value)

	// Check for write scope (should be not authorized)
	name, err := auth.ValidByValue(context.Background(), value, w)
	assert.ErrorIs(err, ErrNotAuthorized)
	assert.Equal(t.Name(), name)

	// Set scope to r,w
	assert.NoError(auth.UpdateScope(context.Background(), t.Name(), r, w))

	// Check for write scope (should be authorized)
	name, err = auth.ValidByValue(context.Background(), value, w)
	assert.NoError(err)
	assert.Equal(t.Name(), name)

	// Check for read scope (should be authorized)
	name, err = auth.ValidByValue(context.Background(), value, r)
	assert.NoError(err)
	assert.Equal(t.Name(), name)

	// Delete the token
	assert.NoError(auth.Expire(context.Background(), t.Name(), true))

	// Close pool
	assert.NoError(auth.Close())
}

func Test_Auth_007(t *testing.T) {
	assert := assert.New(t)
	r, w := auth.ScopeRead, auth.ScopeWrite
	auth := auth.New(pool.New(context.TODO(), uri(t), pool.OptDatabase("auth"), pool.OptTrace(tracefn(t))))
	assert.NotNil(auth)

	// Create X keys
	var keys []string
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%s-%d", t.Name(), i)
		value, err := auth.CreateByte16(context.Background(), key, 0, r, w)
		assert.NoError(err)
		assert.NotNil(value)
		keys = append(keys, key)
	}

	// Fetch all keys
	assert.NoError(auth.List(context.Background(), func(token AuthToken) {
		json, err := json.MarshalIndent(token, "", "  ")
		assert.NoError(err)
		t.Log(string(json))
	}))

	for _, key := range keys {
		// Delete the token
		assert.NoError(auth.Expire(context.Background(), key, true))
	}

	// Close pool
	assert.NoError(auth.Close())
}

///////////////////////////////////////////////////////////////////////////////
// Return URL or skip test

func uri(t *testing.T) *url.URL {
	if uri := os.ExpandEnv(MONGO_URL); uri == "" {
		t.Skip("no MONGO_URL environment variable, skipping test")
	} else if uri, err := url.Parse(uri); err != nil {
		t.Skip("invalid MONGO_URL environment variable, skipping test")
	} else {
		return uri
	}
	return nil
}

func tracefn(t *testing.T) trace.Func {
	return func(ctx context.Context, dur time.Duration, err error) {
		t.Log("TRACE:", trace.DumpContextStr(ctx))
		if dur != 0 {
			t.Log("  DURATION:", dur)
		}
		if err != nil {
			t.Log("  ERROR:", err)
		}
	}
}
