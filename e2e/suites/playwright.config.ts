import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
    testDir: './offergen-tests',
    timeout: 30 * 1000,
    fullyParallel: true,
    forbidOnly: !!process.env.CI,
    retries: 0,
    workers: !!process.env.CI ? 1 : 4,
    reporter: [['html', { host: '0.0.0.0' }]],
    use: {
        baseURL: 'http://offergen',
        trace: 'on',
        extraHTTPHeaders: { 'Authorization': `Basic YWRtaW46a2lza3V0eWE=` },
    },
    projects: [
        { name: 'setup', testMatch: /.*\.setup\.ts/ },
        {
            name: 'chromium',
            use: {
                ...devices['Desktop Chrome'],
            },
            dependencies: ['setup'],
        },
        {
            name: 'firefox',
            use: {
                ...devices['Desktop Firefox'],
            },
            dependencies: ['setup'],
        },
        {
            name: 'webkit',
            use: {
                ...devices['Desktop Safari'],
            },
            dependencies: ['setup'],
        },
    ]
});
