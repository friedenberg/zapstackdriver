package zapstackdriver

import (
	"go.uber.org/zap"
)

func NewLogger(
	serviceContext FieldServiceContext,
	sourceReferences []FieldSourceReference,
	options ...zap.Option,
) (*Logger, error) {
	loggerConfig := NewConfig(serviceContext)
	logger, err := loggerConfig.Build(sourceReferences, options...)

	if err != nil {
		return nil, err
	}

	return logger, nil
}
