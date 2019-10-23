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
		enc.AddString(SDKeyServiceContextService, s.Service)
	}

	if s.Version != "" {
		enc.AddString(SDKeyServiceContextVersion, s.Version)
	}

	return nil
}
