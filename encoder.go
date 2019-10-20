package zapstackdriver

import (
	"go.uber.org/zap/zapcore"
)

var levelToSeverityMap = map[zapcore.Level]string{
	zapcore.DebugLevel:  "DEBUG",
	zapcore.InfoLevel:   "INFO",
	zapcore.WarnLevel:   "WARNING",
	zapcore.ErrorLevel:  "ERROR",
	zapcore.DPanicLevel: "CRITICAL",
	zapcore.PanicLevel:  "ALERT",
	zapcore.FatalLevel:  "EMERGENCY",
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if value, ok := levelToSeverityMap[l]; ok {
		enc.AppendString(value)
	}
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeLevel:    levelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LevelKey:       "severity",
		LineEnding:     zapcore.DefaultLineEnding,
		MessageKey:     "message",
		NameKey:        "logger",
		TimeKey:        "timestamp",
		StacktraceKey:  "stacktrace",
	}
}
