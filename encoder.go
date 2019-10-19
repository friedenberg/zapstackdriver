package zapstackdriver

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
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

func defaultConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeLevel:    levelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LevelKey:       "severity",
		LineEnding:     zapcore.DefaultLineEnding,
		MessageKey:     "message",
		NameKey:        "logger",
		TimeKey:        "timestamp",
	}
}

type encoder struct {
	zapcore.Encoder
}

func newEncoderWithConfig(config zapcore.EncoderConfig) *encoder {
	return &encoder{
		zapcore.NewJSONEncoder(config),
	}
}

func NewEncoder() *encoder {
	return newEncoderWithConfig(defaultConfig())
}

func (e *encoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if entry.Level >= zapcore.ErrorLevel {
		entry, fields = MakeErrorEntry(entry, fields)
	} else {
		sourceLocationField := zap.Object(
			"sourceLocation",
			&fieldSourceLocation{caller{entry.Caller}},
		)

		fields = append(fields, sourceLocationField)
	}

	return e.Encoder.EncodeEntry(entry, fields)
}

func (e *encoder) Clone() zapcore.Encoder {
	return &encoder{e.Encoder.Clone()}
}
