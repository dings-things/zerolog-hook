package zerologhook

import (
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const skip = 4 // skip levels: hook → zerolog → app → actual caller

// CallerHook is a Zerolog hook that automatically adds caller information
// (function name, package name, filename, and line number) to the log entry.
// Each field can be configured individually at the time of creation.
type CallerHook struct {
	WithFunc bool // Include function name
	WithFile bool // Include filename
	WithLine bool // Include line number
	WithPkg  bool // Include package name
}

// funcNameCache is a concurrency-safe map for caching function info by PC (program counter)
var funcNameCache sync.Map // map[uintptr]string

// NewCallerHook creates a new instance of CallerHook
//   - withFunc: whether to include the function name
//   - withFile: whether to include the filename
//   - withLine: whether to include the line number
//   - withPkg: whether to include the package name
//
// example:
//
//	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
//	logger = logger.Hook(NewCallerHook(true, true, true, true))
//	logger.Info().Msg("test log")
//	// Output: {"level":"info","time":...,"func":"captureLog","file":"caller_hook_test.go","line":...,"pkg":"loggerhook_test","message":"test log"}
func NewCallerHook(withFunc, withFile, withLine, withPkg bool) CallerHook {
	return CallerHook{
		WithFunc: withFunc,
		WithFile: withFile,
		WithLine: withLine,
		WithPkg:  withPkg,
	}
}

// Run is called when the hook is executed (implements zerolog.Hook interface)
func (h CallerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// Collect the program counter for the caller
	var pcs [1]uintptr
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return
	}
	pc := pcs[0]

	// Retrieve or store the full function path from/to the cache
	fnVal, ok := funcNameCache.Load(pc)
	var fn string
	if ok {
		fn = fnVal.(string)
	} else {
		frames := runtime.CallersFrames(pcs[:])
		frame, _ := frames.Next()
		fn = frame.Function
		funcNameCache.Store(pc, fn)
	}

	// Split the full function path into function and package names
	split := strings.Split(fn, "/")
	fullFunc := split[len(split)-1]

	funcName := fullFunc
	pkgName := ""
	if i := strings.LastIndex(fullFunc, "."); i != -1 {
		pkgName = fullFunc[:i]
		funcName = fullFunc[i+1:]
	}

	// Add function and package name to log fields if enabled
	if h.WithFunc {
		e.Str("func", funcName)
	}
	if h.WithPkg {
		e.Str("pkg", pkgName)
	}

	// Optionally include filename and line number
	if h.WithFile || h.WithLine {
		frames := runtime.CallersFrames(pcs[:])
		frame, _ := frames.Next()

		if h.WithFile {
			_, filename := path.Split(frame.File)
			e.Str("file", filename)
		}
		if h.WithLine {
			e.Int("line", frame.Line)
		}
	}
}
