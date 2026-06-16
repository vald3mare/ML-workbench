package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/vald3mare/ML-workbench/internal/logger"
)

func main() {
	// создаем логгер для всего приложения, который будет использоваться в разных местах, например в сервисе и в контроллере,
	// а реализация будет в другом месте, например в postgres или в памяти
	l := logger.NewSimpleLogger()

	// загружаем переменные окружения из .env файла, если он есть, и логируем предупреждение, если его нет,
	// но не останавливаем приложение, так как переменные окружения могут быть установлены другими способами, например через систему управления конфигурацией или
	// через переменные окружения операционной системы
	if err := godotenv.Load(); err != nil {
		l.Error(fmt.Sprintf("Warning: could not load .env file: %v", err))
	}

	// запускаем приложение и логируем фатальную ошибку, если она произойдет, и завершаем приложение с кодом 1,
	// так как это означает, что приложение не может продолжать работать из-за ошибки
	if err := run(l); err != nil {
		l.Fatalf("Fatal error: %v", err)
	}
}

// run реализует основную логику приложения, которая включает в себя настройку сервера,
// обработку сигналов для graceful shutdown и запуск сервера в отдельной горутине,
func run(l logger.Logger) error {
	// валидация переменных окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// инициализация сервера
	srv := &http.Server{
		Addr: ":8080",
	}

	// создаем канал для получения ошибок от сервера, он будет использоваться для обработки ошибок в основном потоке,
	errCh := make(chan error, 1)

	// запуск сервера в отдельной горутине, чтобы он не блокировал основной поток, и логируем информацию о запуске сервера,
	go func() {
		l.Info("Server starting on " + srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	// graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// ждем сигнала для graceful shutdown или ошибки от сервера, и обрабатываем их соответственно,
	// логируя информацию о получении сигнала и ошибках,
	select {
	case err := <-errCh:
		if err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-sigCh:
		l.Info(fmt.Sprintf("Received signal: %v", sig))
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}
	}
	return nil
}
