import { expect, test } from '@playwright/test';

test('navigation flows between the public site and the workspace', async ({ page }) => {
	await page.goto('/');
	await page.getByRole('link', { name: 'Open the app' }).click();

	await expect(page.getByText('Request builder', { exact: true })).toBeVisible();
	await expect(page.getByRole('link', { name: 'Templates' })).toBeVisible();

	await page.getByRole('link', { name: 'Templates' }).click();
	await expect(page.getByRole('heading', { name: /safe api templates with category filters, search, and a real launch path/i })).toBeVisible();
	await expect(page.getByText('Launch snapshot', { exact: true })).toBeVisible();

	await page.getByRole('link', { name: /open \/app/i }).first().click();
	await expect(page.getByText('Response viewer', { exact: true })).toBeVisible();

	await page.getByRole('link', { name: 'History', exact: true }).click();
	await expect(page.getByRole('heading', { name: /see the history surface, but keep persistence behind sign-in/i })).toBeVisible();
});
