package mongodb

import (
	// Package imports

	// Namespace imports
	//. "github.com/djthorpe/go-errors"
	. "github.com/mutablelogic/go-accessory"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type collection struct{}

// Ensure *collection implements the Collection interface
var _ Collection = (*collection)(nil)
