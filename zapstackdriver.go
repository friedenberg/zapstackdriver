package zapstackdriver

import (
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"
)

func NewConfig(initialFields ...validatedField) (*Config, error) {
	var errors *multierror.Error

	for _, field := range initialFields {
		errors = multierror.Append(errors, field.validate())
	}

	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:   false,
		DisableCaller: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	return &Config{Config: config}, errors.ErrorOrNil()
}
