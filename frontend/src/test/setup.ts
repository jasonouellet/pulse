import "@testing-library/jest-dom/vitest";
import "vitest-axe/extend-expect";
import { expect } from "vitest";
import * as matchers from "vitest-axe/matchers";

// Enregistre les matchers d'accessibilité (axe) dans Vitest
expect.extend(matchers);
