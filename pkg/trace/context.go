package trace

import (
	"context"
	"fmt"
	"net/url"
	"sync/atomic"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ctxKey uint

type urlOp struct {
	Op
	url string
}

type colOp struct {
	Op
	database, collection string
	matched, modified    int64
}

type execOp struct {
	Op
	statement string
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ctxCol  ctxKey = iota // Collection operation (update, find, ...)
	ctxUrl                // URL operation (connect, disconnect and ping)
	ctxTx                 // Transaction number
	ctxOp                 // Operation (insert, update, delete, find, ...)
	ctxExec               // Exec operation (exec, query, ...
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func WithOp(parent context.Context, op Op) context.Context {
	return context.WithValue(parent, ctxOp, op)
}

func WithUrl(parent context.Context, op Op, url *url.URL) context.Context {
	return context.WithValue(parent, ctxUrl, urlOp{op, redactedUrl(url)})
}

func WithTx(parent context.Context) context.Context {
	return context.WithValue(parent, ctxTx, nextTx())
}

func WithExec(parent context.Context, exec any) context.Context {
	return context.WithValue(parent, ctxExec, execOp{OpExec, fmt.Sprint(exec)})
}

// Return a new context which contains matched and modified placeholders
func WithCollection(parent context.Context, op Op, database, collection string) (context.Context, *int64, *int64) {
	result := &colOp{op, database, collection, -1, -1}
	return context.WithValue(parent, ctxCol, result), &result.matched, &result.modified
}

func DumpContextStr(ctx context.Context) string {
	str := "<trace"
	if tx, ok := ctx.Value(ctxTx).(uint64); ok {
		str += fmt.Sprint(" tx=", tx)
	}
	if op, ok := ctx.Value(ctxOp).(Op); ok {
		str += fmt.Sprintf(" op=%v", op)
	}
	if url, ok := ctx.Value(ctxUrl).(urlOp); ok {
		str += fmt.Sprintf(" op=%v url=%q", url.Op, url.url)
	}
	if exec, ok := ctx.Value(ctxExec).(execOp); ok {
		str += fmt.Sprintf(" op=%v statement=%q", exec.Op, exec.statement)
	}
	if col, ok := ctx.Value(ctxCol).(*colOp); ok {
		str += fmt.Sprintf(" op=%v", col.Op)
		if col.database != "" {
			str += fmt.Sprintf(" database=%q", col.database)
		}
		if col.collection != "" {
			str += fmt.Sprintf(" collection=%q", col.collection)
		}
		if col.matched >= 0 {
			str += fmt.Sprintf(" matched=%d", col.matched)
		}
		if col.modified >= 0 {
			str += fmt.Sprintf(" modified=%d", col.modified)
		}
	}
	return str + ">"
}

// Return the current tx
func Tx(ctx context.Context) uint64 {
	if tx, ok := ctx.Value(ctxTx).(uint64); ok {
		return tx
	}
	return 0
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

var tx uint64

// Return a new transaction number
func nextTx() uint64 {
	return atomic.AddUint64(&tx, 1)
}

// Remove any usernames and passwords before printing out
func redactedUrl(url *url.URL) string {
	url_ := *url // make a copy
	url_.User = nil
	return url_.String()
}
