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
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, api.ErrorReponse("alias is empty"))
			return
		}

		resURL, err := service.GetLongURL(r.Context(), alias)
		if errors.Is(err, urlshortenerservice.ErrURLNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, api.ErrorReponse("not found"))
			return
		}
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, api.ErrorReponse("internal error"))
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, api.URLResponse(resURL))
	}
}
