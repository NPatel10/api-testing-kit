import { expect, test } from '@playwright/test';

test.describe('public site smoke', () => {
	test('landing page surfaces the product promise', async ({ page }) => {
		await page.goto('/');

		await expect(page.getByRole('heading', { name: /a premium API testing workspace that is safe to put in front of people/i })).toBeVisible();
		await expect(page.locator('main').getByRole('link', { name: 'Open app' })).toBeVisible();
	});

	test('docs page keeps the reference content visible', async ({ page }) => {
		await page.goto('/docs');

		await expect(page.getByRole('heading', { name: /reference for the shared request workspace/i })).toBeVisible();
		await expect(page.locator('#overview').getByRole('link', { name: 'Templates' })).toBeVisible();
	});

	test('features page keeps the product structure obvious', async ({ page }) => {
		await page.goto('/features');

		await expect(page.getByRole('heading', { name: /feature reference for the shared api workspace/i })).toBeVisible();
		await expect(page.getByRole('heading', { name: /collections and history/i })).toBeVisible();
	});

	test('templates page stays launchable from the public site', async ({ page }) => {
		await page.goto('/templates');

		await expect(page.getByRole('heading', { name: /templates with filters, previews, and launch paths/i })).toBeVisible();
		await expect(page.getByText('Selection', { exact: true })).toBeVisible();
		await expect(page.getByText('Filter templates', { exact: true })).toBeVisible();
		await expect(page.getByRole('link', { name: /open in \/app/i }).first()).toBeVisible();
	});

	test('case study page keeps the engineering narrative visible', async ({ page }) => {
		await page.goto('/case-study');

		await expect(page.getByRole('heading', { name: /architecture overview for the shared request workspace/i })).toBeVisible();
		await expect(page.getByText('Request pipeline', { exact: true })).toBeVisible();
	});

	test('app workspace exposes the request and response panels', async ({ page }) => {
		await page.goto('/app');

		await expect(page.getByText('Request builder', { exact: true })).toBeVisible();
		await expect(page.getByText('Response viewer', { exact: true })).toBeVisible();
		await expect(page.getByText('Recent guest runs', { exact: true })).toBeVisible();
	});

	test('app history keeps the guest preview state visible', async ({ page }) => {
		await page.goto('/app/history');

		await expect(page.getByRole('heading', { name: /request history preview/i })).toBeVisible();
		await expect(page.getByText('Preview timeline', { exact: true })).toBeVisible();
		await expect(page.locator('main').getByText('Guest preview', { exact: true }).first()).toBeVisible();
	});

	test('collection detail shows the locked guest state', async ({ page }) => {
		await page.goto('/app/collections/saved-workflows');

		await expect(page.locator('main').getByText('Saved workflows', { exact: true }).last()).toBeVisible();
		await expect(page.getByText('Request groups', { exact: true })).toBeVisible();
		await expect(page.getByText('Collection metadata', { exact: true })).toBeVisible();
		await expect(page.getByText('Gated in guest mode', { exact: true }).first()).toBeVisible();
	});
});
