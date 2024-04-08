package querybuilder_test

import (
	"fmt"
	"testing"

	// Packages
	"github.com/stretchr/testify/assert"

	// Import namespaces
	. "github.com/mutablelogic/go-accessory/pkg/querybuilder"

	// Import PG
	_ "github.com/mutablelogic/go-accessory/pkg/pg"
)

func Test_PG_Create_Table_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").CreateTable(), `CREATE TABLE a ()`},
		{N("a").CreateTable().IfNotExists(), `CREATE TABLE IF NOT EXISTS a ()`},
		{N("a").WithSchema("public").CreateTable().Temporary(), `CREATE TEMPORARY TABLE public.a ()`},
		{N("a").WithSchema("public").CreateTable("a", "b", "c").Temporary(), `CREATE TEMPORARY TABLE public.a (a TEXT,b TEXT,c TEXT)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}

func Test_PG_Create_Table_001(t *testing.T) {
	// Test for primary keys
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").CreateTable(N("c1").T("text").PrimaryKey()), `CREATE TABLE a (c1 TEXT PRIMARY KEY)`},
		{N("a").CreateTable(N("c1").T("text").PrimaryKey(), N("c2").T("text").PrimaryKey()), `CREATE TABLE a (c1 TEXT,c2 TEXT,PRIMARY KEY (c1,c2))`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}

func Test_PG_Create_Table_002(t *testing.T) {
	// Test for foreign keys
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{N("a").CreateTable(N("c1").T("text").Foreign("other")), `CREATE TABLE a (c1 TEXT REFERENCES other)`},
		{N("a").CreateTable(N("c1").T("text").Foreign("other", "c1")), `CREATE TABLE a (c1 TEXT REFERENCES other (c1))`},
		{N("a").CreateTable(N("c1").T("text").Unique()), `CREATE TABLE a (c1 TEXT UNIQUE)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
