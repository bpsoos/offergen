import { expect } from '@playwright/test';
import { test } from './authentication';
import { authFile } from './inventory';

test('sets up user with many items', async ({ request, authenticatedHome }) => {
    let createItem = async () => {
        let state = await authenticatedHome.context().storageState();
        const authCookie = state.cookies.find(cookie => cookie.name == "offergen__auth")
        expect(authCookie).toBeDefined();

        for (let i = 0; i < 21; i++) {
            const res = await request.post(
                'http://offergen/inventory/item',
                {
                    data: '{"name":"example item", "price": 5500}',
                    headers: {
                        'Cookie': 'offergen__auth=' + authCookie!.value,
                        'Content-Type': 'application/json',
                        'Accept': 'application/json',
                    }
                }
            )
            expect(res.ok()).toBeTruthy();

            const body = await res.json()
            expect(body.name).toEqual("example item")
            expect(body.price).toEqual(5500)
            expect(body.id).toBeDefined()
        }
    }

    const tasks: Promise<void>[] = [];
    for (let i = 0; i < 8; i++) {
        tasks.push(createItem());
    }
    await Promise.all(tasks);

    await authenticatedHome.context().storageState({ path: authFile });
});
