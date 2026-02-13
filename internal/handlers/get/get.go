package gethandler

import (
	"errors"
	"log/slog"
	"net/http"
	api "url-shortener/internal/handlers"
	urlshortenerservice "url-shortener/internal/service"
	"url-shortener/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewGetHandler(log *slog.Logger, service api.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "invalid request"))

			return
		}

		resURL, err := service.GetLongURL(r.Context(), alias)
		if errors.Is(err, urlshortenerservice.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, api.ErrorReponse(http.StatusNotFound, "not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", logger.ErrorLog(err))

			render.JSON(w, r, api.ErrorReponse(http.StatusInternalServerError, "internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		render.JSON(w, r, api.URLResponse(resURL))
	}
}
