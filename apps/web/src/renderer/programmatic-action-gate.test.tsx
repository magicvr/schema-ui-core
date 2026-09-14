// @vitest-environment jsdom

import { act, useEffect } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { registerCustomComponent, resetCustomComponentsForTests } from "@/renderer/custom-components";
import { RenderPage, useSchemaCrud } from "@/renderer/render.tsx";
import { I18nProvider } from "@/i18n/runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  resetCustomComponentsForTests();
});

afterEach(async () => {
  resetCustomComponentsForTests();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function invokeOnMount(trigger: Record<string, unknown>) {
  return function InvokeOnMount() {
    const crud = useSchemaCrud();
    useEffect(() => {
      crud?.invokeAction(trigger, null);
      // This test component intentionally invokes once; depending on the
      // context object would repeat after feedback/modal state updates.
      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    return null;
  };
}

function renderPage(pageDocument: Record<string, unknown>, context: Record<string, unknown>): HTMLDivElement {
  registerCustomComponent("invoke-on-mount", invokeOnMount(pageDocument.__trigger as Record<string, unknown>));
  const container = globalThis.document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider stored="en-US">
        <RenderPage document={pageDocument as never} context={context} />
      </I18nProvider>,
    );
  });
  return container;
}

async function flush(): Promise<void> {
  await act(async () => {
    await Promise.resolve();
  });
}

describe("programmatic action gate", () => {
  it("denies a permissioned modal invoked with the actionButton node id", async () => {
    const page = {
      __trigger: {
        actionId: "openCreate",
        key: "create-user-button",
        permissionIntent: "edit",
        label: "Create user",
      },
      meta: {
        pageId: "users",
        title: "Users",
        protocolVersion: "2.7",
        requiredCapabilities: ["permissions.inheritance", "actions.page.trigger"],
      },
      actions: {
        openCreate: { type: "modal", content: { type: "section", children: [] } },
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
            props: { actionId: "openCreate", permissionIntent: "edit", label: "Create user" },
          },
          { type: "custom", component: "invoke-on-mount" },
        ],
      },
    };
    const container = renderPage(page, { user: { permissions: [] }, features: {} });
    await flush();

    expect(container.querySelector('[role="dialog"]')).toBeNull();
    expect(container.querySelector('[data-feedback-code="ACTION_NOT_EXECUTED"]')).not.toBeNull();
  });

  it("allows the same permissioned modal only after the page permission is true", async () => {
    const page = {
      __trigger: {
        actionId: "openCreate",
        key: "create-user-button",
        permissionIntent: "edit",
        label: "Create user",
      },
      meta: {
        pageId: "users",
        title: "Users",
        protocolVersion: "2.7",
        requiredCapabilities: ["permissions.inheritance", "actions.page.trigger"],
      },
      actions: {
        openCreate: { type: "modal", content: { type: "section", children: [] } },
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
            props: { actionId: "openCreate", permissionIntent: "edit", label: "Create user" },
          },
          { type: "custom", component: "invoke-on-mount" },
        ],
      },
    };
    const container = renderPage(page, { user: { permissions: ["users.write"] }, features: {} });
    await flush();

    expect(container.querySelector('[role="dialog"][aria-label="Create user"]')).not.toBeNull();
  });

  it("checks local toolbar permissions before a custom export handler", async () => {
    let calls = 0;
    const fetcher = (async () => {
      calls += 1;
      return new Response("csv", { status: 200 });
    }) as typeof fetch;
    const page = {
      __trigger: {
        actionRef: "exportUsers",
        key: "export",
        permissions: { edit: '$context.user.permissions contains "data.export"' },
        label: "Export",
      },
      meta: {
        pageId: "users",
        title: "Users",
        protocolVersion: "2.7",
        requiredCapabilities: ["permissions.inheritance", "actions.page.trigger"],
      },
      actions: { exportUsers: { type: "custom", handler: "export.users" } },
      body: {
        type: "section",
        children: [
          {
            type: "table",
            id: "users-table",
            props: {
              toolbar: [
                {
                  key: "export",
                  label: "Export",
                  actionRef: "exportUsers",
                  permissions: { edit: '$context.user.permissions contains "data.export"' },
                },
              ],
            },
          },
          { type: "custom", component: "invoke-on-mount" },
        ],
      },
    };
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    registerCustomComponent("invoke-on-mount", invokeOnMount(page.__trigger));
    await act(async () => {
      root.render(
        <I18nProvider stored="en-US">
          <RenderPage document={page as never} context={{ user: { permissions: [] }, features: {} }} dataFetcher={fetcher} />
        </I18nProvider>,
      );
    });
    await flush();

    expect(calls).toBe(0);
    expect(container.querySelector('[data-feedback-code="ACTION_NOT_EXECUTED"]')).not.toBeNull();
  });
});
