package pg

import (
	"strings"
	"time"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Opt func(*opts) error

type opts struct {
	// The default timeout
	timeout time.Duration

	// Function to trace calls
	tracefn trace.Func

	// Application name
	applicationName string

	// User and password if not included in the URL
	user, password string

	// Schema
	schema string
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the default timeout
func OptTimeout(v time.Duration) Opt {
	return func(o *opts) error {
		if v <= 0 {
			o.timeout = defaultTimeout
		} else {
			o.timeout = v
		}
		return nil
	}
}

// Set the trace function
func OptTrace(fn trace.Func) Opt {
	return func(o *opts) error {
		o.tracefn = fn
		return nil
	}
}

// Set the application name
func OptApplicationName(v string) Opt {
	return func(o *opts) error {
		o.applicationName = v
		return nil
	}
}

// Set the credentials
func OptCredentials(user, password string) Opt {
	return func(o *opts) error {
		o.user = strings.TrimSpace(user)
		o.password = strings.TrimSpace(password)
		return nil
	}
}

// Set the schema
func OptSchema(v string) Opt {
	return func(o *opts) error {
		o.schema = v
		return nil
	}
}
