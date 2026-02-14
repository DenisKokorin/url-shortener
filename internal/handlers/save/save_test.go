package savehandler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	api "url-shortener/internal/handlers"
	mock_api "url-shortener/internal/handlers/mocks"
	savehandler "url-shortener/internal/handlers/save"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	tests := []struct {
		name       string
		alias      string
		url        string
		httpStatus int
		respError  string
		mockError  error
	}{
		{
			name:       "Success",
			alias:      "test_alias",
			url:        "http://www.google.com/",
			httpStatus: 200,
		},
		{
			name:       "Empty URL",
			url:        "",
			alias:      "",
			httpStatus: 400,
			respError:  "url is empty",
		},
		{
			name:       "Invalid URL",
			url:        "invalid URL",
			alias:      "",
			httpStatus: 400,
			respError:  "invalid request",
		},
		{
			name:       "SaveURL Error",
			alias:      "",
			url:        "https://google.com",
			httpStatus: 500,
			respError:  "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock_api.NewMockService(ctrl)

			if tt.httpStatus != 400 {
				mockService.EXPECT().GetShortURL(gomock.Any(), tt.url).Return(tt.alias, tt.mockError)
			}

			handler := savehandler.NewSaveHandler(slog.Default(), mockService)

			input := fmt.Sprintf(`{"url": "%s"}`, tt.url)

			req, err := http.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			var resp api.Response

			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))

			require.Equal(t, tt.alias, resp.Alias)
		})
	}
}
