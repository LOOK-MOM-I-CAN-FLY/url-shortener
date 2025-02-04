package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DB.URL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Выполнение миграций
	if err := repository.RunMigrations(cfg.DB.MigrationsPath, cfg.DB.URL); err != nil {
		logger.Fatal("Migrations failed", zap.Error(err))
	}

	// Инициализация слоёв
	repo := repository.NewPostgresRepository(pool, logger)
	cache := repository.NewRedisCache(cfg.Redis.Address, cfg.Redis.Password, logger)
	svc := service.NewShortenerService(repo, cache, logger)
	h := handler.NewHandler(svc, logger)

	// Настройка роутера
	r := chi.NewRouter()
	r.Use(handler.RequestID)
	r.Use(handler.LogMiddleware(logger))

	r.Post("/api/v1/urls", h.CreateShortURL)
	r.Get("/{shortCode}", h.Redirect)
	r.Get("/api/v1/urls/{code}", h.GetURLInfo)

	// Запуск сервера
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
	}
}
