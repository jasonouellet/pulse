import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import sonarReporter from 'vitest-sonar-reporter';

export default defineConfig({
    plugins: [react()],
    test: {
        globals: true,
        environment: 'jsdom',
        setupFiles: './src/test/setup.ts',
        reporters: [
            'default',
            [
                'vitest-sonar-reporter',
                {
                    outputFile: 'test-report-frontend.xml',
                },
            ],
        ],
        coverage: {
            provider: 'v8',
            reporter: ['text', 'lcov', 'clover'],
            reportsDirectory: './coverage',
            exclude: [
                'node_modules/',
                'src/test/',
                '**/*.d.ts',
                '**/*.config.*',
                '**/snapshot*',
            ],
        },
    },
});
