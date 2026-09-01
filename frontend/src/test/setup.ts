import "@testing-library/jest-dom/vitest";
import "vitest-axe/extend-expect";
import { expect, vi } from "vitest";
import * as matchers from "vitest-axe/matchers";

// Enregistre les matchers d'accessibilité (axe) dans Vitest
expect.extend(matchers);

// Mock de window.matchMedia pour JSDOM
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});
