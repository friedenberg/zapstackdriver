package zapstackdriver

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCore(w io.Writer) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			// ~~https://cloud.google.com/error-reporting/docs/formatting-error-messages~~
			// https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
			NewEncoder(),
			zapcore.AddSync(w),
			core,
		)
	})
}
