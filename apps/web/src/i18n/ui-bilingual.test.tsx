// @vitest-environment jsdom

/**
 * S2 bilingual coverage (GOAL-003 · C1–C4).
 *
 * Drives the REAL shipped surfaces — LoginPage / App shell / SchemaTable /
 * FormControls / manifest key resolution — through the real catalogs and the
 * real module schema fixtures, under both locales. Also covers M4 (missing
 * key observable + safe fallback).
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { LoginPage } from "@/app/LoginPage";
import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { MISSING_TRANSLATION_EVENT, resetMissingTranslationReports } from "@/i18n/catalog";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const __dir = dirname(fileURLToPath(import.meta.url));
const MANIFEST_PATH = resolve(__dir, "../test-fixtures/app-manifest.admin.json");
const USERS_SCHEMA_PATH = resolve(
  __dir,
  "../../../api/internal/modules/users/schema/users.json",
);

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
  resetMissingTranslationReports();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
  resetMissingTranslationReports();
});

function adminManifest(): AppManifest {
  return validateAppManifest(JSON.parse(readFileSync(MANIFEST_PATH, "utf8")));
}

function usersSchemaDocument(): Record<string, unknown> {
  return JSON.parse(readFileSync(USERS_SCHEMA_PATH, "utf8"));
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

function fetcherFor(documents: Record<string, unknown>): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname.startsWith("/api/schema/")) {
      const pageId = pathname.slice("/api/schema/".length);
      const document = documents[pageId];
      if (document === undefined) {
        return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 });
      }
      return new Response(JSON.stringify(document), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/users")) {
      return new Response(JSON.stringify(USERS), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
  }) as typeof fetch;
}

async function renderAt(
  path: string,
  element: React.ReactElement,
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(element);
  });
  return container;
}

function shell(children: React.ReactNode, locale: "zh-CN" | "en-US") {
  return (
    <I18nProvider stored={locale} browserLanguages={["en-US"]}>
      {children}
    </I18nProvider>
  );
}

/** Admin context with permissions so schema-gated toolbar/row actions enable. */
function adminContext() {
  return {
    user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
  };
}

// ── C1 · fixed UI chrome under both locales ──────────────────────────────────

describe("S2 · fixed UI chrome (C1)", () => {
  it("login page renders zh-CN chrome with the zh locale", async () => {
    const container = await renderAt(
      "/",
      shell(<LoginPage onLogin={async () => undefined} />, "zh-CN"),
    );
    expect(container.querySelector("h3")?.textContent).toBe("登录");
    expect(container.querySelector('label[for="username"]')?.textContent).toBe("用户名");
    expect(container.querySelector('label[for="password"]')?.textContent).toBe("密码");
    expect(container.querySelector('button[type="submit"]')?.textContent).toBe("登录");
    expect(container.textContent).toContain("本地开发种子账号");
    expect(document.documentElement.lang).toBe("zh-CN");
  });

  it("login page renders en-US chrome with the en locale", async () => {
    const container = await renderAt(
      "/",
      shell(<LoginPage onLogin={async () => undefined} />, "en-US"),
    );
    expect(container.querySelector("h3")?.textContent).toBe("Sign in");
    expect(container.querySelector('label[for="username"]')?.textContent).toBe("Username");
    expect(document.documentElement.lang).toBe("en-US");
  });

  it("maps stable login error codes to localized text", async () => {
    const container = await renderAt(
      "/",
      shell(
        <LoginPage
          onLogin={async () => {
            throw new Error("boom");
          }}
        />,
        "zh-CN",
      ),
    );
    const username = container.querySelector<HTMLInputElement>('input[name="username"]')!;
    const password = container.querySelector<HTMLInputElement>('input[name="password"]')!;
    const form = container.querySelector("form")!;
    await act(async () => {
      username.value = "alice";
      password.value = "secret";
      form.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    });
    // The generic failure path (non-AuthError) maps to the generic key.
    expect(container.querySelector('[role="alert"]')?.textContent).toBe("登录失败");
  });
});

// ── C2 · manifest titleKey/labelKey real resolution ──────────────────────────

describe("S2 · manifest key resolution (C2)", () => {
  it("resolves the users page title through titleKey under zh-CN", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    expect(container.querySelector("h1")?.textContent).toContain("用户");
    // Nav label resolved through labelKey.
    expect(container.textContent).toContain("用户");
  });

  it("resolves the users page title through titleKey under en-US", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "en-US",
      ),
    );
    expect(container.querySelector("h1")?.textContent).toContain("Users");
  });

  it("falls back to the literal title when titleKey is absent from the manifest", async () => {
    const manifest = adminManifest();
    const page = manifest.pages.find((entry) => entry.pageId === "users")!;
    delete (page as { titleKey?: string }).titleKey;
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={manifest}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    expect(container.querySelector("h1")?.textContent).toContain("Users");
  });
});

