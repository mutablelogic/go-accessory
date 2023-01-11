package trace

import (
	"context"
	"fmt"
	"net/url"
	"sync/atomic"
	"time"
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
}

///////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	ctxCol   ctxKey = iota // Collection operation (update, find, ...)
	ctxUrl                 // URL operation (connect, disconnect and ping)
	ctxTx                  // Transaction number
	ctxDelta               // Delta time
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func WithUrl(parent context.Context, op Op, url *url.URL) context.Context {
	return context.WithValue(parent, ctxUrl, urlOp{op, redactedUrl(url)})
}

func WithTx(parent context.Context) context.Context {
	return context.WithValue(parent, ctxTx, nextTx())
}

func WithDelta(parent context.Context, delta time.Duration) context.Context {
	return context.WithValue(parent, ctxDelta, delta)
}

func WithCol(parent context.Context, op Op, database, collection string, delta time.Duration) context.Context {
	return context.WithValue(parent, ctxCol, colOp{op, database, collection})
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
	}
	if col, ok := ctx.Value(ctxCol).(colOp); ok {
		str += fmt.Sprintf(" op=%v", col.Op)
		if col.database != "" {
			str += fmt.Sprintf(" database=%q", col.database)
		}
		if col.collection != "" {
			str += fmt.Sprintf(" collection=%q", col.collection)
		}
	}
	if delta, ok := ctx.Value(ctxDelta).(time.Duration); ok && delta > 0 {
		str += fmt.Sprint(" delta=", delta.Truncate(time.Millisecond))
	}
	return str + ">"
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
