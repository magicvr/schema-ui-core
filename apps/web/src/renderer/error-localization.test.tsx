// @vitest-environment jsdom

/**
 * S4 frontend localization floor (GOAL-005 · C4):
 * - auth-client attaches the active locale as Accept-Language.
 * - readResourceApiError parses messageKey/params from the envelope.
 * - The renderer's feedback region prefers the catalog entry (current locale)
 *   over the server message, and falls back to the server message when the key
 *   is unknown (never renders the raw key).
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { authFetch } from "@/account/auth-client";
import { I18nProvider, setActiveLocale } from "@/i18n/runtime";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";
import { readResourceApiError, ResourceApiError } from "@/renderer/resource";

const __dir = dirname(fileURLToPath(import.meta.url));
const MANIFEST_PATH = resolve(__dir, "../test-fixtures/app-manifest.admin-dogfood.json");
const USERS_SCHEMA_PATH = resolve(
  __dir,
  "../../../api/modules/users/schema/users.json",
);

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
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

function usersSchemaDocument(): Record<string, unknown> {
  return JSON.parse(readFileSync(USERS_SCHEMA_PATH, "utf8"));
}

function adminContext() {
  return {
    user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
  };
}

const USERS = {
  items: [
    { id: "usr-1", username: "alice", name: "Alice", roles: ["admin"], updatedAt: "2026-08-03T00:00:00.000Z" },
    { id: "usr-2", username: "bob", name: "Bob", roles: ["editor"], updatedAt: "2026-08-03T11:00:00.000Z" },
  ],
  total: 2,
  page: 1,
  pageSize: 10,
};

// ── C4 · Accept-Language on auth requests ─────────────────────────────────────

describe("S4 · auth client sends the active locale", () => {
  it("attaches Accept-Language to authed fetches", async () => {
    const seen: string[] = [];
    globalThis.fetch = (async (_input: RequestInfo | URL, init?: RequestInit) => {
      seen.push((init?.headers as Headers | undefined)?.get("Accept-Language") ?? "none");
      return new Response(JSON.stringify({ error: "X", message: "m" }), { status: 401 });
    }) as typeof fetch;
    setActiveLocale("zh-CN");
    await authFetch("/api/me");
    expect(seen[0]).toBe("zh-CN");
    setActiveLocale("en-US");
  });
});

// ── C4 · envelope parsing ─────────────────────────────────────────────────────

describe("S4 · readResourceApiError envelope", () => {
  it("parses messageKey and params from the VP-007 envelope", async () => {
    const response = new Response(
      JSON.stringify({
        error: "INVALID_SITE_TITLE",
        message: "站点标题不能为空",
        messageKey: "error.invalidSiteTitle",
        params: { field: "siteTitle" },
        correlation_id: "req-r1-001",
      }),
      { status: 400, headers: { "Content-Type": "application/json" } },
    );
    const error = await readResourceApiError(response, "settings update");
    expect(error).toBeInstanceOf(ResourceApiError);
    expect(error.code).toBe("INVALID_SITE_TITLE");
    expect(error.messageKey).toBe("error.invalidSiteTitle");
    expect(error.params).toEqual({ field: "siteTitle" });
    expect(error.correlationId).toBe("req-r1-001");
    expect(error.message).toContain("request req-r1-001");
    expect(error.message).toContain("HTTP 400");
  });

  it("keeps legacy envelopes working without messageKey", async () => {
    const response = new Response(JSON.stringify({ error: "FORBIDDEN", message: "nope" }), {
      status: 403,
      headers: { "Content-Type": "application/json" },
    });
    const error = await readResourceApiError(response, "list");
    expect(error.messageKey).toBeUndefined();
    expect(error.code).toBe("FORBIDDEN");
    expect(error.fieldErrors).toEqual([]);
  });

  it("falls back to the response header when the body omits correlation_id", async () => {
    const response = new Response(JSON.stringify({ error: "FORBIDDEN", message: "nope" }), {
      status: 403,
      headers: { "Content-Type": "application/json", "X-Request-ID": "header-r1-001" },
    });
    const error = await readResourceApiError(response, "list");
    expect(error.correlationId).toBe("header-r1-001");
    expect(error.message).toContain("request header-r1-001");
  });

  it("parses fieldErrors from the GOAL-014 envelope (A-003 F-002)", async () => {
    const response = new Response(
      JSON.stringify({
        error: "INVALID_PATCH_FIELD",
        message: "invalid patch field: name must not be empty",
        messageKey: "error.invalidPatchField",
        fieldErrors: [{ field: "name", reason: "must not be empty" }],
      }),
      { status: 400, headers: { "Content-Type": "application/json" } },
    );
    const error = await readResourceApiError(response, "dict type update");
    expect(error.code).toBe("INVALID_PATCH_FIELD");
    expect(error.fieldErrors).toEqual([{ field: "name", reason: "must not be empty" }]);
  });
});

// ── C4 · renderer feedback localization floor ─────────────────────────────────

describe("S4 · feedback region localization", () => {
  it("renders the catalog entry by messageKey under zh-CN", async () => {
    const documents = { users: usersSchemaDocument() };
    const fetcher = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const pathname = new URL(String(input), "http://test.local").pathname;
      if (pathname.startsWith("/api/schema/")) {
        return new Response(JSON.stringify(documents.users), { status: 200 });
      }
      if (pathname.startsWith("/api/users")) {
        if (init?.method === "PATCH" || init?.method === "POST" || init?.method === "DELETE") {
          return new Response(
            JSON.stringify({
              error: "INVALID_SITE_TITLE",
              message: "siteTitle must not be empty",
              messageKey: "error.invalidSiteTitle",
            }),
            { status: 400, headers: { "Content-Type": "application/json" } },
          );
        }
        return new Response(JSON.stringify(USERS), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }
      return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
    }) as typeof fetch;
    window.history.replaceState({}, "", "/users");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider stored="zh-CN" browserLanguages={["en-US"]}>
          <App
            manifest={adminManifest()}
            navigationContext={adminContext()}
            schemaFetcher={fetcher}
            resourceFetcher={fetcher}
          />
        </I18nProvider>,
      );
    });
    // A 400 on a page-level write action with messageKey: the feedback region
    // must show the catalog text, not the server English message and not the raw key.
    // W11 · U-05: delete lives in the row "⋯" overflow menu.
    const moreTrigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label]',
    );
    expect(moreTrigger).not.toBeNull();
    await act(async () => {
      moreTrigger!.click();
    });
    const deleteButton = [...document.body.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("删除"),
    );
    expect(deleteButton).toBeDefined();
    await act(async () => {
      deleteButton!.click();
    });
    // Row delete declares a confirm; confirm first, then the request runs.
    const confirmButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("确认"),
    );
    expect(confirmButton).toBeDefined();
    await act(async () => {
      confirmButton!.click();
    });
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const alertAfter = container.querySelector('[role="alert"]');
    const codeEl = alertAfter?.querySelector("[data-feedback-code]");
    expect(codeEl?.getAttribute("data-feedback-code")).toBe("INVALID_SITE_TITLE");
    expect(alertAfter?.textContent).toContain("站点标题不能为空");
    expect(alertAfter?.textContent).not.toContain("INVALID_SITE_TITLE");
    expect(alertAfter?.textContent).not.toContain("siteTitle must not be empty");
  });

  it("falls back to the server message when the messageKey is unknown", async () => {
    const fetcher = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const pathname = new URL(String(input), "http://test.local").pathname;
      if (pathname.startsWith("/api/schema/")) {
        return new Response(JSON.stringify(usersSchemaDocument()), { status: 200 });
      }
      if (pathname.startsWith("/api/users")) {
        if (init?.method === "DELETE") {
          return new Response(
            JSON.stringify({
              error: "MYSTERY_CODE",
              message: "custom server detail",
              messageKey: "error.no.such.key",
            }),
            { status: 400, headers: { "Content-Type": "application/json" } },
          );
        }
        return new Response(JSON.stringify(USERS), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }
      return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
    }) as typeof fetch;
    window.history.replaceState({}, "", "/users");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider stored="zh-CN" browserLanguages={["en-US"]}>
          <App
            manifest={adminManifest()}
            navigationContext={adminContext()}
            schemaFetcher={fetcher}
            resourceFetcher={fetcher}
          />
        </I18nProvider>,
      );
    });
    const moreTrigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label]',
    );
    await act(async () => {
      moreTrigger!.click();
    });
    const deleteButton = [...document.body.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("删除"),
    );
    await act(async () => {
      deleteButton!.click();
    });
    const confirmButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("确认"),
    );
    await act(async () => {
      confirmButton!.click();
    });
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const alertAfter = container.querySelector('[role="alert"]');
    const codeEl = alertAfter?.querySelector("[data-feedback-code]");
    // Unknown key → the translator returns the key itself; the feedback
    // region detects that and keeps the server message (never a raw key).
    expect(codeEl?.getAttribute("data-feedback-code")).toBe("MYSTERY_CODE");
    expect(alertAfter?.textContent).toContain("custom server detail");
    expect(alertAfter?.textContent).not.toContain("MYSTERY_CODE");
  });
});
