package zapstackdriver

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewProduction(options ...zap.Option) (*zap.Logger, error) {
	stdoutOption := NewCore(os.Stdout)
	options = append([]zap.Option{stdoutOption}, options...)
	return zap.NewProduction(options...)
}

func NewCore(w io.Writer) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			NewEncoder(),
			zapcore.AddSync(w),
			core,
		)
	})
}
