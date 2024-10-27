import { test } from './authentication';
import { authFile } from './inventory';

test('sets up user with many items', async ({ browser, authenticatedHome }) => {
    let createItem = async () => {
        const c = await browser.newContext({ storageState: await authenticatedHome.context().storageState() })
        const p = await c.newPage();
        await p.goto("/");
        await p.getByText('Inventory', { exact: true }).click()
        for (let i = 0; i < 21; i++) {
            await p.getByText('Create new', { exact: true }).click()
            await p.getByPlaceholder('item name...', { exact: true }).fill('example item');
            await p.getByPlaceholder('price...', { exact: true }).fill('5500');
            await p.getByText('Create', { exact: true }).click();
        }
    }

    const tasks: Promise<void>[] = [];
    for (let i = 0; i < 8; i++) {
        tasks.push(createItem());
    }
    await Promise.all(tasks);

    await authenticatedHome.context().storageState({ path: authFile });
});
