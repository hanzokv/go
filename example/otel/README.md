# go-kv OpenTelemetry Monitoring with Uptrace

This example demonstrates how to instrument and monitor KV operations in Go applications using
OpenTelemetry and [Uptrace](https://github.com/uptrace/uptrace), providing comprehensive
observability into your KV performance and operations.

## Overview

This integration provides:

- **Distributed tracing** for KV operations
- **Performance monitoring** with latency and throughput metrics
- **Error tracking** and debugging capabilities
- **Visual dashboards** for KV health monitoring
- **Production-ready** observability stack with Docker

## Prerequisites

- Go 1.19+
- Docker and Docker Compose
- Basic understanding of KV and OpenTelemetry

## Quick Start

### 1. Clone and Navigate

```bash
git clone https://github.com/kv/go-kv.git
cd example/otel
```

### 2. Start the Monitoring Stack

Launch KV and Uptrace services:

```bash
docker compose up -d
```

This starts:

- KV server on `localhost:6379`
- Uptrace APM on `http://localhost:14318`

### 3. Verify Services

Check that Uptrace is running properly:

```bash
docker compose logs uptrace
```

Look for successful startup messages without errors.

### 4. Run the Example

Execute the instrumented KV client:

```bash
go run client.go
```

You should see output similar to:

```
trace: http://localhost:14318/traces/ee029d8782242c8ed38b16d961093b35
```

Click the trace URL to view detailed operation traces in Uptrace.

![KV trace visualization](./image/kv-trace.png)

### 5. Explore the Dashboard

Open the Uptrace UI at [http://localhost:14318](http://localhost:14318/metrics/1) to explore:

- **Traces**: Individual KV operation details
- **Metrics**: Performance statistics and trends
- **Logs**: Application and system logs
- **Service Map**: Visual representation of dependencies

## Advanced Monitoring Setup

### KV Performance Metrics

For production environments, enable comprehensive KV monitoring by installing the OpenTelemetry
Collector:

The [OpenTelemetry Collector](https://uptrace.dev/opentelemetry/collector) acts as a telemetry agent
that:

- Pulls performance metrics directly from KV
- Collects system-level statistics
- Forwards data to Uptrace via OTLP protocol

When configured, Uptrace automatically generates a KV dashboard:

![KV performance dashboard](./image/metrics.png)

### Key Metrics Monitored

- **Connection Statistics**: Active connections, connection pool utilization
- **Command Performance**: Operation latency, throughput, error rates
- **Memory Usage**: Memory consumption, key distribution
- **Replication Health**: Master-slave sync status and lag

### Logs and Debugging

View service logs:

```bash
# All services
docker compose logs

# Specific service
docker compose logs kv
docker compose logs uptrace
```

## Additional Resources

- [Complete go-kv Monitoring Guide](https://kv.uptrace.dev/guide/go-kv-monitoring.html)
- [OpenTelemetry Go Instrumentation](https://uptrace.dev/get/opentelemetry-go/tracing)
- [Uptrace Open Source APM](https://uptrace.dev/get/hosted/open-source-apm)
