import { expect, test, type Page } from "@playwright/test";

import { signInAsAdmin } from "./sign-in";

// W14 R-001 (GOAL-015 / workspace-009): browser-level production-wiring smoke.
//
// GOAL-013 F-010 moved GET /api/schema/{pageId} behind auth middleware, and the
// production entry initially shipped WITHOUT passing the authed transport as
// schemaFetcher — every page-document request went out anonymous (401) and
// EVERY page rendered PageSchemaErrorSurface ("无法显示此页面"), while all unit
// tests stayed green because they inject explicit fetchers. The vitest wiring
// lock (auth-gate.wiring.test.tsx) guards the assembly point with a stubbed
// fetch boundary; THIS spec closes R-001 by asserting the real network path in
// a real Chromium against the real API: after sign-in, every /api/schema/*
// request the SPA emits must carry a Bearer Authorization header and succeed.

interface CapturedSchemaRequest {
  url: string;
  authorization: string | null;
}

function collectSchemaTraffic(page: Page): {
  requests: CapturedSchemaRequest[];
  statuses: number[];
} {
  const requests: CapturedSchemaRequest[] = [];
  const statuses: number[] = [];
  page.on("request", (request) => {
    if (new URL(request.url()).pathname.startsWith("/api/schema/")) {
      // Playwright lower-cases header names.
      requests.push({ url: request.url(), authorization: request.headers()["authorization"] ?? null });
    }
  });
  page.on("response", (response) => {
    if (new URL(response.url()).pathname.startsWith("/api/schema/")) {
      statuses.push(response.status());
    }
  });
  return { requests, statuses };
}

test("signed-in page-schema requests ride the Bearer transport end-to-end (W14 R-001)", async ({ page }) => {
  test.setTimeout(90_000);

  await signInAsAdmin(page);

  // Attach AFTER sign-in so the captured traffic is purely the authenticated
  // shell's page-schema loads, then force a full app boot onto a manifest page
  // (fresh document → restoreSession → AuthGate → App → loadPageDocument).
  //
  // The SPA fires its schema fetches AFTER the navigation's load event (React
  // mount effects), so the collector alone races the assertions. Wait
  // deterministically for the FIRST schema response, then let the collector
  // settle via expect.poll before making "every capture" claims.
  const { requests, statuses } = collectSchemaTraffic(page);
  const firstSchemaResponse = page.waitForResponse(
    (response) => new URL(response.url()).pathname.startsWith("/api/schema/"),
  );
  await page.goto("/users");
  const first = await firstSchemaResponse;

  // Positive control on the deterministic capture itself…
  expect(first.status(), "first page-schema request failed").toBe(200);
  expect(
    first.request().headers()["authorization"] ?? "",
    "first page-schema request carried no Bearer token",
  ).toMatch(/^Bearer .+/);

  // …and let any remaining boot-time schema traffic settle.
  await expect
    .poll(() => statuses.length, { timeout: 15_000 })
    .toBeGreaterThanOrEqual(1);

  // THE lock: EVERY observed page-schema request carries a Bearer access
  // token, and the API accepts them all (an anonymous 401 would surface as a
  // non-200 here).
  expect(requests.length, "expected at least one /api/schema/* request after sign-in").toBeGreaterThan(0);
  for (const request of requests) {
    expect(
      request.authorization,
      `missing/invalid Authorization header on ${request.url}`,
    ).toMatch(/^Bearer .+/);
  }
  for (const status of statuses) {
    expect(status, "page-schema request failed").toBe(200);
  }

  // The page must render its real content, not the fail-closed error surface
  // (either locale of shell.pageSchemaError.title).
  await expect(
    page.getByRole("heading", { name: /can't be displayed|无法显示此页面/ }),
  ).toHaveCount(0);
});
