package zapstackdriver

import "go.uber.org/zap"

type Config struct {
	zap.Config
}

func (c *Config) Build(options ...zap.Option) (*Logger, error) {
	logger, err := c.Config.Build(options...)
	return &Logger{SugaredLogger: logger.Sugar()}, err
}
