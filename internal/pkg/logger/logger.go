package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

type Fields map[string]interface{}

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger()
}

func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
		"./logs/go-service/app.log",
	}
	config.ErrorOutputPaths = []string{
		"stderr",
		"./logs/go-service/error.log",
	}

	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "msg"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.StacktraceKey = "stacktrace"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// Create logs directory if it doesn't exist
	os.MkdirAll("./logs/go-service", 0755)

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	return &Logger{logger.Sugar()}
}

func Default() *Logger {
	return defaultLogger
}

func (l *Logger) WithFields(fields Fields) *Logger {
	var zapFields []interface{}
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	return &Logger{l.SugaredLogger.With(zapFields...)}
}

func (l *Logger) WithTraceID(traceID string) *Logger {
	return l.WithFields(Fields{"trace_id": traceID})
}

func (l *Logger) WithSpanID(spanID string) *Logger {
	return l.WithFields(Fields{"span_id": spanID})
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	if traceID := GetTraceIDFromContext(ctx); traceID != "" {
		return l.WithTraceID(traceID)
	}
	return l
}

// Context helpers
type contextKey string

const traceIDKey contextKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// Gin middleware for logging
func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			traceID := uuid.New().String()

			// Add trace ID to context
			param.Keys = map[string]interface{}{
				"trace_id": traceID,
			}

			logger := Default().WithFields(Fields{
				"trace_id": traceID,
				"method":   param.Method,
				"path":     param.Path,
				"status":   param.StatusCode,
				"duration": param.Latency.Milliseconds(),
				"ip":       param.ClientIP,
				"user_agent": param.Request.UserAgent(),
			})

			if param.StatusCode >= 400 {
				if param.StatusCode >= 500 {
					logger.Error("HTTP request completed with server error")
				} else {
					logger.Warn("HTTP request completed with client error")
				}
			} else {
				logger.Info("HTTP request completed")
			}

			return ""
		},
		Output: os.Stdout,
	})
}

// Recovery middleware with logging
func GinRecovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(os.Stderr, func(c *gin.Context, recovered interface{}) {
		traceID := ""
		if c.Keys != nil {
			if id, exists := c.Keys["trace_id"].(string); exists {
				traceID = id
			}
		}

		logger := Default().WithFields(Fields{
			"trace_id": traceID,
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"panic":    recovered,
		})

		logger.Error("Panic recovered")
		c.AbortWithStatus(500)
	})
}

// Global logger functions
func Info(msg string, fields ...Fields) {
	if len(fields) > 0 {
		defaultLogger.WithFields(fields[0]).Info(msg)
	} else {
		defaultLogger.Info(msg)
	}
}

func Error(msg string, fields ...Fields) {
	if len(fields) > 0 {
		defaultLogger.WithFields(fields[0]).Error(msg)
	} else {
		defaultLogger.Error(msg)
	}
}

func Warn(msg string, fields ...Fields) {
	if len(fields) > 0 {
		defaultLogger.WithFields(fields[0]).Warn(msg)
	} else {
		defaultLogger.Warn(msg)
	}
}

func Debug(msg string, fields ...Fields) {
	if len(fields) > 0 {
		defaultLogger.WithFields(fields[0]).Debug(msg)
	} else {
		defaultLogger.Debug(msg)
	}
}
