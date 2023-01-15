package main

import (

	// Packages
	iface "github.com/mutablelogic/go-server"
	task "github.com/mutablelogic/go-server/pkg/task"

	// Namespace imports
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type t struct {
	task.Task
	Pool
}

var _ iface.Task = (*t)(nil)

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

// Create a new pool task from plugin configuration
func NewWithPool(pool Pool) (iface.Task, error) {
	t := new(t)
	t.Pool = pool
	return t, nil
}
