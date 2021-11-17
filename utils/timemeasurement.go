package utils

import (
	"fmt"
	"go.uber.org/zap"
	"time"
)

type TimeMeasurement struct {
	logger *zap.Logger
}

func NewTimeMeasurement(logger *zap.Logger) *TimeMeasurement {
	return &TimeMeasurement{
		logger: logger,
	}
}

//use with defer
func (t *TimeMeasurement) TimeTrack(start time.Time, contextName string, methodName string) {
	elapsed := time.Since(start)
	t.logger.Info(fmt.Sprintf("%s %s execution took: %s", contextName, methodName, elapsed),
		zap.String("ContextName", contextName),
		zap.String("ExecutingMethod", methodName),
		zap.Int64("ExecutionTook", elapsed.Nanoseconds()),
	)
}
