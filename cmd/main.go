package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/config"
	"url-shortener/pkg/logger"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info(
		"starting app",
		slog.String("env", cfg.Env),
	)

	router := chi.NewRouter()

	router.Post("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("created")) })
	router.Get("/{alias}", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

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
}
