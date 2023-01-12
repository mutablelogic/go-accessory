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
	OpCommit
	OpRollback
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
		return "None"
	case OpConnect:
		return "Connect"
	case OpDisconnect:
		return "Disconnect"
	case OpPing:
		return "Ping"
	case OpInsert:
		return "Insert"
	case OpInsertMany:
		return "InsertMany"
	case OpDelete:
		return "Delete"
	case OpDeleteMany:
		return "DeleteMany"
	case OpFind:
		return "Find"
	case OpFindMany:
		return "FindMany"
	case OpUpdate:
		return "Update"
	case OpUpdateMany:
		return "UpdateMany"
	case OpCommit:
		return "Commit"
	case OpRollback:
		return "Rollback"
	default:
		return "[?? Invalid Operation value]"
	}
}
