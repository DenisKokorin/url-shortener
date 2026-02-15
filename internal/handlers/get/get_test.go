package gethandler_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	api "url-shortener/internal/handlers"
	gethandler "url-shortener/internal/handlers/get"
	mock_api "url-shortener/internal/handlers/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	tests := []struct {
		name      string
		alias     string
		url       string
		respError string
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},
		{
			name:      "Url not found",
			alias:     "test_alias",
			url:       "",
			respError: "url not found",
		},
		{
			name:      "Internal error",
			alias:     "test_alias",
			url:       "",
			respError: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_api.NewMockService(ctrl)

			mockService.EXPECT().GetLongURL(gomock.Any(), tt.alias).Return(tt.url, nil)

			handler := gethandler.NewGetHandler(slog.Default(), mockService)

			r := chi.NewRouter()
			r.Get("/{alias}", handler)

			ts := httptest.NewServer(r)
			defer ts.Close()

			path := "/" + tt.alias
			req, err := http.NewRequest(http.MethodGet, path, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			require.NoError(t, err)

			body := rr.Body.String()

			var resp api.LongURLReponse

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tt.url, resp.URL)
		})
	}
}
