package zapstackdriver

import (
	"errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type serviceContext struct {
	service string
	version string
}

func (s serviceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if s.service == "" {
		return errors.New("service name is mandatory")
	}

	enc.AddString("service", s.service)
	enc.AddString("version", s.version)

	return nil
}

func WithContext(service string, version string) zap.Option {
	serviceCtx := serviceContext{
		service: service,
		version: version,
	}

	return zap.Fields(
		zap.Object("serviceContext", serviceCtx),
	)
}
