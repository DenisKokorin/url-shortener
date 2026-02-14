package urlshortenerservice_test

import (
	"log/slog"
	"testing"
	urlshortenerservice "url-shortener/internal/service"
	mock_urlshortenerservice "url-shortener/internal/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetShortURL(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mock_urlshortenerservice.NewMockStorage(ctrl)
			mockGenerator := mock_urlshortenerservice.NewMockAliasGenerator(ctrl)

			mockGenerator.EXPECT().Generate().Return(tt.alias, nil)
			mockStorage.EXPECT().SaveURL(gomock.Any(), tt.url, tt.alias).Return(nil)

			service := urlshortenerservice.NewURLShortenerService(slog.Default(), mockStorage, mockGenerator)

			res, err := service.GetShortURL(t.Context(), tt.url)

			require.NoError(t, err)

			require.Equal(t, tt.alias, res)
		})
	}
}

func TestGetLongURL(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mock_urlshortenerservice.NewMockStorage(ctrl)
			mockGenerator := mock_urlshortenerservice.NewMockAliasGenerator(ctrl)

			mockStorage.EXPECT().GetLongURL(gomock.Any(), tt.alias).Return(tt.url, nil)

			service := urlshortenerservice.NewURLShortenerService(slog.Default(), mockStorage, mockGenerator)

			res, err := service.GetLongURL(t.Context(), tt.alias)

			require.NoError(t, err)

			require.Equal(t, tt.url, res)
		})
	}
}
