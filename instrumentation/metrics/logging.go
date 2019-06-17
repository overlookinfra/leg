package metrics

import (
	"context"

	logging "github.com/puppetlabs/insights-logging"
)

var (
	defaultLogger = logging.Builder().At("horsehead", "instrumentation", "metrics")
)

func log(ctx context.Context) logging.Logger {
	return defaultLogger.With(ctx).Build()
}
