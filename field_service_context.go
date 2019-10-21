package zapstackdriver

import (
	"go.uber.org/zap/zapcore"
)

type FieldServiceContext struct {
	Service string
	Version string
}

func (s FieldServiceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if s.Service != "" {
		enc.AddString("service", s.Service)
	}

	if s.Version != "" {
		enc.AddString("version", s.Version)
	}

	return nil
}
