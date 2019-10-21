package zapstackdriver

import (
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap/buffer"
)

var (
	_pool           = buffer.NewPool()
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}

	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func MakeStacktrace(offset int, filter StacktraceFilter, formatter StacktraceFormatter) string {
	pc := getProgramCounters(offset + 1)

	buffer := _pool.Get()
	defer buffer.Free()

	field := FieldStacktrace{
		callers:   pc.pcs,
		writer:    buffer,
		filter:    filter,
		formatter: formatter,
	}

	field.Format()

	return buffer.String()
}

func MakeDefaultStacktrace(offset int) string {
	return MakeStacktrace(
		offset+1,
		&StacktraceDefaultFilter{},
		&StacktraceDefaultFormatter{},
	)
}

func getProgramCounters(offset int) *programCounters {
	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int

	for {
		numFrames = runtime.Callers(offset, programCounters.pcs)

		if numFrames < len(programCounters.pcs) {
			break
		}
		// Don't put the too-short counter slice back into the pool; this lets
		// the pool adjust if we consistently take deep stacktraces.
		programCounters = newProgramCounters(len(programCounters.pcs) * offset)
	}

	return programCounters
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}

type bufferWriter interface {
	AppendBool(bool)
	AppendByte(byte)
	AppendFloat(float64, int)
	AppendInt(int64)
	AppendString(string)
	AppendUint(uint64)
}

type StacktraceFilter interface {
	ShouldSkip(runtime.Frame) bool
}

type StacktraceFormatter interface {
	BeginFormatting(writer bufferWriter)
	FormatFrame(writer bufferWriter, frame runtime.Frame, index int)
	JoinFrames(writer bufferWriter, firstIndex int, secondIndex int)
	EndFormatting(writer bufferWriter, index int)
}

type FieldStacktrace struct {
	callers   []uintptr
	writer    bufferWriter
	filter    StacktraceFilter
	formatter StacktraceFormatter
}

func (f *FieldStacktrace) Format() {
	i := 0
	shouldSkipFrame := true
	frames := runtime.CallersFrames(f.callers)

	f.formatter.BeginFormatting(f.writer)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if shouldSkipFrame && f.filter.ShouldSkip(frame) {
			continue
		} else {
			shouldSkipFrame = false
		}

		f.formatter.FormatFrame(f.writer, frame, i)
		oldIndex := i
		i++

		if more {
			f.formatter.JoinFrames(f.writer, oldIndex, i)
		}
	}

	f.formatter.EndFormatting(f.writer, i)
}

type StacktraceDefaultFilter struct{}

func (f *StacktraceDefaultFilter) ShouldSkip(runtime.Frame) bool {
	return false
}

type StacktraceDefaultFormatter struct{}

func (f *StacktraceDefaultFormatter) BeginFormatting(writer bufferWriter) {}

func (f *StacktraceDefaultFormatter) FormatFrame(writer bufferWriter, frame runtime.Frame, index int) {
	writer.AppendString(" at ")
	fileWithoutBase := strings.TrimPrefix(frame.File, basepath)
	writer.AppendString(fileWithoutBase)
	writer.AppendByte(':')
	writer.AppendInt(int64(frame.Line))
	writer.AppendByte(':')
	writer.AppendInt(int64(frame.Line))
}

func (f *StacktraceDefaultFormatter) JoinFrames(writer bufferWriter, firstIndex int, secondIndex int) {
	writer.AppendByte('\n')
}

func (f *StacktraceDefaultFormatter) EndFormatting(writer bufferWriter, lastIndex int) {}
