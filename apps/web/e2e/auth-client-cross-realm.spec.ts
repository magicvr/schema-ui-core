import { expect, test, type Request as PlaywrightRequest } from "@playwright/test";

interface CapturedRequest {
  url: string;
  method: string;
  headers: Record<string, string>;
}

interface CrossRealmResult {
  constructorsAreForeign: boolean;
  sameOriginStatus: number;
  crossOriginRequestStatus: number;
  crossOriginUrlStatus: number;
  authEndpointStatus: number;
}

const TEST_ACCESS_TOKEN = "cross-realm-access";
const TEST_REFRESH_TOKEN = "cross-realm-refresh";
const CROSS_ORIGIN = "https://cross-realm.example.test";

interface RealmWindow extends Window {
  Request: typeof Request;
  URL: typeof URL;
}

function captureRequest(request: PlaywrightRequest): CapturedRequest {
  return {
    url: request.url(),
    method: request.method(),
    headers: request.headers(),
  };
}

function header(request: CapturedRequest, name: string): string | undefined {
  return request.headers[name.toLowerCase()];
}

test("authFetch preserves same-origin credentials across a real iframe realm", async ({ page }) => {
  await page.goto("/?e2e=auth-client-cross-realm");

  const captured: CapturedRequest[] = [];
  let captureStarted = false;

  await page.exposeFunction("__markCrossRealmAuthTestStarted", () => {
    captured.length = 0;
    captureStarted = true;
  });

  page.on("request", (request) => {
    if (!captureStarted) {
      return;
    }
    const url = new URL(request.url());
    if (
      url.pathname === "/api/account/sessions" ||
      url.pathname === "/api/auth/login" ||
      url.pathname === "/api/auth/refresh"
    ) {
      captured.push(captureRequest(request));
    }
  });

  await page.route(/\/api\/account\/sessions(?:\?.*)?$/, async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      headers: {
        "access-control-allow-origin": "*",
        "access-control-allow-headers": "*",
        "access-control-allow-methods": "GET, OPTIONS",
      },
      body: JSON.stringify({ items: [] }),
    });
  });

  await page.route(/\/api\/auth\/login(?:\?.*)?$/, async (route) => {
    await route.fulfill({
      status: 401,
      contentType: "application/json",
      body: JSON.stringify({ error: "UNAUTHORIZED" }),
    });
  });

  // A refresh response is provided as a safety net: if the auth endpoint is
  // misclassified, the test still completes and can assert that this route
  // was unexpectedly requested.
  await page.route(/\/api\/auth\/refresh(?:\?.*)?$/, async (route) => {
    await route.fulfill({
      status: 200,
      contentType: "application/json",
      body: JSON.stringify({
        accessToken: "unexpected-refresh-access",
        refreshToken: "unexpected-refresh-token",
      }),
    });
  });

  const result = await page.evaluate(
    async ({ accessToken, refreshToken, crossOrigin }) => {
      const mainWindow = window as unknown as RealmWindow;
      const markStarted = (window as unknown as Window & {
        __markCrossRealmAuthTestStarted: () => Promise<void>;
      }).__markCrossRealmAuthTestStarted;

      const iframe = document.createElement("iframe");
      iframe.src = "about:blank";
      document.body.appendChild(iframe);

      try {
        const foreignWindow = iframe.contentWindow as unknown as RealmWindow | null;
        if (foreignWindow === null) {
          throw new Error("cross-realm iframe window is unavailable");
        }
        if (foreignWindow.Request === mainWindow.Request || foreignWindow.URL === mainWindow.URL) {
          throw new Error("iframe did not provide distinct Request and URL constructors");
        }

        const [{ authFetch }, tokens] = await Promise.all([
          // @ts-expect-error Vite resolves this absolute source-module URL at runtime.
          import("/src/account/auth-client.ts"),
          // @ts-expect-error Vite resolves this absolute source-module URL at runtime.
          import("/src/account/tokens.ts"),
        ]);
        tokens.setAccessToken(accessToken);
        tokens.setRefreshToken(refreshToken);

        await markStarted();

        const sameOriginRequest = new foreignWindow.Request(
          `${window.location.origin}/api/account/sessions`,
        );
        const crossOriginRequest = new foreignWindow.Request(
          `${crossOrigin}/api/account/sessions`,
        );
        const crossOriginUrl = new foreignWindow.URL(
          `${crossOrigin}/api/account/sessions`,
        );
        const authEndpointRequest = new foreignWindow.Request(
          `${window.location.origin}/api/auth/login`,
        );

        const sameOriginResponse = await authFetch(sameOriginRequest);
        const crossOriginRequestResponse = await authFetch(crossOriginRequest);
        const crossOriginUrlResponse = await authFetch(crossOriginUrl);
        const authEndpointResponse = await authFetch(authEndpointRequest);

        return {
          constructorsAreForeign:
            foreignWindow.Request !== mainWindow.Request && foreignWindow.URL !== mainWindow.URL,
          sameOriginStatus: sameOriginResponse.status,
          crossOriginRequestStatus: crossOriginRequestResponse.status,
          crossOriginUrlStatus: crossOriginUrlResponse.status,
          authEndpointStatus: authEndpointResponse.status,
        } satisfies CrossRealmResult;
      } finally {
        iframe.remove();
      }
    },
    {
      accessToken: TEST_ACCESS_TOKEN,
      refreshToken: TEST_REFRESH_TOKEN,
      crossOrigin: CROSS_ORIGIN,
    },
  );

  expect(result).toEqual({
    constructorsAreForeign: true,
    sameOriginStatus: 200,
    crossOriginRequestStatus: 200,
    crossOriginUrlStatus: 200,
    authEndpointStatus: 401,
  });

  const sessionRequests = captured.filter(
    (request) => new URL(request.url).pathname === "/api/account/sessions",
  );
  expect(sessionRequests).toHaveLength(3);
  expect(new URL(sessionRequests[0].url).origin).toBe(new URL(page.url()).origin);
  expect(header(sessionRequests[0], "authorization")).toBe(`Bearer ${TEST_ACCESS_TOKEN}`);
  expect(header(sessionRequests[0], "x-refresh-token")).toBe(TEST_REFRESH_TOKEN);

  expect(new URL(sessionRequests[1].url).origin).toBe(CROSS_ORIGIN);
  expect(header(sessionRequests[1], "authorization")).toBeUndefined();
  expect(header(sessionRequests[1], "x-refresh-token")).toBeUndefined();

  expect(new URL(sessionRequests[2].url).origin).toBe(CROSS_ORIGIN);
  expect(header(sessionRequests[2], "authorization")).toBeUndefined();
  expect(header(sessionRequests[2], "x-refresh-token")).toBeUndefined();

  const authEndpointRequests = captured.filter(
    (request) => new URL(request.url).pathname === "/api/auth/login",
  );
  expect(authEndpointRequests).toHaveLength(1);
  expect(header(authEndpointRequests[0], "authorization")).toBe(`Bearer ${TEST_ACCESS_TOKEN}`);

  const refreshRequests = captured.filter(
    (request) => new URL(request.url).pathname === "/api/auth/refresh",
  );
  expect(refreshRequests, "a same-origin auth endpoint 401 must not trigger refresh").toHaveLength(0);
});
