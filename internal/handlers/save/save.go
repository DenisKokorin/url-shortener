package savehandler

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	api "url-shortener/internal/handlers"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

func NewSaveHandler(log *slog.Logger, service api.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req api.Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "empty request"))
			return
		}
		if err != nil {
			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, "failed to decode request"))
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			render.JSON(w, r, api.ErrorReponse(http.StatusBadRequest, validateErr.Error()))
			return
		}

		alias, err := service.GetShortURL(r.Context(), req.URL)
		if err != nil {
			render.JSON(w, r, api.ErrorReponse(http.StatusInternalServerError, "internal error"))
			return
		}

		render.JSON(w, r, api.ResponseOK(alias))
	}
}
