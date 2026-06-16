package logger

import (
	"log"
	"os"
)

// Делаем интерфейс для логгера, чтобы можно было легко заменить реализацию,
// например для тестов или для использования другой библиотеки логирования,
// и не менять код в других местах, а только там где используется логгер
type Logger interface {
	Info(msg string)
	Error(msg string)
	Fatalf(format string, v ...interface{})
}

// SimpleLogger это простая реализация интерфейса Logger,
// которая использует стандартную библиотеку log для логирования сообщений в консоль,
type SimpleLogger struct {
	info  *log.Logger
	error *log.Logger
}

// Info logs an informational message.
func (s *SimpleLogger) Info(msg string) {
	s.info.Println(msg)
}

// Error logs an error message.
func (s *SimpleLogger) Error(msg string) {
	s.error.Println(msg)
}

// Fatalf logs a fatal message and exits.
func (s *SimpleLogger) Fatalf(format string, v ...interface{}) {
	s.error.Fatalf(format, v...)
}

// NewSimpleLogger creates a new SimpleLogger instance.
// Following Go idiom: constructor returns concrete type when creating simple objects.
func NewSimpleLogger() *SimpleLogger {
	return &SimpleLogger{
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}
