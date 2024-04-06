package pg

import (
	"context"
	"errors"
	"fmt"

	// Packages
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
)

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Tx interface {
	BeginTx(context.Context) error
	RollbackTx(context.Context) error
	CommitTx(context.Context) error
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (conn *conn) BeginTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	if _, err := conn.Exec(ctx, "BEGIN"); err != nil {
		return err
	} else if tx := trace.Tx(ctx); tx == 0 {
		return nil
	} else if _, err := conn.Exec(ctx, "SAVEPOINT p"+fmt.Sprint(tx)); err != nil {
		conn.Exec(ctx, "ROLLBACK")
		return err
	}

	// Return success
	return nil
}

func (conn *conn) RollbackTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	if tx := trace.Tx(ctx); tx > 0 {
		if _, err := conn.Exec(ctx, "ROLLBACK TO SAVEPOINT p"+fmt.Sprint(tx)); err != nil {
			return err
		}
	} else {
		if _, err := conn.Exec(ctx, "ROLLBACK"); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (conn *conn) CommitTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	// Release the savepoint
	var result error
	if tx := trace.Tx(ctx); tx > 0 {
		if _, err := conn.Exec(ctx, "RELEASE SAVEPOINT p"+fmt.Sprint(tx)); err != nil {
			result = errors.Join(result, err)
		}
	}

	// Commit the transaction
	if _, err := conn.Exec(ctx, "COMMIT"); err != nil {
		result = errors.Join(result, err)
	}

	// Return any errors
	return result
}
