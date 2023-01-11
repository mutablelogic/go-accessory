package trace

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Op uint

///////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	OpNone Op = iota
	OpConnect
	OpDisconnect
	OpPing
	OpTransaction
	OpInsert
	OpInsertMany
	OpDelete
	OpDeleteMany
	OpFind
	OpFindMany
	OpUpdate
	OpUpdateMany
	OpUpsert
	OpUpsertMany
)

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (o Op) String() string {
	switch o {
	case OpNone:
		return "OpNone"
	case OpConnect:
		return "OpConnect"
	case OpDisconnect:
		return "OpDisconnect"
	case OpPing:
		return "OpPing"
	case OpTransaction:
		return "OpTransaction"
	case OpInsert:
		return "OpInsert"
	case OpInsertMany:
		return "OpInsertMany"
	case OpDelete:
		return "OpDelete"
	case OpDeleteMany:
		return "OpDeleteMany"
	case OpFind:
		return "OpFind"
	case OpFindMany:
		return "OpFindMany"
	case OpUpdate:
		return "OpUpdate"
	case OpUpdateMany:
		return "OpUpdateMany"
	default:
		return "[?? Invalid Operation value]"
	}
}
