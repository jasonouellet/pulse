import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router";
import App from "./App.tsx";
import "./index.css";
import { initFrontendObservability } from "./lib/observability.ts";
import './i18n/config'; // <-- Initialise i18next at startup

// Initialise OTEL tracing if activated
if (import.meta.env.VITE_ENABLE_OTEL === "true") {
  initFrontendObservability();
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </React.StrictMode>,
);
