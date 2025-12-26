package middleware

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupLoggerMiddleware(logFilePath string, level string) (gin.HandlerFunc, error) {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	slogLevel := parseLogLevel(level)
	logger := slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slogLevel,
	}))

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		errorMsg := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		logAttrs := []any{
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.String("user-agent", userAgent),
			slog.Duration("latency", latency),
		}

		if errorMsg != "" {
			logAttrs = append(logAttrs, slog.String("error", errorMsg))
		}

		switch {
		case status >= 500:
			logger.Error("request failed", logAttrs...)
		case status >= 400:
			logger.Warn("client error", logAttrs...)
		default:
			logger.Info("incoming request", logAttrs...)
		}
	}, nil
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
