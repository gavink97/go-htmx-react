import { test, expect } from '@playwright/test';

test('api', async ({ request }) => {
    const base = 'http://localhost:8080'
    const data = await request.get(`${base}/api/data`);
    
    expect(data.ok()).toBeTruthy();
    expect(await data.json()).toMatchObject(expect.objectContaining({
        body: "I love avira so much",
        status: "Successful",
    }));
});
