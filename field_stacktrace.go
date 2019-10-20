package zapstackdriver

import (
	"runtime"
	"sync"
)

var (
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}
)

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
	AppendFloat(float64, bitSize int)
	AppendInt(int64)
	AppendString(string)
	AppendUint(uint64)
}

type stackFrameFilter interface {
	ShouldSkip(runtime.Frame) bool
}

func takeStacktrace(writer bufferWriter, filter stackFrameFilter, callers []uintptr) {
	i := 0
	shouldSkipFrame := true
	frames := runtime.CallersFrames(callers)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if shouldSkipFrame && filter.ShouldSkip(frame) {
			continue
		} else {
			shouldSkipFrame = false
		}

		if i != 0 {
			writer.AppendByte('\n')
		}

		i++

		writer.AppendString(frame.Function)
		writer.AppendByte('\n')
		writer.AppendByte('\t')
		writer.AppendString(frame.File)
		writer.AppendByte(':')
		writer.AppendInt(int64(frame.Line))
	}
}
