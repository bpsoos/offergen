import { test } from './authentication';

export const authFile = './../.auth/user_with_many_items.json';

export const filledInventoryTest = test.extend({
    page: async ({ browser }, use) => {
        const context = await browser.newContext({ storageState: authFile });
        const page = await context.newPage();
        await page.goto("/");

        await use(page)

        await context.close();
    },
});

