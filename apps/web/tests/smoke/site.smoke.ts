import { expect, test } from '@playwright/test';

test.describe('public site smoke', () => {
	test('landing page surfaces the product promise', async ({ page }) => {
		await page.goto('/');

		await expect(page.getByRole('heading', { name: /premium API testing workspace/i })).toBeVisible();
		await expect(page.getByRole('link', { name: 'Open the app' })).toBeVisible();
	});

	test('app workspace exposes the request and response panels', async ({ page }) => {
		await page.goto('/app');

		await expect(page.locator('main').getByText('Request builder', { exact: true }).first()).toBeVisible();
		await expect(page.locator('main').getByText('Response viewer', { exact: true }).first()).toBeVisible();
	});

	test('docs page keeps the quick start guidance visible', async ({ page }) => {
		await page.goto('/docs');

		await expect(page.getByRole('heading', { name: /learn the product in a few minutes/i })).toBeVisible();
		await expect(page.locator('#quick-start').getByRole('link', { name: 'Open /app' })).toBeVisible();
	});
});
