// Package middleware provides HTTP middleware functions.
package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/example/dormitory-management/pkg/logger"
)

// Logger creates a structured logging middleware.
func Logger(zapLogger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []logger.Field{
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.Int("status", status),
			logger.Duration("latency", latency),
			logger.String("client_ip", c.ClientIP()),
		}

		if query != "" {
			fields = append(fields, logger.String("query", query))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, logger.String("errors", c.Errors.String()))
		}

		if status >= 500 {
			zapLogger.Error("Server error", fields...)
		} else if status >= 400 {
			zapLogger.Warn("Client error", fields...)
		} else {
			zapLogger.Info("Request completed", fields...)
		}
	}
}

// Recovery creates a panic recovery middleware.
func Recovery(zapLogger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				zapLogger.Error("Panic recovered",
					logger.Any("error", err),
					logger.String("path", c.Request.URL.Path),
				)
				c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
