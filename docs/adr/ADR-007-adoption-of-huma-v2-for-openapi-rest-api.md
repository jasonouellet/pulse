# Adoption of Huma v2 for OpenAPI v3 REST API, Request Validation, and Documentation

* Status: accepted
* Deciders: PULSE Core Engineering Team
* Date: 2026-08-28

Technical Story: [PULSE-ADR-007] Transition from raw Chi with swaggo annotations to a type-safe, schema-driven REST framework.

## Context and Problem Statement

Project PULSE requires a robust, performant, and self-documenting REST API layer across all domain modules (`core`, `tournament`, `scheduling`, `evaluation`). In earlier iterations, the API was routed using raw Chi controllers with manual JSON marshaling and comment-based Swagger documentation generated via `swaggo/swag`.

This approach introduced significant operational friction:

* Documentation drift between Go implementation structs and docstring comments.
* Repetitive boilerplate code for HTTP request body parsing, query parameter extraction, and validation.
* Lack of runtime JSON Schema enforcement on incoming requests.

How can we achieve type-safe HTTP routing, automated runtime request validation, and guaranteed OpenAPI v3 specification generation without coupling our core domain logic to transport details?

## Decision Drivers

* **Type Safety & DX:** Need for strongly-typed HTTP inputs/outputs tied directly to Go structs.
* **Automated OpenAPI Specs:** Elimination of manual docstring annotations (`// @Param`, `// @Success`) in favor of code-as-contract.
* **Architecture Hexagonal Alignment:** Transport adapters (`internal/<module>/adapters/http`) must remain decoupled from core domain logic (`internal/<module>/domain`).
* **Performance & Middleware Support:** Native compatibility with Chi router, Prometheus metrics, and OpenTelemetry tracing.
* **Runtime Validation:** Built-in JSON Schema validation for path parameters, query strings, and JSON payloads.

## Considered Options

* **Option 1:** Huma v2 framework built on top of Chi (`github.com/danielgtaylor/huma/v2`)
* **Option 2:** Raw Chi router + `swaggo/swag` docstring generation
* **Option 3:** `oapi-codegen` (OpenAPI-first schema generation to Go boilerplate)

## Decision Outcome

Chosen option: **Option 1: Huma v2 framework built on top of Chi**, because it provides the optimal balance of code-first developer velocity, reflection-based OpenAPI v3 generation, and seamless integration with our existing Chi routing and OpenTelemetry instrumentation.

### Positive Consequences

* **Single Source of Truth:** Go request and response structs natively define both the runtime behavior and the OpenAPI v3 spec (`/openapi.json` and interactive `/docs` UI).
* **Zero Docstring Drift:** Complete elimination of fragile `swaggo/swag` annotations.
* **Automated Validation:** Path, query, and body parameters are automatically validated against struct tags (e.g., `maxLength`, `minimum`, `format:"uuid"`) before reaching handler logic.
* **Cleaner Handlers:** Reduced HTTP boilerplate code by over 40% across adapter packages.

### Negative Consequences

* **Framework Coupling:** Transport adapters depend directly on `huma/v2` primitives (`huma.Context`, `huma.Operation`).
* **Reflect Overhead:** Minimal reflection cost during route registration at application startup (negligible at runtime).

## Pros and Cons of the Options

### Option 1: Huma v2 with Chi (Selected)

Huma v2 acts as a modern API framework wrapper around standard Go HTTP routers (Chi in our case).

* Good, because request parsing, input validation, and output serialization are handled automatically via Go struct definitions.
* Good, because OpenAPI v3 JSON specs and Swagger UI / Stoplight Elements documentation are served dynamically at zero extra maintenance cost.
* Good, because it supports custom middleware and integrates seamlessly with our OpenTelemetry tracer and Prometheus metrics HTTP handlers.
* Neutral, requires learning Huma-specific struct tags (`doc`, `example`, `readOnly`) for enriched documentation.
* Bad, adds a external framework dependency to the HTTP adapter layer.

### Option 2: Raw Chi Router + Swaggo Annotations

The initial baseline where routes were mapped manually with Chi and documented via `swaggo/swag` CLI comments.

* Good, zero abstraction layer over standard `http.HandlerFunc`.
* Bad, comment annotations (`// @Summary`, `// @Failure`) frequently go out of sync with actual Go struct changes.
* Bad, requires manual `json.Decode()`, error handling, and validation logic in every single endpoint.

### Option 3: OpenAPI-First with `oapi-codegen`

Designing the OpenAPI YAML spec first and generating Go server boilerplate using code generators.

* Good, strict contract-first design methodology.
* Bad, high friction in rapid prototyping phases where domain models evolve frequently.
* Bad, generated code artifacts increase repository complexity and build setup overhead.

## Validation

The decision was validated during the `feature/core` module refactoring:

* Successfully migrated `/api/v1/core/users` endpoints to Huma v2 handlers.
* Verified automatic Swagger UI availability at `/docs`.
* Validated input constraint rejections (420/400 validation errors) without writing custom validation code.
* Unit test coverage confirmed with Huma test utilities (`huma-test`).

## Links

* [Huma v2 Official Documentation](https://huma.rocks/)
* [ADR-003: Modular Monolith Architecture](./0003-modular-monolith-architecture.md)
* [ADR-006: Global Observability and Telemetry](./0006-observability-opentelemetry.md)
