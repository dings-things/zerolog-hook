# zerolog-hook
[![Go Reference](https://pkg.go.dev/badge/github.com/dings-things/zerolog-hook.svg)](https://pkg.go.dev/github.com/dings-things/zerolog-hook)
[![Go Report Card](https://goreportcard.com/badge/github.com/dings-things/zerolog-hook)](https://goreportcard.com/report/github.com/dings-things/zerolog-hook)
[![License: MIT](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/dings-things/zerolog-hook/blob/main/LICENSE)
[![Test](https://github.com/dings-things/zerolog-hook/actions/workflows/test.yml/badge.svg)](https://github.com/dings-things/zerolog-hook/actions/workflows/test.yml)

A lightweight and performant caller info hook for [zerolog](https://github.com/rs/zerolog) that enriches log events with function name, package name, file name, and line number — with **minimal runtime overhead**.

> ⚡ Up to 3x faster than the built-in `.Caller()` method in zerolog. See [Benchmarks](#benchmarks).

---

## Features

✅ Adds structured caller information:  
&nbsp;&nbsp;&nbsp;&nbsp;• Function name  
&nbsp;&nbsp;&nbsp;&nbsp;• Package name  
&nbsp;&nbsp;&nbsp;&nbsp;• File name  
&nbsp;&nbsp;&nbsp;&nbsp;• Line number  

✅ Individually configurable fields

✅ Zero external dependencies

✅ Thread-safe caching for runtime function resolution

---

## Installation

```bash
go get github.com/dings-things/zerolog-hook
```

---

## Quick Start

```go
import (
    "os"

    "github.com/rs/zerolog"
    zerohook "github.com/dings-things/zerolog-hook"
)

func main() {
    logger := zerolog.New(os.Stdout).
        With().
        Timestamp().
        Logger().
        Hook(zerohook.NewCallerHook(
            true,  // WithFunc
            true,  // WithFile
            true,  // WithLine
            true,  // WithPkg
        ))

    logger.Info().Msg("Hello from logger")
}
```

Sample output:

```json
{
  "level": "info",
  "time": "2025-04-18T12:00:00Z",
  "func": "main",
  "file": "main.go",
  "line": 14,
  "pkg": "main",
  "message": "Hello from logger"
}
```

---

## Use Cases

- When you need **structured log context** for tracing and debugging
- To **replace zerolog’s `.Caller()`** with a faster and more flexible alternative
- In **performance-sensitive** environments (high-frequency log generation)

---

## Benchmarks

Tested on:

```
OS:     Windows 10
CPU:    12th Gen Intel(R) Core(TM) i5-12600
Go:     go1.21+
```

```bash
go test -bench=. -benchmem
```

| Method                | Ops/sec   | ns/op  | B/op | allocs/op |
|----------------------|-----------|--------|------|-----------|
| `zerolog.Caller()`   | ~1.12 M   | 1014   | 880  | 7         |
| `zerolog-hook`       | ~3.74 M   | 308.7  | 580  | 5         |

✅ **~3x faster and ~34% less memory** usage than zerolog's built-in `.Caller()`

---

## API

```go
func NewCallerHook(
    withFunc bool,
    withFile bool,
    withLine bool,
    withPkg bool,
) zerolog.Hook
```

- `withFunc`: include `"func"` field
- `withFile`: include `"file"` field
- `withLine`: include `"line"` field
- `withPkg`: include `"pkg"` field

---

## Testing

```bash
go test -v ./...
```

Unit and benchmark tests are included to validate correctness and performance.

---

## License

MIT © 2025 [dings-things](https://github.com/dings-things)

---

## Contributing

Pull requests are welcome!  
Please open an issue first for major changes or discussions.