package log

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger, _ = zap.Config{
	Encoding:    "json",
	Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
	OutputPaths: []string{"stdout"},
	EncoderConfig: zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "@timestamp",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	},
}.Build()

func GetLogger() *zap.Logger {
	return logger
}

type HttpContext struct {
	echo.Context
}

func CreateContextFromEcho(ctx context.Context, echo echo.Context) context.Context {
	httpContext := HttpContext{echo}
	return context.WithValue(ctx, "httpContext", httpContext)
}

func CreateZapFieldsFromContext(ctx context.Context) []zap.Field {
	httpContext, ok := ctx.Value("httpContext").(HttpContext)
	if ok {
		return CreateZapFieldsFromEcho(httpContext)
	}
	return nil
}

func CreateZapFieldsFromEcho(ctx echo.Context, additionalFields ...zap.Field) []zap.Field {
	fields := make([]zap.Field, 0)
	if correlationId := ctx.Request().Header.Get("x-correlation-id"); correlationId != "" {
		fields = append(fields, zap.String("fields.CorrelationId", correlationId))
	}
	if userId := ctx.Request().Header.Get("user-id"); userId != "" {
		fields = append(fields, zap.String("fields.UserId", userId))
	}
	if clientId := ctx.Request().Header.Get("client-id"); clientId != "" {
		fields = append(fields, zap.String("fields.ClientId", clientId))
	}
	if tenantId := ctx.Request().Header.Get("tenant-id"); tenantId != "" {
		fields = append(fields, zap.String("fields.TenantId", tenantId))
	}
	if appKey := ctx.Request().Header.Get("app-key"); appKey != "" {
		fields = append(fields, zap.String("fields.AppKey", appKey))
	}
	if email := ctx.Request().Header.Get("email"); email != "" {
		fields = append(fields, zap.String("fields.Email", email))
	}
	for _, k := range additionalFields {
		fields = append(fields, k)
	}
	return fields
}
