package ast

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Context struct{}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewContext() *Context {
	return &Context{}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS