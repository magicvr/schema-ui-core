import { expect, test } from "@playwright/test";

// Host failure browser-level conformance (ADR-0036 D7 / spec 10 §3.8):
// production Host surface for maintenance / protocol-rejected / route
// not-found, focus management, live-region announcement classes and recovery
// actions. The bootstrap document is intercepted on the REAL well-known
// entry; the manifest must NOT be fetched after a terminal availability gate.

const BOOTSTRAP_PATH = "/.well-known/schema-ui/host-bootstrap.json";

function maintenanceDocument() {
  return {
    bootstrapVersion: "1.0",
    requiredCapabilities: ["host.bootstrap"],
    manifest: { url: "/.well-known/schema-ui/app-manifest.json" },
    availability: { mode: "maintenance", retryAfterSeconds: 60, messageKey: "hostFailure.maintenance" },
  };
}

function protocolRejectedDocument() {
  return {
    bootstrapVersion: "1.0",
    requiredCapabilities: ["app.manifest"],
    manifest: { url: "/.well-known/schema-ui/app-manifest.json" },
    availability: { mode: "normal" },
  };
}

test("maintenance terminal: renders before manifest fetch, polite announcement, focus on title", async ({ page }) => {
  let manifestFetched = false;
  await page.route("**/.well-known/schema-ui/app-manifest.json", async (route) => {
    manifestFetched = true;
    await route.abort();
  });
  await page.route(`**${BOOTSTRAP_PATH}`, async (route) => {
    await route.fulfill({ contentType: "application/json", body: JSON.stringify(maintenanceDocument()) });
  });

  await page.goto("/");

  const title = page.getByRole("heading", { name: "Scheduled maintenance" });
  await expect(title).toBeVisible();
  // Stage order: availability terminal before manifest fetch (ADR-0035 D3).
  expect(manifestFetched).toBe(false);
  // Focus moved to the unique error title on first terminal entry.
  await expect.poll(() => page.evaluate(() => document.activeElement?.id)).toBe("host-failure-title");
  // maintenance uses a polite live region (not assertive).
  const region = page.locator('[aria-live="polite"][role="status"]');
  await expect(region).toHaveText(/Scheduled maintenance/);
  // Retry rebuilds the whole application instance (reload semantics).
  await expect(page.getByRole("button", { name: "Retry" })).toBeVisible();
});

test("protocol-rejected terminal: assertive announcement, no continue-render action", async ({ page }) => {
  await page.route(`**${BOOTSTRAP_PATH}`, async (route) => {
    await route.fulfill({ contentType: "application/json", body: JSON.stringify(protocolRejectedDocument()) });
  });

  await page.goto("/");

  await expect(page.getByRole("heading", { name: "Incompatible application" })).toBeVisible();
  const region = page.locator('[aria-live="assertive"][role="status"]');
  await expect(region).toHaveText(/Incompatible application/);
  // protocol-rejected must not offer "continue rendering": no recovery
  // actions beyond the (empty) allowed set.
  await expect(page.getByRole("button", { name: "Retry" })).toHaveCount(0);
});

test("route not-found: HOST_ROUTE_NOT_FOUND surface and home recovery with focus", async ({ page }) => {
  await page.route(`**${BOOTSTRAP_PATH}`, async (route) => {
    await route.fulfill({ contentType: "application/json", body: JSON.stringify(maintenanceDocument()) });
  });

  // Authenticate via the dev seed to reach the shell router.
  await page.goto("/");
  await page.getByRole("heading", { name: "Scheduled maintenance" }).waitFor();
  // Maintenance blocks sign-in (correct stage order); unroute and reload to
  // exercise the router path with a real normal bootstrap.
  await page.unroute(`**${BOOTSTRAP_PATH}`);
  await page.goto("/");
  await page.getByLabel("Username").fill("admin");
  await page.getByLabel("Password").fill("admin");
  await page.getByRole("button", { name: "Sign in" }).click();
  await expect(page).toHaveURL(/\/users$/);

  await page.goto("/definitely-not-a-manifest-page");
  await expect(page.getByRole("heading", { name: "Page not found" })).toBeVisible();
  await expect.poll(() => page.evaluate(() => document.activeElement?.id)).toBe("host-failure-title");

  await page.getByRole("button", { name: "Return home" }).click();
  await expect(page).toHaveURL(/\/users$/);
  await expect(page.getByRole("heading", { name: "Users" })).toBeVisible();
  // Recovery lands focus on the restored page's main heading.
  await expect.poll(() => page.evaluate(() => document.activeElement?.textContent ?? "")).toBe("Users");
});

test("normal bootstrap document from the real entry boots the shell", async ({ page }) => {
  // The Go API serves the real document (same-origin, public entry).
  const bootstrap = await page.request.get(BOOTSTRAP_PATH);
  expect(bootstrap.status()).toBe(200);
  const document = await bootstrap.json();
  expect(document.bootstrapVersion).toBe("1.0");
  expect(document.manifest.url).toBe("/.well-known/schema-ui/app-manifest.json");
  expect(document.manifest.sha256).toMatch(/^[0-9a-f]{64}$/);

  await page.goto("/");
  await expect(page.getByRole("heading", { name: "Sign in" })).toBeVisible();
});
