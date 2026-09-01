import { WebTracerProvider } from "@opentelemetry/sdk-trace-web";
import { SimpleSpanProcessor } from "@opentelemetry/sdk-trace-web";
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http";
import { FetchInstrumentation } from "@opentelemetry/instrumentation-fetch";
import { registerInstrumentations } from "@opentelemetry/instrumentation";

export function initFrontendObservability() {
  const exporter = new OTLPTraceExporter({
    url:
      import.meta.env.VITE_OTEL_COLLECTOR_URL ||
      "http://localhost:4318/v1/traces",
  });

  const provider = new WebTracerProvider({
    spanProcessors: [new SimpleSpanProcessor(exporter)],
  });

  provider.register();

  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: [/localhost:8080/, /\/api\/v1\/.*/],
      }),
    ],
  });

  console.log("[PULSE] OpenTelemetry Web Tracing initialized.");
}
