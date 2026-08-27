import "@testing-library/jest-dom/vitest";
import { expect } from "vitest";
import * as matchers from "vitest-axe/matchers";

// Enregistre les matchers d'accessibilité (axe) dans Vitest
expect.extend(matchers);
