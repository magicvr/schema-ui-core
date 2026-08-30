// golden-web 渲染闭环探针（SSR）：仅经 npm 包 @schema-ui/renderer 渲染 schema 页面文档。
// fixture 形状 = 主线协议校验管线（page.schema.json / RenderPageDocument）验证的页面文档结构。
import { renderToString } from "react-dom/server";
import React from "react";
import assert from "node:assert";
import { I18nProvider, RenderPage, registerCustomComponent } from "@schema-ui/renderer";

// 真实主线形态的最小表单页（对标 dev/examples form 类页面：字段 + 默认值 + 反应式门控语义
// 由 RenderPage 内部 reaction-engine 处理；本 fixture 保持零数据源依赖 = 纯静态渲染）。
const pageDoc = {
  meta: {
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "form.controls.advanced"],
  },
  body: {
    type: "form",
    id: "golden-form",
    props: {
      fields: [
        { id: "name", label: "Name", type: "input", defaultValue: "Golden" },
        {
          id: "role",
          label: "Role",
          type: "input",
          defaultValue: "admin",
          reactions: [{ when: "$deps.name == 'Golden'", fulfill: { value: "operator" } }],
        },
      ],
    },
  },
};

// 扩展接缝（包导出面）应可注册（v0.1 语义验证）。
registerCustomComponent("golden-probe", () => React.createElement("span", null, "golden"));

const html = renderToString(
  React.createElement(
    I18nProvider,
    { locale: "zh-CN" },
    React.createElement(RenderPage, { document: pageDoc, context: {} }),
  ),
);

assert.ok(html.includes("Name"), "label Name rendered");
assert.ok(html.includes("Golden"), "defaultValue rendered");
assert.ok(!html.includes("hybridui-error"), "no error surface");
console.log("golden-web render probe PASS · html bytes:", Buffer.byteLength(html));