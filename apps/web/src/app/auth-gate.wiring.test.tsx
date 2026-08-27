// @vitest-environment jsdom

/**
 * W14 F-001 · production-wiring regression lock（「测试装配 ≠ 生产装配」）。
 *
 * History: GOAL-013 F-010 moved `/api/schema/{pageId}` behind auth middleware,
 * but the production entry never passed an authed schemaFetcher into `<App>` —
 * every page-document request went out anonymous (HTTP 401) and EVERY page
 * rendered 无法显示此页面 (PageSchemaErrorSurface). All 30+ existing tests
 * inject an explicit schemaFetcher, so the suite stayed green while the real
 * assembly was broken: the entry point's createRoot side effects made its
 * wiring invisible to tests.
 *
 * This file mounts the REAL production gate (`AuthGate`, extracted verbatim
 * from main.tsx in W14 F-001) inside the REAL `AuthProvider` with only the
 * transport boundary mocked, then asserts what the gate actually hands to
 * `<App>`. If the authed wiring is dropped or de-graded again, THIS fails —
 * not production.
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

const mocks = vi.hoisted(() => ({
  restoreSession: vi.fn(),
}));

// Partial mock: keep the REAL authFetch export so the identity assertion below
// proves the gate passes the genuine Bearer / refresh-on-401 transport.
vi.mock("@/account/auth-client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/account/auth-client")>();
  return {
    ...actual,
    restoreSession: mocks.restoreSession,
  };
});

const captured = vi.hoisted(() => ({
  appProps: null as Record<string, unknown> | null,
}));

// Capture the props the PRODUCTION gate passes to the shell <App>.
vi.mock("@/app/App", () => ({
  App: (props: Record<string, unknown>) => {
    captured.appProps = props;
    return null;
  },
}));

import type { AuthSession } from "@/account/auth-client";
import { authFetch } from "@/account/auth-client";
import { AuthProvider } from "@/account/AuthContext";
import { AuthGate } from "@/app/AuthGate";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const __dir = dirname(fileURLToPath(import.meta.url));
const MANIFEST_PATH = resolve(__dir, "../test-fixtures/app-manifest.admin-dogfood.json");

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
  mocks.restoreSession.mockReset();
  captured.appProps = null;
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
});

function adminManifest(): AppManifest {
  return validateAppManifest(JSON.parse(readFileSync(MANIFEST_PATH, "utf8")));
}

function sessionFor(): { kind: "session"; session: AuthSession } {
  return {
    kind: "session",
    session: {
      user: { id: "user-admin", name: "Admin", roles: ["admin"], permissions: ["dashboard.view"] },
      features: { menu_dashboard: true },
    },
  };
}

/** Mounts the production gate tree exactly as main.tsx assembles it. */
function renderProductionGate(): void {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider stored="zh-CN">
        <AuthProvider>
          <AuthGate manifest={adminManifest()} />
        </AuthProvider>
      </I18nProvider>,
    );
  });
}

async function flushAuth(): Promise<void> {
  await act(async () => {
    await Promise.resolve();
  });
}

describe("production gate → shell wiring (W14 F-001)", () => {
  it("hands the REAL authed transports to <App> for an authenticated session", async () => {
    mocks.restoreSession.mockResolvedValue(sessionFor());
    renderProductionGate();
    await flushAuth();

    expect(captured.appProps).not.toBeNull();
    const props = captured.appProps as Record<string, unknown>;

    // THE lock: page-schema documents (/api/schema/{pageId}, GOAL-013 F-010)
    // must ride the same Bearer / refresh-on-401 transport as every other API
    // call. A bare fetch here = anonymous 401 on every page again.
    expect(props.schemaFetcher).toBe(authFetch);
    expect(props.schemaFetcher).not.toBe(globalThis.fetch);

    // The resource path keeps its config-aware authed wrapper (GOAL-011):
    // still a function, still not the bare global fetch.
    expect(typeof props.resourceFetcher).toBe("function");
    expect(props.resourceFetcher).not.toBe(globalThis.fetch);

    // Session context reaches the shell expression engine intact.
    const context = props.navigationContext as {
      user?: { id?: string };
      features?: Record<string, boolean>;
    };
    expect(context?.user?.id).toBe("user-admin");
    expect(context?.features).toEqual({ menu_dashboard: true });
    expect(props.onLogout).toBeTypeOf("function");
    expect((props.currentUser as { id?: string })?.id).toBe("user-admin");
  });

  it("mounts no shell and renders the login surface when unauthenticated", async () => {
    mocks.restoreSession.mockResolvedValue({ kind: "none" });
    renderProductionGate();
    await flushAuth();

    // <App> must NOT mount without a session…
    expect(captured.appProps).toBeNull();
    // …and the anonymous surface is the LoginPage (form + username input).
    expect(document.querySelector("form")).not.toBeNull();
    expect(document.querySelector('input[name="username"]')).not.toBeNull();
  });
});
