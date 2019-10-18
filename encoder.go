package zapstackdriver

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func defaultConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		CallerKey:      "caller",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch l {
			case zapcore.DebugLevel:
				enc.AppendString("DEBUG")
			case zapcore.InfoLevel:
				enc.AppendString("INFO")
			case zapcore.WarnLevel:
				enc.AppendString("WARNING")
			case zapcore.ErrorLevel:
				enc.AppendString("ERROR")
			case zapcore.DPanicLevel:
				enc.AppendString("CRITICAL")
			case zapcore.PanicLevel:
				enc.AppendString("ALERT")
			case zapcore.FatalLevel:
				enc.AppendString("EMERGENCY")
			}
		},
		EncodeTime: zapcore.ISO8601TimeEncoder,
		LevelKey:   "severity",
		LineEnding: zapcore.DefaultLineEnding,
		MessageKey: "message",
		NameKey:    "logger",
		TimeKey:    "timestamp",
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
