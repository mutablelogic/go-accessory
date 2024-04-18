package ast

///////////////////////////////////////////////////////////////////////////////
// TYPES

// Classifies the kind of node
type Kind uint

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	Any Kind = iota
	Number
	Bool
	String
	Ident
	BinaryOp
	UnaryOp
	Func
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// IsKindValue returns true if the kind is a scalar value
func IsKindValue(kind Kind) bool {
	switch kind {
	case String, Number, Bool, Ident:
		return true
	default:
		return false
	}
}

// IsKindExpr returns true if the kind is an value (a scalar value) or an operation (plus, minus etc)
func IsKindExpr(kind Kind) bool {
	switch {
	case IsKindValue(kind):
		return true
	case kind == BinaryOp || kind == UnaryOp:
		return true
	default:
		return false
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (k Kind) String() string {
	switch k {
	case Number:
		return "Number"
	case Bool:
		return "Bool"
	case String:
		return "String"
	case Ident:
		return "Ident"
	case BinaryOp:
		return "BinaryOp"
	case UnaryOp:
		return "UnaryOp"
	case Func:
		return "Func"
	default:
		return "[?? Invalid Kind value]"
	}

}
