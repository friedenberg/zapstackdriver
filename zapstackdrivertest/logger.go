package zapstackdrivertest

import (
	"github.com/friedenberg/zapstackdriver"
	"go.uber.org/zap/zaptest"
)

func NewLogger(t zaptest.TestingT, opts ...zaptest.LoggerOption) *zapstackdriver.Logger {
	l := zaptest.NewLogger(t, opts...)
	return &zapstackdriver.Logger{SugaredLogger: l.Sugar()}
}
