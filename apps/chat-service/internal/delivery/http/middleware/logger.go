package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Специальная структура-обертка, которая "шпионит" за статус-кодом
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

// Переопределяем метод WriteHeader, чтобы запомнить статус-код
func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger возвращает мидлварь, которая логирует информацию о каждом HTTP-запросе
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Фиксируем время начала запроса
		start := time.Now()

		// 2. Передаем управление следующему хендлеру в цепочке (нашему бизнес-коду)
		// Оборачиваем оригинальный ResponseWriter, ставим дефолтный статус 200
		wrappedWriter := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Передаем обертку дальше
		next.ServeHTTP(wrappedWriter, r)

		// 3. После того как хендлер отработал, считаем время выполнения
		duration := time.Since(start)

		// 4. Пишем красивый структурированный лог
		slog.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrappedWriter.statusCode),
			slog.Duration("duration", duration),
		)
	})
}
