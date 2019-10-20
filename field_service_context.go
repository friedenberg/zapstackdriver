package zapstackdriver

import (
	"errors"

	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap/zapcore"
)

var (
	errServiceEmpty = errors.New("service name is required, but was empty")
	errVersionEmpty = errors.New("version name is required, but was empty")
)

type serviceContext struct {
	service string
	version string
}

func (s serviceContext) validate() error {
	var result *multierror.Error

	if s.service == "" {
		result = multierror.Append(result, errServiceEmpty)
	}

	if s.version == "" {
		result = multierror.Append(result, errVersionEmpty)
	}

	return result.ErrorOrNil()
}

func (s serviceContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	err := s.validate()

	if err != nil {
		return err
	}

	enc.AddString("service", s.service)
	enc.AddString("version", s.version)

	return nil
}
