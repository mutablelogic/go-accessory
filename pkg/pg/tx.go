package pg

import (
	"context"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (conn *conn) BeginTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	if _, err := conn.Exec(ctx, "BEGIN"); err != nil {
		return err
	} else {
		return nil
	}
}

func (conn *conn) RollbackTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	if _, err := conn.Exec(ctx, "ROLLBACK"); err != nil {
		return err
	} else {
		return nil
	}
}

func (conn *conn) CommitTx(ctx context.Context) error {
	conn.Mutex.Lock()
	defer conn.Mutex.Unlock()

	if _, err := conn.Exec(ctx, "COMMIT"); err != nil {
		return err
	} else {
		return nil
	}
}
