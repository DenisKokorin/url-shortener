package memory_test

import (
	"testing"
	"url-shortener/internal/storage/memory"

	"github.com/stretchr/testify/require"
)

func TestMemoryStorage(t *testing.T) {
	tests := []struct {
		name  string
		alias string
		url   string
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "http://www.google.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			storage := memory.NewMemoryStorage()

			err := storage.SaveURL(t.Context(), tt.url, tt.alias)

			require.NoError(t, err)

			res, err := storage.GetLongURL(t.Context(), tt.alias)

			require.NoError(t, err)

			require.Equal(t, tt.url, res)
		})
	}
}
