import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
    await page.goto('/');

    await expect(page).toHaveTitle('Offergen: inventory manager and offering generator');
});

test('can generate menu', async ({ page }) => {
    await page.goto('/');

    await page.getByPlaceholder("item name...", { exact: true }).fill("example item 1");
    await page.getByPlaceholder("price...", { exact: true }).fill("1900");
    await page.getByText("Add", { exact: true }).click();
    await page.getByText("example item 1", { exact: true }).isVisible();
    await page.getByText("1900 USD", { exact: true }).isVisible();

    await page.getByPlaceholder("item name...", { exact: true }).fill("example item 2");
    await page.getByPlaceholder("price...", { exact: true }).fill("2500");
    await page.getByText("Add", { exact: true }).click();
    await page.getByText("example item 2", { exact: true }).isVisible();
    await page.getByText("2400 USD", { exact: true }).isVisible();

    await page.getByPlaceholder("item name...", { exact: true }).fill("example item 3");
    await page.getByPlaceholder("price...", { exact: true }).fill("3400");
    await page.getByText("Add", { exact: true }).click();
    await page.getByText("example item 3", { exact: true }).isVisible();
    await page.getByText("3400 USD", { exact: true }).isVisible();

    page.on("dialog", dialog => dialog.accept());

    await page
        .locator('#itemsContainer div')
        .filter({ hasText: 'example item 33400 USD' })
        .locator('#itemRowDeleteButton').click();

    await page.getByText("Generate offering", { exact: true }).click();
    await page.getByText("Offering", { exact: true }).isVisible();
    await page.getByText("example item 3", { exact: true }).isVisible();
    await page.getByText("3400 USD", { exact: true }).isVisible();
    await page.getByText("example item 1", { exact: true }).isVisible();
    await page.getByText("1900 USD", { exact: true }).isVisible();
});