// ── C3 · schema-driven text resolution ───────────────────────────────────────

describe("S2 · schema key resolution (C3)", () => {
  it("renders users schema labels via labelKey under zh-CN", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    // Column headers resolved from schema labelKey.
    const text = container.textContent ?? "";
    expect(text).toContain("用户名");
    expect(text).toContain("姓名");
    expect(text).toContain("更新时间");
    // Toolbar label from labelKey.
    expect(text).toContain("新建用户");
  });

  it("renders users schema labels under en-US", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "en-US",
      ),
    );
    const text = container.textContent ?? "";
    expect(text).toContain("Username");
    expect(text).toContain("New user");
  });

  it("resolves modal form field labels and submit labels under zh-CN (real fixture actions)", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    const newUserButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("新建用户"),
    )!;
    await act(async () => {
      newUserButton.click();
    });
    const text = container.textContent ?? "";
    expect(text).toContain("用户名");
    expect(text).toContain("创建用户");
  });

  it("resolves user recordView title and field labels under zh-CN", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });
    const nameCells = Array.from(container.querySelectorAll("td")).filter((td) =>
      td.textContent?.includes("Alice"),
    );
    expect(nameCells.length).toBeGreaterThan(0);
    await act(async () => {
      nameCells[0]!.closest("tr")?.click();
    });
    const panel = container.querySelector('[data-record-view="panel"]');
    expect(panel).not.toBeNull();
    expect(panel?.getAttribute("aria-label")).toBe("用户详情");
    expect(panel?.querySelector("h2")?.textContent).toBe("用户详情");
    const labels = [...(panel?.querySelectorAll("dt") ?? [])].map((el) => el.textContent);
    expect(labels).toContain("用户名");
    expect(labels).toContain("姓名");
    expect(labels).toContain("创建时间");
    expect(labels).not.toContain("username");
    expect(labels).not.toContain("updatedAt");
  });

  it("resolves delete confirm through confirmKey under zh-CN", async () => {
    const documents = { users: usersSchemaDocument() };
    const container = await renderAt(
      "/users",
      shell(
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor(documents)}
          resourceFetcher={fetcherFor(documents)}
        />,
        "zh-CN",
      ),
    );
    // Wait for the list to load, then click the row Delete action.
    // W11 · U-05: delete lives in the row "⋯" overflow menu.
    const moreTrigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label]',
    );
    expect(moreTrigger).not.toBeNull();
    await act(async () => {
      moreTrigger!.click();
    });
    const deleteButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("删除"),
    );
    expect(deleteButton).toBeDefined();
    await act(async () => {
      deleteButton!.click();
    });
    expect(container.textContent).toContain("删除该用户？");
  });
});

// ── C4 · M4 missing-key flow ─────────────────────────────────────────────────

describe("S2 · M4 missing translation key (C4)", () => {
  it("reports the missing key, renders the literal fallback, and keeps the flow usable", async () => {
    const schema = usersSchemaDocument();
    // Manufacture a missing key on the first column label: no catalog entry in
    // either locale, so the literal label must be the final fallback.
    const table = (schema.body as { children: Array<{ type?: string; props?: { columns?: unknown[] } }> }).children.find(
      (child) => child.type === "table",
    )!;
    const columns = table.props!.columns!;
    (columns[1] as { labelKey?: string }).labelKey = "no.such.key.anywhere";
    const events: Array<{ locale: string; key: string }> = [];
    const handler = (event: Event) => {
      events.push((event as CustomEvent).detail);
    };
    window.addEventListener(MISSING_TRANSLATION_EVENT, handler);
    try {
      const documents = { users: schema };
      const container = await renderAt(
        "/users",
        shell(
          <App manifest={adminManifest()} schemaFetcher={fetcherFor(documents)} resourceFetcher={fetcherFor(documents)} />,
          "zh-CN",
        ),
      );
      // Literal text fallback (protocol text) — never empty, never a crash.
      expect(container.textContent).toContain("Username");
      expect(events.some((event) => event.key === "no.such.key.anywhere")).toBe(true);
      // The main flow still works: rows render.
      expect(container.textContent).toContain("alice");
      // Toolbar still renders (labeled via a valid key).
      expect(container.textContent).toContain("新建用户");
    } finally {
      window.removeEventListener(MISSING_TRANSLATION_EVENT, handler);
    }
  });
});
