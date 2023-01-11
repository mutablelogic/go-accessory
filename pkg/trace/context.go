package trace

import (
	"context"
	"fmt"
	"runtime/internal/atomic"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ctxKey uint

type urlOp struct {
	Op
	url   string
	delta time.Duration
}

type colOp struct {
	Op
	database, collection string
	delta                time.Duration
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ctxCol ctxKey = iota // Collection operation (update, find, ...)
	ctxUrl               // URL operation (connect, disconnect and ping)
	ctxTx                // Transaction number
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func WithUrl(parent context.Context, op Op, url string, delta time.Duration) context.Context {
	return context.WithValue(parent, ctxUrl, urlOp{op, url, delta})
}

func WithTx(parent context.Context, op Op, tx uint) context.Context {
	return context.WithValue(parent, ctxTx, nextTx())
}

func WithCol(parent context.Context, op Op, database, collection string, delta time.Duration) context.Context {
	return context.WithValue(parent, ctxCol, colOp{op, database, collection, delta})
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func DumpContextStr(ctx context.Context) string {
	str := "<trace"
	if tx, ok := ctx.Value(ctxTx).(uint64); ok {
		str += fmt.Sprint(" tx=", tx)
	}
	if url, ok := ctx.Value(ctxUrl).(urlOp); ok {
		str += fmt.Sprintf(" op=%v url=%q", url.Op, url.url)
		str += fmt.Sprint(" delta=", url.delta.Truncate(time.Millisecond))
	}
	if col, ok := ctx.Value(ctxCol).(colOp); ok {
		str += fmt.Sprintf(" op=%v", col.Op)
		if col.database != "" {
			str += fmt.Sprintf(" database=%q", col.database)
		}
		if col.collection != "" {
			str += fmt.Sprintf(" collection=%q", col.collection)
		}
		str += fmt.Sprint(" delta=", col.delta.Truncate(time.Millisecond))
	}
	return str + ">"
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

var tx atomic.Uint64

// Return a new transaction number
func nextTx() uint64 {
	return tx.Add(1)
}
