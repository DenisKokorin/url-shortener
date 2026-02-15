package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/config"
	gethandler "url-shortener/internal/handlers/get"
	savehandler "url-shortener/internal/handlers/save"
	logmiddleware "url-shortener/internal/middleware"
	urlshortenerservice "url-shortener/internal/service"
	"url-shortener/internal/storage"
	"url-shortener/pkg/generator"
	"url-shortener/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	storage, err := storage.GetStorageFromConfig(cfg)
	if err != nil {
		log.Error("failed to init storage", logger.ErrorLog(err))
		os.Exit(1)
	}

	generator := generator.NewAliasGenerator(cfg.AliasLength)

	service := urlshortenerservice.NewURLShortenerService(log, storage, generator)

	log.Info(
		"starting app",
		slog.String("env", cfg.Env),
	)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(logmiddleware.NewLogMiddleware(log))
	router.Use(middleware.Recoverer)

	router.Post("/", savehandler.NewSaveHandler(log, service))
	router.Get("/{alias}", gethandler.NewGetHandler(log, service))

	log.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.IdleTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {

		log.Error("failed to stop server", logger.ErrorLog(err))

		return
	}

	log.Info("server stopped")
}
