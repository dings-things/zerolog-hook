package zerologhook_test

import (
	"io"
	"testing"

	hook "github.com/dings-things/zerolog-hook"
	"github.com/rs/zerolog"
)

/*
Benchmark results
This custom hook minimizes runtime calls and improves performance
compared to the built-in zerolog Caller() method.

> go test -bench=. -benchmem
goos: windows
goarch: amd64
pkg: test/zerolog/hook
cpu: 12th Gen Intel(R) Core(TM) i5-12600
BenchmarkBuiltinCaller-12                1124965              1014 ns/op             880 B/op          7 allocs/op
BenchmarkCustomCallerHook-12             3744060               308.7 ns/op           580 B/op          5 allocs/op
*/

// Dummy logging function using zerolog's built-in Caller()
func logWithBuiltinCaller() {
	logger := zerolog.New(io.Discard).With().Caller().Logger()
	logger.Info().Msg("builtin caller")
}

// Dummy logging function using the custom caller hook
func logWithCustomCallerHook() {
	logger := zerolog.New(io.Discard).With().
		Timestamp().Logger()
	logger.Hook(hook.NewCallerHook(true, true, true, true))
	logger.Info().Msg("custom hook")
}

// Benchmark using zerolog's built-in Caller()
func BenchmarkBuiltinCaller(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logWithBuiltinCaller()
	}
}

// Benchmark using the custom CallerHook
func BenchmarkCustomCallerHook(b *testing.B) {
	for i := 0; i < b.N; i++ {
		logWithCustomCallerHook()
	}
}
