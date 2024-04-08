/*
ForeignKey() is a factory method that returns a foreign key constraint:

	ForeignKey().References("other_table") => REFERENCES other_table
	ForeignKey("b","c").References("other_table","c1","c2") => FOREIGN KEY (b, c) REFERENCES other_table (c1, c2)
	ForeignKey("a").References("other_table") => FOREIGN KEY a REFERENCES other_table
	ForeignKey().References("other_table","c1") => REFERENCES other_table (c1)
	ForeignKey().References("other_table").OnDeleteRestrict() => REFERENCES other_table (c1) ON DELETE RESTRICT
	ForeignKey().References("other_table").OnDeleteCascade() => REFERENCES other_table (c1) ON DELETE CASCADE
	ForeignKey().References("other_table").OnDeleteNoAction() => REFERENCES other_table (c1) ON DELETE NO ACTION

PrimaryKey() is a factory method that returns a primary key constraint:

	PrimaryKey() => PRIMARY KEY
	PrimaryKey("a","b") => PRIMARY KEY (a,b)
*/
package querybuilder

///////////////////////////////////////////////////////////////////////////////
// TYPES

type foreignKey struct {
	flags
	name    name
	foreign []any
	columns []any
}

///////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func ForeignKey(v ...any) foreignKey {
	// TODO: Convert string to N()
	return foreignKey{columns: v}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (q foreignKey) String() string {
	if len(q.columns) == 0 && q.name.name == "" && len(q.foreign) == 0 {
		return ""
	}
	return "FOREIGN KEY " + join(q.columns...) + " REFERENCES " + join(q.foreign...)
}
