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

func Test_PG_Constraints_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{Key(), `PRIMARY KEY`},
		{Key("a", "b"), `PRIMARY KEY (a,b)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}

func Test_PG_Constraints_001(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{Key().Foreign("a"), `REFERENCES a`},
		{Key().Foreign("a", "c1", "c2"), `REFERENCES a (c1,c2)`},
		{Key("a").Foreign("a"), `FOREIGN KEY (a) REFERENCES a`},
		{Key("a", "b").Foreign("a", "c1", "c2"), `FOREIGN KEY (a,b) REFERENCES a (c1,c2)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}

func Test_PG_Constraints_002(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{Key().Unique(), `UNIQUE`},
		{Key("a", "b").Unique(), `UNIQUE (a,b)`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}

func Test_PG_Constraints_003(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{Key().Foreign("other_table").OnDeleteCascade(), `REFERENCES other_table ON DELETE CASCADE`},
		{Key().Foreign("other_table").OnDeleteNoAction(), `REFERENCES other_table ON DELETE NO ACTION`},
		{Key().Foreign("other_table").OnDeleteRestrict(), `REFERENCES other_table ON DELETE RESTRICT`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
