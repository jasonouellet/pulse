# [ADR-006]: Observability Architecture Strategy via OpenTelemetry (OTEL)

## Status

Accepted

## Context

Project PULSE operates as a modular monolith in Go (with a planned microservices transition), a React SPA frontend, PostgreSQL databases, and Redis caches. In distributed systems and complex tournament management scenarios (e.g., live score tracking, multi-field schedule updates), full visibility into system behavior, latency bottlenecks, and error propagation is critical.

We require a unified observability standard covering the Three Pillars of Observability:

1. **Traces**: End-to-end request flows across Browser -> API -> PostgreSQL/Redis.
2. **Metrics**: Performance metrics (HTTP request rates, Go runtime stats, DB connection pool health, React Web Vitals).
3. **Logs**: Structured JSON logging correlated with Trace IDs.

## Decision

We adopt **OpenTelemetry (OTEL)** as the vendor-neutral telemetry standard across all layers of the platform (Frontend, Backend, Database).

### 1. Backend Layer (Go)

* **Tracing**: Instrument HTTP endpoints via `go.opentelemetry.io/contrib/instrumentation/net/http/otelchi`.
* **Database Tracing**: Wrap `pgx/v5` PostgreSQL connections via `otelsql` / `otelpgx` to capture query duration, DDL/DML statements, and schema scope without exposing sensitive data values.
* **Structured Logging**: Bridge Go's native `log/slog` to OpenTelemetry using `go.opentelemetry.io/contrib/bridges/otelslog`. Automatically inject `trace_id` and `span_id` into every log entry.
* **Metrics**: Collect Go runtime metrics (memory, goroutines) and HTTP request histograms via `go.opentelemetry.io/otel/sdk/metric`.

### 2. Frontend Layer (React + TypeScript)

* **Browser Tracing**: Instrument Vite/React SPA using `@opentelemetry/sdk-trace-web` and `@opentelemetry/instrumentation-fetch`.
* **W3C Trace Context Propagation**: Inject `traceparent` and `tracestate` HTTP headers into API calls (`/api/v1/*`) to stitch browser user actions directly to Go backend traces.
* **Real User Monitoring (RUM)**: Capture Core Web Vitals (LCP, FID, CLS) and unhandled React error boundaries as telemetry events.

### 3. Database & Cache Layer (PostgreSQL & Redis)

* **DB Span Generation**: Database calls generate child spans within the parent HTTP handler span.
* **Database Collector**: Deploy `otel-collector` with `postgresql` receiver to query `pg_stat_statements` for long-term query performance tracking.
* **Redis Tracing**: Instrument Redis client operations via `github.com/redis/go-redis/extra/redisotel/v9`.

### 4. Telemetry Pipeline Architecture

* All spans, metrics, and logs are pushed via OTLP (OpenTelemetry Protocol over gRPC/HTTP) to an **OpenTelemetry Collector** container (`otel-collector`).
* The collector routes telemetry data to selected backends (e.g., Grafana Tempo for traces, Prometheus for metrics, Loki / Elasticsearch for logs).

## Consequences

### Positive Impacts

* **End-to-End Traceability**: Ability to trace a single click on the React UI down to the exact SQL query executed in `core` or `tournament` schemas.
* **Zero Vendor Lock-in**: Standardized OTLP telemetry allows swapping backend storage engines (Grafana, Jaeger, Datadog) without code changes.
* **Schema Isolation Compliance**: Span attributes explicitly record the target PostgreSQL schema (`core`, `tournament`, `scheduling`) to audit ADR-003 boundary compliance.

### Trade-offs / Constraints

* **Slight CPU/Network Overhead**: Tracing introduces minor latency and memory overhead. Mitigated by applying head-based probabilistic sampling (e.g., 100% sampling in dev, 10% in production).
* **Bundle Size (Frontend)**: OpenTelemetry web packages add ~30KB to the React frontend bundle.
