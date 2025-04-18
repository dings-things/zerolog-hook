package zerologhook_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	hook "github.com/dings-things/zerolog-hook"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func captureLog(withFunc, withFile, withLine, withPkg bool) map[string]any {
	var buf bytes.Buffer
	multi := io.MultiWriter(&buf, os.Stdout)

	logger := zerolog.New(multi).With().Timestamp().Logger()
	logger = logger.Hook(hook.NewCallerHook(withFunc, withFile, withLine, withPkg))
	logger.Info().Msg("test log")

	var logLine map[string]any
	err := json.Unmarshal(buf.Bytes(), &logLine)
	if err != nil {
		panic("log is not valid JSON: " + buf.String())
	}

	return logLine
}

func TestCallerHook_OnlyFunc(t *testing.T) {
	log := captureLog(true, false, false, false)
	assert.Contains(t, log, "func")
	assert.NotContains(t, log, "file")
	assert.NotContains(t, log, "line")
	assert.NotContains(t, log, "pkg")
	assert.Equal(t, "captureLog", log["func"])
}

func TestCallerHook_OnlyFile(t *testing.T) {
	log := captureLog(false, true, false, false)
	assert.Contains(t, log, "file")
	assert.NotContains(t, log, "func")
	assert.NotContains(t, log, "line")
	assert.NotContains(t, log, "pkg")
	assert.IsType(t, "", log["file"])
}

func TestCallerHook_OnlyLine(t *testing.T) {
	log := captureLog(false, false, true, false)
	assert.Contains(t, log, "line")
	assert.NotContains(t, log, "func")
	assert.NotContains(t, log, "file")
	assert.NotContains(t, log, "pkg")
	assert.IsType(t, float64(0), log["line"])
}

func TestCallerHook_OnlyPkg(t *testing.T) {
	log := captureLog(false, false, false, true)
	assert.Contains(t, log, "pkg")
	assert.NotContains(t, log, "func")
	assert.NotContains(t, log, "file")
	assert.NotContains(t, log, "line")
	assert.IsType(t, "", log["pkg"])
}

func TestCallerHook_AllEnabled(t *testing.T) {
	log := captureLog(true, true, true, true)
	assert.Contains(t, log, "func")
	assert.Contains(t, log, "file")
	assert.Contains(t, log, "line")
	assert.Contains(t, log, "pkg")
}

func TestCallerHook_NoneEnabled(t *testing.T) {
	log := captureLog(false, false, false, false)
	assert.NotContains(t, log, "func")
	assert.NotContains(t, log, "file")
	assert.NotContains(t, log, "line")
	assert.NotContains(t, log, "pkg")
}

func TestCallerHook_FuncAndFile(t *testing.T) {
	log := captureLog(true, true, false, false)
	assert.Contains(t, log, "func")
	assert.Contains(t, log, "file")
	assert.NotContains(t, log, "line")
	assert.NotContains(t, log, "pkg")
}
