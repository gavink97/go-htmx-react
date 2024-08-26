import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

// for more information on accessibility testing
// https://playwright.dev/docs/accessibility-testing

// name tests
test.describe('homepage', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:8080/');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations).toEqual([]);
    });
});

test.describe('about page', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/about');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations).toEqual([]);
    });
});

test.describe('register page', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/register');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations).toEqual([]);
    });
});

test.describe('register account', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/register');
        await page.getByPlaceholder('name@company.com').click();
        await page.getByPlaceholder('name@company.com').fill('gavin@gmail.com');
        await page.getByPlaceholder('••••••••').click();
        await page.getByPlaceholder('••••••••').fill('password');
        await page.getByPlaceholder('••••••••').press('Enter');
    });
});

test.describe('login page', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/login');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations).toEqual([]);
    });
});

test.describe('login', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/login');
        await page.getByPlaceholder('name@company.com').click();
        await page.getByPlaceholder('name@company.com').fill('gavin@gmail.com');
        await page.getByPlaceholder('••••••••').click();
        await page.getByPlaceholder('••••••••').fill('password');
        await page.getByPlaceholder('••••••••').press('Enter');
    });
});

test.describe('logout', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/login');
        await page.getByPlaceholder('name@company.com').click();
        await page.getByPlaceholder('name@company.com').fill('gavin@gmail.com');
        await page.getByPlaceholder('••••••••').click();
        await page.getByPlaceholder('••••••••').fill('password');
        await page.getByPlaceholder('••••••••').press('Enter');
        await page.getByRole('button', { name: 'Logout' }).click();
    });
});

test.describe('404 page', () => {
    test('test', async ({ page }) => {
        await page.goto('http://localhost:3000/404');
        const accessibilityScanResults = await new AxeBuilder({ page }).analyze();
        expect(accessibilityScanResults.violations).toEqual([]);
    });
});
