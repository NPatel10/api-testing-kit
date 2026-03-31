import { expect, test } from '@playwright/test';

test('navigation flows between the public site and the workspace', async ({ page }) => {
	await page.goto('/');
	await page.locator('main').getByRole('link', { name: 'Open app' }).click();

	await expect(page.getByText('Request builder', { exact: true })).toBeVisible();
	await expect(page.getByRole('link', { name: 'Templates' })).toBeVisible();

	await page.getByRole('link', { name: 'Templates' }).click();
	await expect(page.getByRole('heading', { name: /templates with filters, previews, and launch paths/i })).toBeVisible();
	await expect(page.getByText('Selection', { exact: true })).toBeVisible();

	await page.getByRole('link', { name: /open \/app/i }).first().click();
	await expect(page.getByText('Response viewer', { exact: true })).toBeVisible();

	await page.getByRole('link', { name: 'History', exact: true }).click();
	await expect(page.getByRole('heading', { name: /request history preview/i })).toBeVisible();
});
