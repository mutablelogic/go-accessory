package ast

import "reflect"

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Node is a node in the graph
type Node interface {
	// Kind returns the node kind
	Kind() Kind

	// Children returns the child nodes
	Children() []Node

	// Type returns the type of the node, or nil
	Type() reflect.Type

	// Eval returns the value of the node
	Eval(*Context) (any, error)
}
