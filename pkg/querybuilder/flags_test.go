package querybuilder

import (
	"fmt"
	"testing"

	// Packages
	"github.com/stretchr/testify/assert"
)

func Test_Flags_000(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		In       any
		Expected string
	}{
		{temporary, `TEMPORARY`},
		{unlogged, `UNLOGGED`},
		{ifNotExists, `IF NOT EXISTS`},
		{temporary | ifNotExists, `TEMPORARY IF NOT EXISTS`},
	}
	for _, test := range tests {
		assert.Equal(test.Expected, fmt.Sprint(test.In))
	}
}
