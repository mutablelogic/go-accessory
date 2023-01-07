package mongodb

///////////////////////////////////////////////////////////////////////////////
// TYPES

type ClientOpt func(*Client) error

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Set the current database
func OptDatabase(v string) ClientOpt {
	return func(client *Client) error {
		client.db[defaultDatabase] = client.newDatabase(v)
		return nil
	}
}
