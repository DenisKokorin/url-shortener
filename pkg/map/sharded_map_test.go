package shardedmap_test

import (
	"testing"
	shardedmap "url-shortener/pkg/map"

	"github.com/stretchr/testify/require"
)

func TestShardedMap(t *testing.T) {
	tests := []struct {
		name  string
		key   any
		value any
	}{
		{
			name:  "string",
			key:   "test_alias",
			value: "http://www.google.com/",
		},
		{
			name:  "int",
			key:   5,
			value: 10,
		},
		{
			name:  "bool",
			key:   true,
			value: true,
		},
		{
			name:  "string key and int value",
			key:   "key",
			value: 5,
		},
		{
			name:  "int key and string value",
			key:   5,
			value: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sm := shardedmap.NewShardedMap()
			err := sm.Save(t.Context(), tt.key, tt.value)

			require.NoError(t, err)

			v, err := sm.Get(t.Context(), tt.key)

			require.NoError(t, err)

			require.Equal(t, tt.value, v)
		})
	}
}
