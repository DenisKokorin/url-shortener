package gethandler

import (
	"errors"
	"log/slog"
	"net/http"
	api "url-shortener/internal/handlers"
	urlshortenerservice "url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewGetHandler(log *slog.Logger, service api.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "alias is empty"))
			return
		}

		resURL, err := service.GetLongURL(r.Context(), alias)
		if errors.Is(err, urlshortenerservice.ErrURLNotFound) {
			render.JSON(w, r, api.ErrorReponse(http.StatusNotFound, "not found"))
			return
		}
		if err != nil {
			render.JSON(w, r, api.ErrorReponse(http.StatusInternalServerError, "internal error"))
			return
		}

		render.JSON(w, r, api.URLResponse(resURL))
	}
}
