package accessory

///////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Query represents any kind of SQL Query
type Query interface {
	// Return the SQL query that can be executed
	Query() string
}

// Name represents an SQL name (table name, column name)
type Name interface {
	Query

	// Use a specific schema name
	WithSchema(string) Name

	// Use a specific alias name
	WithAlias(string) Name

	// Add a DESC clause or ORDER
	WithDesc() Name

	// Return a column with a specific type
	Column(string) Column
}

// Column is a column in a table with a type
type Column interface {
	Name

	// Set NOT NULL
	NotNull() Column

	// Identify as a primary key
	WithPrimary() Column

	// Add AUTO_INCREMENT keyword
	WithAutoIncrement() Column
}
