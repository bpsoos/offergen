import { randomEmail, signUp } from '../authentication';
import { test } from '@playwright/test'

export const signInOrSignUpText = 'Sign in or sign up';


test('tries to sign up with invalid passcode', async ({ page }) => {
    await page.goto('/');

    await page.getByText(signInOrSignUpText, { exact: true }).click();

    const emailInput = page.getByPlaceholder('john.doe@example.com', { exact: true });
    await emailInput.fill(randomEmail());
    await page.getByText('Continue', { exact: true }).click();
    await page.getByText('Confirm', { exact: true }).click();

    const passcodeInput = page.getByPlaceholder('123456');
    await passcodeInput.fill('123123');
    await page.getByText('Submit', { exact: true }).click();
    await page.getByText('Invalid passcode', { exact: true }).isVisible();
});

test('tries to sign up with invalid email address', async ({ page }) => {
    await page.goto('/');

    await page.getByText(signInOrSignUpText, { exact: true }).click();

    const emailInput = page.getByPlaceholder('john.doe@example.com', { exact: true });
    await emailInput.fill(randomEmail());
    await page.getByText('email: validation error', { exact: true }).isVisible();
});

test('signs up and deletes user', async ({ page, request }) => {
    const email = randomEmail();
    await signUp(page, request, email);

    await page.getByText('Profile', { exact: true }).click();
    await page.getByText(email).isVisible();

    page.on("dialog", dialog => dialog.accept())
    await page.getByText('Delete account', { exact: true }).click();
    await page.getByText(signInOrSignUpText).click();

    const emailInput = page.getByPlaceholder('john.doe@example.com', { exact: true });
    await emailInput.fill(email);
    await page.getByText('Confirm', { exact: true }).isVisible();
});

