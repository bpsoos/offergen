import { expect, test as base, Page, APIRequestContext, } from "@playwright/test";

type Home = {
    authenticatedHome: Page;
};

export const test = base.extend<Home>({
    authenticatedHome: async ({ page, request }, use) => {
        const emailAddress = randomEmail();

        await signUp(page, request, emailAddress);

        await page.getByText('Log out', { exact: true }).click();
        await expect(page.getByText(signInOrSignUpText, { exact: true })).toHaveCount(1);

        await signIn(page, request, emailAddress);

        await expect(page.getByText('Log out', { exact: true })).toHaveCount(1);
        await use(page);
    }
});

export const signInOrSignUpText = 'Sign in or sign up';

export async function signUp(page: Page, request: APIRequestContext, email: string) {
    await page.goto('/');

    await page.getByText(signInOrSignUpText, { exact: true }).click();

    let emailInput = page.getByPlaceholder('john.doe@example.com', { exact: true });
    let emailAddress = email;
    await emailInput.fill(emailAddress);
    await page.getByText('Continue', { exact: true }).click();

    await page.getByText('Confirm', { exact: true }).click();

    const passcodeInputForSignUp = page.getByPlaceholder('123456');
    await passcodeInputForSignUp.isVisible();
    let passcodeForSignUp: string | null = '';
    await expect(async () => {
        const resp = await getEmails(request);
        passcodeForSignUp = getPasscode(resp.mailItems, emailAddress, 'verify')
        expect(passcodeForSignUp).not.toBeNull();
    }).toPass(
        {
            timeout: 2_000,
            intervals: [500, 500, 1_000, 1_000, 2_000],
        }
    );
    await passcodeInputForSignUp.fill(passcodeForSignUp);
    await page.getByText('Submit', { exact: true }).click();

    await page.getByText('Log out', { exact: true }).isVisible();
}

export async function signIn(page: Page, request: APIRequestContext, email: string) {
    await page.goto('/');

    await page.getByText(signInOrSignUpText, { exact: true }).click();

    let emailInput = page.getByPlaceholder('john.doe@example.com', { exact: true });
    let emailAddress = email;
    await emailInput.fill(emailAddress);
    await page.getByText('Continue', { exact: true }).click();

    const passcodeInputForSignUp = page.getByPlaceholder('123456');
    await passcodeInputForSignUp.isVisible();
    let passcodeForSignUp: string | null = '';
    await expect(async () => {
        const resp = await getEmails(request);
        passcodeForSignUp = getPasscode(resp.mailItems, emailAddress, 'login')
        expect(passcodeForSignUp).not.toBeNull();
    }).toPass(
        {
            timeout: 2_000,
            intervals: [500, 500, 1_000],
        }
    );
    await passcodeInputForSignUp.fill(passcodeForSignUp);
    await page.getByText('Submit', { exact: true }).click();

    await page.getByText('Log out', { exact: true }).isVisible();
}


export function randomEmail(): string {
    return Math.random().toString().substring(2) + '@email.com';
}

export function getPasscode(data: Array<email>, address: string, filter: string): string | null {
    const filterRe = new RegExp(filter);
    const emails = data.filter(
        (email) => {
            if (email.toAddresses[0] != address) {
                return false
            }

            return filterRe.test(email.subject)
        }
    )
    emails.sort((a, b) => {
        return new Date(a.dateSent).valueOf() - new Date(b.dateSent).valueOf();
    })
    const email = emails[0];

    if (email == null) {
        return null
    }
    const subject = email.subject
    if (subject == null) {
        return null
    }

    const re = /passcode (\d+)/;
    const execResults = re.exec(subject)!;

    return execResults[1];
}

export async function getEmails(request: APIRequestContext): Promise<emailResponse> {
    const getMailsResp = await request.get('http://smtp:8085/mail');
    return getMailsResp.json();
}

export interface email {
    toAddresses: Array<string>,
    subject: string,
    dateSent: string,
}

export interface emailResponse {
    mailItems: Array<email>,
}
