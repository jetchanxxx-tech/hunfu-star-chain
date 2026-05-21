import { test, expect } from '@playwright/test';

test.describe('API Demo Data', () => {
  test('GET /api/v1/demo/home returns valid data', async ({ request }) => {
    const res = await request.get('/api/v1/demo/home');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('nickname');
    expect(body).toHaveProperty('week');
    expect(typeof body.week).toBe('number');
    expect(body).toHaveProperty('events');
    expect(Array.isArray(body.events)).toBe(true);
    expect(body).toHaveProperty('reports');
    expect(Array.isArray(body.reports)).toBe(true);
    expect(body).toHaveProperty('stats');
  });

  test('demo events have expected fields', async ({ request }) => {
    const res = await request.get('/api/v1/demo/home');
    const body = await res.json();
    if (body.events.length > 0) {
      const event = body.events[0];
      expect(event).toHaveProperty('event_type');
      expect(event).toHaveProperty('event_date');
    }
  });

  test('demo reports have expected fields', async ({ request }) => {
    const res = await request.get('/api/v1/demo/home');
    const body = await res.json();
    if (body.reports.length > 0) {
      const report = body.reports[0];
      expect(report).toHaveProperty('report_type');
      expect(report).toHaveProperty('report_date');
    }
  });

  test('demo works without any auth', async ({ request }) => {
    const res = await request.get('/api/v1/demo/home', {
      headers: { Authorization: '' },
    });
    expect(res.status()).toBe(200);
  });
});
