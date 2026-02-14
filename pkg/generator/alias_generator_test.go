package generator_test

import (
	"testing"
	"url-shortener/pkg/generator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAliasGenerator(t *testing.T) {
	tests := []struct {
		name        string
		aliasLength int
	}{
		{
			name:        "len 5",
			aliasLength: 5,
		},
		{
			name:        "len 10",
			aliasLength: 10,
		},
		{
			name:        "len 20",
			aliasLength: 20,
		},
		{
			name:        "len 40",
			aliasLength: 40,
		},
		{
			name:        "len 0",
			aliasLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := generator.NewAliasGenerator(tt.aliasLength)
			s, err := gen.Generate()

			require.NoError(t, err)

			assert.Len(t, s, tt.aliasLength)
		})
	}
}
