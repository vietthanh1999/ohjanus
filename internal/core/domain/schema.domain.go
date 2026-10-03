package domain

// Column describes a single table column.
type Column struct {
	Name     string
	Type     string
	Nullable bool
}

// Table describes a database table.
type Table struct {
	Name       string
	Columns    []Column
	PrimaryKey []string
}

// Schema groups tables under a schema name.
type Schema struct {
	Name   string
	Tables []Table
}

// Connection is the metadata the agent is allowed to see.
// It never contains credentials or network addresses.
type Connection struct {
	Name     string
	Driver   string
	ReadOnly bool
}

// ConnectionMeta carries the allow/deny configuration for one connection alias.
type ConnectionMeta struct {
	Connection     Connection
	AllowedSchemas []string
	DeniedTables   []string
	AllowedTables  []string
	RowLimit       int
}
