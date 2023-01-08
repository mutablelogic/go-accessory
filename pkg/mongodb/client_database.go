package mongodb

import (
	"context"
)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (client *client) Collections(ctx context.Context) ([]string, error) {
	return client.Database(defaultDatabase).Collections(ctx)
}

func (client *client) Insert(ctx context.Context, document any) (string, error) {
	return client.Database(defaultDatabase).Insert(ctx, document)
}
