package main

import (
	"context"
	"net/url"

	// Packages
	pool "github.com/mutablelogic/go-accessory/pkg/pool"
	iface "github.com/mutablelogic/go-server"
	task "github.com/mutablelogic/go-server/pkg/task"

	// Namespace imports
	. "github.com/djthorpe/go-errors"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Plugin struct {
	task.Plugin
	Url_ string `json:"url,omitempty"` // Database URL
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	defaultName = "pool"
)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new pool task from plugin configuration
func (p Plugin) New(ctx context.Context, provider iface.Provider) (iface.Task, error) {
	if err := p.HasNameLabel(); err != nil {
		return nil, err
	}

	// Get the URL
	if url, err := p.Url(); err != nil {
		return nil, err
	} else if pool := pool.New(ctx, url); pool == nil {
		return nil, ErrBadParameter.With(url)
	} else {
		return NewWithPool(pool)
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func WithLabel(label string) Plugin {
	return Plugin{
		Plugin: task.WithLabel(defaultName, label),
	}
}

func (p Plugin) Name() string {
	if name := p.Plugin.Name(); name != "" {
		return name
	} else {
		return defaultName
	}
}

func (p Plugin) WithUrl(url string) Plugin {
	p.Url_ = url
	return p
}

func (p Plugin) Url() (*url.URL, error) {
	if p.Url_ == "" {
		return nil, ErrBadParameter.With("Url")
	} else if url, err := url.Parse(p.Url_); err != nil {
		return nil, err
	} else {
		return url, nil
	}
}
