package pg_test

import (
	"context"
	"testing"

	// Packages
	pg "github.com/mutablelogic/go-accessory/pkg/pg"
	trace "github.com/mutablelogic/go-accessory/pkg/trace"
	assert "github.com/stretchr/testify/assert"
)

func Test_Tx_001(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("Skipping, connection error")
	}
	defer func() {
		assert.NoError(c.Close())
	}()

	ctx := trace.WithTx(context.TODO())
	assert.NoError(c.(pg.Tx).BeginTx(ctx))
	assert.NoError(c.(pg.Tx).RollbackTx(ctx))
}

func Test_Tx_002(t *testing.T) {
	assert := assert.New(t)
	c, err := pg.Open(context.TODO(), uri(t), tracefn(t), pg.OptApplicationName(t.Name()))
	assert.NoError(err)
	if c == nil {
		t.Skip("Skipping, connection error")
	}
	defer func() {
		assert.NoError(c.Close())
	}()

	ctx := trace.WithTx(context.TODO())
	assert.NoError(c.(pg.Tx).BeginTx(ctx))
	assert.NoError(c.(pg.Tx).CommitTx(ctx))
}
