// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { defaultBranding } from "@/app/branding";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest } from "@/protocol/app-manifest";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  window.history.replaceState({}, "", "/home");
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

function manifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "palette-app", name: "Palette app", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
      { pageId: "users", title: "Users", schemaUrl: "/schema/users", route: "/users" },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [
        {
          label: "Admin",
          items: [
            {
              pageRef: "users",
              label: "Users",
              visibleWhen: { when: "$context.features.menu_users == true" },
            },
          ],
        },
      ],
    },
  });
}

function schemaDocument(pageId: string, body: Record<string, unknown>, actions?: Record<string, unknown>) {
  return {
    meta: {
      pageId,
      title: pageId,
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    ...(actions === undefined ? {} : { actions }),
    body,
  };
}

const HOME_DOCUMENT = schemaDocument("home", {
  type: "section",
  children: [{ type: "text", props: { text: "Home body" } }],
});

const USERS_DOCUMENT = {
  meta: {
    pageId: "users",
    title: "Users",
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation", "actions.page.trigger", "permissions.inheritance"],
  },
  actions: {
    openCreate: {
      type: "modal",
      content: {
        type: "form",
        id: "create-user-form",
        permissionCascade: { keys: ["edit"] },
        permissions: { edit: '$context.user.permissions contains "users.write"' },
        props: {
          fields: [{ id: "name", type: "input", label: "Name" }],
          submitLabel: "Create",
        },
      },
    },
  },
  body: {
    type: "section",
    id: "users-root",
    permissionCascade: { keys: ["edit"] },
    permissions: { edit: '$context.user.permissions contains "users.write"' },
    children: [
      {
        type: "actionButton",
        id: "create-user-button",
        props: {
          label: "Create user",
          actionId: "openCreate",
          permissionIntent: "edit",
        },
      },
    ],
  },
};

function schemaFetcher(input: RequestInfo | URL): Promise<Response> {
  const pathname = new URL(String(input), "http://test.local").pathname;
  const document = pathname === "/schema/home" ? HOME_DOCUMENT : pathname === "/schema/users" ? USERS_DOCUMENT : undefined;
  return Promise.resolve(
    document === undefined
      ? new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 })
      : new Response(JSON.stringify(document), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        }),
  );
}

async function renderApp(permissions: string[]): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <App
          manifest={manifest()}
          navigationContext={{
            user: { id: "u1", permissions },
            features: { menu_users: true },
          }}
          schemaFetcher={schemaFetcher as typeof fetch}
          branding={defaultBranding()}
        />
      </I18nProvider>,
    );
  });
  return container;
}

async function flush(): Promise<void> {
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 0));
  });
}

describe("App command palette integration", () => {
  it("opens from Ctrl+K, indexes the visible page action, and preserves internal navigation", async () => {
    const container = await renderApp(["users.write"]);
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "k", ctrlKey: true, bubbles: true }));
    });
    await flush();

    expect(container.querySelector('[role="dialog"]')).not.toBeNull();
    const action = [...container.querySelectorAll('[role="option"]')].find((option) =>
      option.textContent?.includes("Create user"),
    ) as HTMLButtonElement | undefined;
    expect(action).not.toBeUndefined();

    await act(async () => {
      action?.click();
    });
    await flush();
    await flush();

    expect(window.location.pathname).toBe("/users");
    expect(container.querySelector('[role="dialog"][aria-label="Create user"]')).not.toBeNull();
  });

  it("opens a visible page and focuses its heading, while ignoring Ctrl+K in the search input", async () => {
    const container = await renderApp(["users.write"]);
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-command-palette-trigger]")?.click();
    });
    await flush();
    const pageOption = [...container.querySelectorAll('[role="option"]')].find((option) =>
      option.textContent?.trim().startsWith("Users"),
    ) as HTMLButtonElement | undefined;
    expect(pageOption).not.toBeUndefined();
    await act(async () => {
      pageOption?.click();
    });
    await flush();
    expect(window.location.pathname).toBe("/users");
    expect(document.activeElement).toBe(container.querySelector("#page-title"));
    const adminGroupButton = container.querySelector<HTMLButtonElement>(
      'section[data-navigation-group="Admin"] button',
    );
    expect(adminGroupButton?.getAttribute("aria-expanded")).toBe("true");

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-command-palette-trigger]")?.click();
    });
    await flush();
    const input = container.querySelector<HTMLInputElement>('[role="combobox"]')!;
    await act(async () => {
      input.dispatchEvent(new KeyboardEvent("keydown", { key: "k", ctrlKey: true, bubbles: true }));
    });
    expect(container.querySelector('[role="dialog"]')).not.toBeNull();
    await act(async () => {
      input.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
    });
  });

  it("opens from Meta+K outside editable targets", async () => {
    const container = await renderApp(["users.write"]);
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "k", metaKey: true, bubbles: true }));
    });
    await flush();
    expect(container.querySelector('[role="dialog"]')).not.toBeNull();
  });

  it("does not expose a denied action while retaining the visible page result", async () => {
    const container = await renderApp([]);
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-command-palette-trigger]")?.click();
    });
    await flush();

    const options = [...container.querySelectorAll('[role="option"]')];
    expect(options.some((option) => option.textContent?.includes("Users"))).toBe(true);
    expect(options.some((option) => option.textContent?.includes("Create user"))).toBe(false);
  });
});
