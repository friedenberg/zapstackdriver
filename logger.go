package zapstackdriver

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultCallerSkipOffset = 2
)

type Logger struct {
	*zap.SugaredLogger
	callerSkipOffset    int
	errorContextRequest *FieldHTTPRequest
	errorContextUser    string
	sourceReferences    []FieldSourceReference
}

func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger.With(args...),
		callerSkipOffset:    l.callerSkipOffset,
		errorContextRequest: l.errorContextRequest,
		errorContextUser:    l.errorContextUser,
		sourceReferences:    l.sourceReferences,
	}
}

func (l *Logger) WithCallerSkipOffset(offset int) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger,
		callerSkipOffset:    offset,
		errorContextRequest: l.errorContextRequest,
		errorContextUser:    l.errorContextUser,
		sourceReferences:    l.sourceReferences,
	}
}

func (l *Logger) WithRequest(request *http.Request) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger,
		callerSkipOffset:    l.callerSkipOffset,
		errorContextRequest: &FieldHTTPRequest{Request: *request},
		errorContextUser:    l.errorContextUser,
		sourceReferences:    l.sourceReferences,
	}
}

func (l *Logger) WithUser(user string) *Logger {
	return &Logger{
		SugaredLogger:       l.SugaredLogger,
		callerSkipOffset:    l.callerSkipOffset,
		errorContextRequest: l.errorContextRequest,
		errorContextUser:    user,
		sourceReferences:    l.sourceReferences,
	}
}

func (l *Logger) SetResponseStatusCode(statusCode int) {
	if l.errorContextRequest != nil {
		l.errorContextRequest.ResponseStatusCode = statusCode
	}
}

func (l *Logger) CallerSkipOffset() int {
	return defaultCallerSkipOffset + l.callerSkipOffset
}

func (l *Logger) callerWithAddedOffset(offset int) Caller {
	caller := zapcore.NewEntryCaller(runtime.Caller(l.CallerSkipOffset() + offset))
	return Caller{EntryCaller: caller}
}

func (l *Logger) withNonErrorContext() *zap.SugaredLogger {
	caller := l.callerWithAddedOffset(1)
	fieldSourceLocation := FieldSourceLocation{Caller: caller}

	return l.SugaredLogger.With(SDKeySourceLocation, fieldSourceLocation)
}

func (l *Logger) withErrorContext() *zap.SugaredLogger {
	const errorEntryKeyContext = "context"

	caller := l.callerWithAddedOffset(1)

	errorContext := FieldErrorContext{
		HTTPRequest:      l.errorContextRequest,
		User:             l.errorContextUser,
		ReportLocation:   FieldReportLocation{Caller: caller},
		SourceReferences: l.sourceReferences,
	}

	return l.SugaredLogger.With(
		SDKeyType, SDValueErrorType,
		errorEntryKeyContext, errorContext,
	)
}

func (l *Logger) appendStacktraceString(msg string, offset int) string {
	stacktrace := MakeDefaultStacktrace(l.CallerSkipOffset() + 1 + offset)
	return fmt.Sprintf("%s\n%s", msg, stacktrace)
}

func (l *Logger) appendStacktrace(v interface{}, offset int) interface{} {
	if msg, ok := v.(string); ok {
		return l.appendStacktraceString(msg, offset+1)
	}

	return v
}

func (l *Logger) Debug(args ...interface{}) {
	l.withNonErrorContext().Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.withNonErrorContext().Info(args...)
}

func (l *Logger) Error(args ...interface{}) {
	args[0] = l.appendStacktrace(args[0], 1)
	l.withErrorContext().Error(args...)
}

func (l *Logger) DPanic(args ...interface{}) {
	args[0] = l.appendStacktrace(args[0], 1)
	l.withErrorContext().DPanic(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	args[0] = l.appendStacktrace(args[0], 1)
	l.withErrorContext().Panic(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	args[0] = l.appendStacktrace(args[0], 1)
	l.withErrorContext().Fatal(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.withNonErrorContext().Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.withNonErrorContext().Infof(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	template = l.appendStacktraceString(template, 1)
	l.withErrorContext().Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	template = l.appendStacktraceString(template, 1)
	l.withErrorContext().DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	template = l.appendStacktraceString(template, 1)
	l.withErrorContext().Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	template = l.appendStacktraceString(template, 1)
	l.withErrorContext().Fatalf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.withNonErrorContext().Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.withNonErrorContext().Infow(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	msg = l.appendStacktraceString(msg, 1)
	l.withErrorContext().Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	msg = l.appendStacktraceString(msg, 1)
	l.withErrorContext().DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	msg = l.appendStacktraceString(msg, 1)
	l.withErrorContext().Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	msg = l.appendStacktraceString(msg, 1)
	l.withErrorContext().Fatalw(msg, keysAndValues...)
}

func (l *Logger) StdLogger() *log.Logger {
	return zap.NewStdLog(l.SugaredLogger.Desugar())
}
