import { expect } from '@playwright/test';
import { test } from './authentication';
import { filledInventoryTest } from './inventory';

test('can add and delete items from inventory', async ({ authenticatedHome }) => {
    let page = authenticatedHome;
    await page.goto("/");
    await page.getByText('Inventory', { exact: true }).click()
    await page.getByText('Create new', { exact: true }).click()
    await page.getByPlaceholder('item name...', { exact: true }).fill('example item 1');
    await page.getByPlaceholder('price...', { exact: true }).fill('1900');
    await page.getByText('Create', { exact: true }).click();
    await expect(page.getByText('example item 1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('1900 USD', { exact: true })).toHaveCount(1);

    await page.getByText('Create new', { exact: true }).click()
    await page.getByPlaceholder('item name...', { exact: true }).fill('example item 2');
    await page.getByPlaceholder('price...', { exact: true }).fill('2500');
    await page.getByText('Create', { exact: true }).click();
    await expect(page.getByText('example item 2', { exact: true })).toHaveCount(1);
    await expect(page.getByText('2500 USD', { exact: true })).toHaveCount(1);

    await page.getByText('Create new', { exact: true }).click()
    await page.getByPlaceholder('item name...', { exact: true }).fill('example item 3');
    await page.getByPlaceholder('price...', { exact: true }).fill('3400');
    await page.getByText('Create', { exact: true }).click();
    await expect(page.getByText('example item 3', { exact: true })).toHaveCount(1);
    await expect(page.getByText('3400 USD', { exact: true })).toHaveCount(1);

    page.on('dialog', dialog => dialog.accept());
    const deleteButtonLocator = page
        .locator('#inventoryItemsContainer .item')
        .filter({ hasText: 'example item 22500 USD' })
        .locator('.delete-button')

    await deleteButtonLocator.click();

    await expect(page.getByText('example item 3', { exact: true })).toHaveCount(1);
    await expect(page.getByText('3400 USD', { exact: true })).toHaveCount(1);
    await expect(page.getByText('example item 1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('1900 USD', { exact: true })).toHaveCount(1);
    await expect(page.getByText('example item 2', { exact: true })).toHaveCount(0);
    await expect(page.getByText('2400 USD', { exact: true })).toHaveCount(0);
});

filledInventoryTest('can paginate inventory', async ({ page }) => {
    await page.goto("/");
    await page.getByText('Inventory', { exact: true }).click()

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('2', { exact: true })).toHaveCount(1);
    await expect(page.getByText('3', { exact: true })).toHaveCount(1);
    await expect(page.getByText('4', { exact: true })).toHaveCount(1);
    await expect(page.getByText('5', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.getByText('6', { exact: true }).click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('4', { exact: true })).toHaveCount(1);
    await expect(page.getByText('5', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('8', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.getByText('17', { exact: true }).click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('11', { exact: true })).toHaveCount(1);
    await expect(page.getByText('12', { exact: true })).toHaveCount(1);
    await expect(page.getByText('13', { exact: true })).toHaveCount(1);
    await expect(page.getByText('14', { exact: true })).toHaveCount(1);
    await expect(page.getByText('15', { exact: true })).toHaveCount(1);
    await expect(page.getByText('16', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-backward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('8', { exact: true })).toHaveCount(1);
    await expect(page.getByText('9', { exact: true })).toHaveCount(1);
    await expect(page.getByText('10', { exact: true })).toHaveCount(1);
    await expect(page.getByText('11', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-forward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('11', { exact: true })).toHaveCount(1);
    await expect(page.getByText('12', { exact: true })).toHaveCount(1);
    await expect(page.getByText('13', { exact: true })).toHaveCount(1);
    await expect(page.getByText('14', { exact: true })).toHaveCount(1);
    await expect(page.getByText('15', { exact: true })).toHaveCount(1);
    await expect(page.getByText('16', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.getByText('1', { exact: true }).click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('2', { exact: true })).toHaveCount(1);
    await expect(page.getByText('3', { exact: true })).toHaveCount(1);
    await expect(page.getByText('4', { exact: true })).toHaveCount(1);
    await expect(page.getByText('5', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-forward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('8', { exact: true })).toHaveCount(1);
    await expect(page.getByText('9', { exact: true })).toHaveCount(1);
    await expect(page.getByText('10', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-forward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('10', { exact: true })).toHaveCount(1);
    await expect(page.getByText('11', { exact: true })).toHaveCount(1);
    await expect(page.getByText('12', { exact: true })).toHaveCount(1);
    await expect(page.getByText('13', { exact: true })).toHaveCount(1);
    await expect(page.getByText('14', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-backward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('8', { exact: true })).toHaveCount(1);
    await expect(page.getByText('9', { exact: true })).toHaveCount(1);
    await expect(page.getByText('10', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);

    await page.locator('.jump-backward').click();

    await expect(page.getByText('1', { exact: true })).toHaveCount(1);
    await expect(page.getByText('2', { exact: true })).toHaveCount(1);
    await expect(page.getByText('3', { exact: true })).toHaveCount(1);
    await expect(page.getByText('4', { exact: true })).toHaveCount(1);
    await expect(page.getByText('5', { exact: true })).toHaveCount(1);
    await expect(page.getByText('6', { exact: true })).toHaveCount(1);
    await expect(page.getByText('7', { exact: true })).toHaveCount(1);
    await expect(page.getByText('17', { exact: true })).toHaveCount(1);
});
