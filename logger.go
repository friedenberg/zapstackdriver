package zapstackdriver

import (
	"net/http"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultCallerSkipOffset      = 2
	stackdriverKeySourceLocation = "sourceLocation"
	stackdriverKeyType           = "@type"
	stackdriverValueErrorType    = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"
)

type Logger struct {
	*zap.SugaredLogger
	errorContextRequest *HttpRequest
	errorContextUser    string
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{SugaredLogger: l.SugaredLogger.With(args...)}
}

func (l *Logger) WithRequest(request *http.Request) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger,
		errorContextRequest: &HttpRequest{Request: request},
		errorContextUser:    l.errorContextUser,
	}
}

func (l *Logger) WithUser(user string) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger,
		errorContextRequest: l.errorContextRequest,
		errorContextUser:    user,
	}
}

func (l *Logger) SetResponseStatusCode(statusCode int) {
	if l.errorContextRequest != nil {
		l.errorContextRequest.ResponseStatusCode = statusCode
	}
}

func (l *Logger) CallerSkipOffset() int {
	return defaultCallerSkipOffset
}

func (l *Logger) callerWithAddedOffset(offset int) Caller {
	caller := zapcore.NewEntryCaller(runtime.Caller(l.CallerSkipOffset() + offset))
	return Caller{EntryCaller: caller}
}

func (l *Logger) withNonErrorContext() *zap.SugaredLogger {

	caller := l.callerWithAddedOffset(1)
	fieldSourceLocation := FieldSourceLocation{Caller: caller}

	return l.SugaredLogger.With(stackdriverKeySourceLocation, &fieldSourceLocation)
}

func (l *Logger) withErrorContext() *zap.SugaredLogger {
	const errorEntryKeyContext = "context"

	caller := l.callerWithAddedOffset(1)

	errorContext := FieldErrorContext{
		HttpRequest:    l.errorContextRequest,
		User:           l.errorContextUser,
		ReportLocation: &FieldReportLocation{Caller: caller},
	}

	return l.SugaredLogger.With(
		stackdriverKeyType, stackdriverValueErrorType,
		errorEntryKeyContext, &errorContext,
	)
}

func (l *Logger) Error(args ...interface{}) {
	l.withErrorContext().Error(args...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.withErrorContext().DPanic(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.withErrorContext().Panic(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.withErrorContext().Fatal(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.withErrorContext().Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.withErrorContext().DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.withErrorContext().Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.withErrorContext().Fatalf(template, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.withErrorContext().Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.withErrorContext().DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.withErrorContext().Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.withErrorContext().Fatalw(msg, keysAndValues...)
}
