import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import SonarReporter from "vitest-sonar-reporter";

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: "jsdom",
    setupFiles: "./src/test/setup.ts",
    reporters: [
      "default",
      new SonarReporter({
        outputFile: "test-report-frontend.xml",
      }),
    ],
    coverage: {
      provider: "v8",
      reporter: ["text", "lcov", "clover"],
      reportsDirectory: "./coverage",
      exclude: [
        "node_modules/",
        "src/test/",
        "**/*.d.ts",
        "**/*.config.*",
        "**/snapshot*",
      ],
    },
  },
});
