# OpenTelemetry instrumentation for go-kv

## Installation

```bash
go get github.com/kv/go-kv/extra/kvotel/v9
```

## Usage

Tracing is enabled by adding a hook:

```go
import (
    "github.com/kv/go-kv/v9"
    "github.com/kv/go-kv/extra/kvotel/v9"
)

rdb := rdb.NewClient(&rdb.Options{...})

// Enable tracing instrumentation.
if err := kvotel.InstrumentTracing(rdb); err != nil {
	panic(err)
}

// Enable metrics instrumentation.
if err := kvotel.InstrumentMetrics(rdb); err != nil {
	panic(err)
}
```

See [example](../../example/otel) and
[Monitoring Go KV Performance and Errors](https://kv.uptrace.dev/guide/go-kv-monitoring.html)
for details.
