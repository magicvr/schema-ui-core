// @vitest-environment jsdom
//
// S2 非领域化接入演练（Web/Renderer-Shell 侧，VP-008 §最小可枚举证据面 6）：
// 证明新增标准模块的页面**不需要修改 Renderer/Shell 中央业务注册**——只需
// (1) 模块在 Manifest 贡献 page 条目、(2) 宿主经 `/api/schema/{pageId}` 提供
// 页面文档，Renderer 的通用 schema-driven 路径即渲染该页面。
//
// 本测试在内存 Manifest 中加入一个「宿主从未硬编码过的」probe 页面（pageId
// `probe-items`），并注入 probe schema 文档；App 在 `/probe-items` 处渲染该页，
// 证明中央业务注册无 probe 分支（无 Renderer/Shell 改动）。
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest, type NavigationContext } from "@/protocol/app-manifest";

// 宿主基线 Manifest（home / catalog / catalog-detail）+ 一个 probe 页面。
// probe 的 pageId 在 Renderer 中央业务注册中不存在——本测试证明它仍按
// schema-driven 路径渲染。
function drillManifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "integration", name: "Integration", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
      {
        pageId: "probe-items",
        title: "Probe Items",
        schemaUrl: "/schema/probe-items",
        route: "/probe-items",
      },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [
        { pageRef: "probe-items", label: "Probe Items" },
      ],
    },
  });
}

function probeSchemaDocument() {
  return {
    meta: {
      pageId: "probe-items",
      title: "Probe Items",
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    body: {
      type: "section",
      children: [
        { type: "text", props: { text: "Probe schema body" } },
      ],
    },
  };
}

function probeSchemaFetcher() {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname === "/schema/probe-items") {
      return new Response(JSON.stringify(probeSchemaDocument()), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), {
      status: 404,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

async function renderDrill(path: string, navigationContext?: NavigationContext) {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <App
          manifest={drillManifest()}
          navigationContext={navigationContext}
          schemaFetcher={probeSchemaFetcher()}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("S2 access drill · Renderer/Shell renders a new module page generically", () => {
  it("renders the probe page from manifest + schema without central registration", async () => {
    const container = await renderDrill("/probe-items", {
      user: { roles: ["admin"] },
      features: {},
    });
    expect(container.querySelector("h1")?.textContent).toBe("Probe Items");
    expect(container.textContent).toContain("Probe schema body");
    // 侧边导航也出现 probe 条目（来自 Manifest navigation，非中央注册）。
    expect(container.textContent).toContain("Probe Items");
  });
});
