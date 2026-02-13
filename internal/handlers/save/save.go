package savehandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	api "url-shortener/internal/handlers"
	"url-shortener/pkg/logger"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func NewSaveHandler(log *slog.Logger, service api.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {

			//log.Error("request body is empty")

			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "empty request"))

			return
		}
		if err != nil {
			//log.Error("failed to decode request body", logger.ErrorLog(err))

			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "failed to decode request"))

			return
		}

		//log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", logger.ErrorLog(err))

			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, validateErr.Error()))

			return
		}

		alias, err := service.GetShortURL(r.Context(), req.URL)
		if err != nil {
			//log.Error("failed to get short url", logger.ErrorLog(err))

			render.JSON(w, r, api.ErrorReponse(http.StatusInternalServerError, "internal error"))

			return
		}

		//log.Info("url added", slog.String("alias", alias))

		render.JSON(w, r, api.ResponseOK(alias))
	}
}
