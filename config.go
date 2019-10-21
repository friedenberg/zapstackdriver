package zapstackdriver

import "go.uber.org/zap"

type Config struct {
	zap.Config
}

func (c *Config) Build(sourceReferences []FieldSourceReference, options ...zap.Option) (*Logger, error) {
	zaplogger, err := c.Config.Build(options...)

	if err != nil {
		return nil, err
	}

	logger := &Logger{
		SugaredLogger:    zaplogger.Sugar(),
		sourceReferences: sourceReferences,
	}

	return logger, nil
}

func NewConfig(serviceContext FieldServiceContext) *Config {
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "json",
		EncoderConfig:     NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		InitialFields:     map[string]interface{}{"serviceContext": serviceContext},
	}

	return &Config{Config: config}
}
