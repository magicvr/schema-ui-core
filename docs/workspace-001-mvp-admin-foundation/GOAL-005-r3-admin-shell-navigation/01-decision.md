---
id: GOAL-005-r3-admin-shell-navigation
doc: decision
status: done
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.2.1
---

# 决策记录 · GOAL-005

## 信息需求与阶段门禁

P-005 信息台账维护在 [00-meta.md](00-meta.md)。D-005 已将 `I-005-001` 至 `I-005-005` 的可核对结论落盘并标为 `verified`；后续协议、路由或 shell 边界变化须重新核验对应证据，不得把本目标的 R3 子集扩大为完整协议支持。

父目标的 `I-PROTO-002` 与 `I-PROTO-003` 分别属于 R4/R5 的后续门禁。本目标不修改其状态，也不以 `I-PROTO-001=verified` 替代它们。

## D-001 · 立项 R3 Admin 外壳与导航子目标

**日期**：2026-07-31
**状态**：accepted

**决定**：

在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-005-r3-admin-shell-navigation`，将 R3 范围限定为 App manifest 装载、Admin shell、导航入口和路由语义规划/实施；Root 仍保持 `active`，纲领进度仍为 `2/6`。

**为什么**：

- Root 路线图已把 R3 定义为“Admin 外壳与导航”，R2 的 D-009/A-006 已记录冻结后进入 R3 规划。
- R1 已明确不包含 Admin 导航壳、多业务路由和协议 Renderer 全量；R3 需要有独立的范围、信息项和审计台账。
- 协议清单将 `D-APP` 的 React 主责定义为装载与导航壳，能够为本目标提供范围锚点，但不能替代本地实现决策。

**未选方案**：

- 在 R1 目标中直接补入 Admin shell 或业务路由：会越过既有 R1/R3 边界。
- 把账号权限、Renderer 全量或业务范例一起纳入本目标：会提前吞并 R4/R5，且绕过其 required 信息门禁。
- 以硬编码业务路由树代替 manifest/导航契约：无法证明与 `S-09`、`app-navigation` 语义一致。

**影响**：

本目标进入 `active` 规划阶段；当前不改 `apps/web`，不改变父目标或协议门禁状态。

## D-002 · 采用“信息就绪 → 方案冻结 → 实施 → 验证关门”的 R3 路线

日期：2026-07-31
**状态**：accepted

**决定**：

先处理 `I-005-001` 至 `I-005-005`，再冻结 manifest 最小子集、路由映射、默认/fallback/active-route 语义和 shell 产品边界；在此之前不把待确认取舍写成实现事实。实现完成后必须补结构/行为/运行时证据和阶段自审，才可讨论 `done`。

**为什么**：

当前 Web 仍是 R1 单页占位，`main.tsx` 没有 router，仓库也没有本地 manifest loader、navigation 或 shell 实现。现有协议资料只固定了上游路径与验证方向，不能据此推断本地路由行为。

**未选方案**：

- 先写一个临时 shell，再事后补 manifest 与路由语义：会把未知项伪装成已决定行为。
- 仅以构建成功或页面可打开作为 R3 关门证据：无法覆盖 manifest/navigation 契约。

**影响**：

开放 required 信息项是 R3 方案冻结和受影响实施门禁；必要时须按 P-004 由用户裁决 fixed、accepted-residual 或 user-overruled，不能静默放行。

## D-003 · 采用固定上游资料与明确的验证证据边界

日期：2026-07-31

状态：accepted

**决定**：

R3 规划以 `protocol-inventory-v2.7.0.md` 登记的 source commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 为资料定位锚点，优先对照 `S-09`、`app-manifest` 和 `app-navigation`；未来验收须记录 schema/fixture 或等价可核对证据。`conformance/reference-js`、`reference-python` 和 `runner` 仅可作为参考，不能单独证明兼容。

**为什么**：

协议清单明确区分 structural contract、behavioral fixture 和 excluded reference/runner。将资料版本和证据类别写入计划，可以避免把“找到路径”或“调用参考实现”误写成 R3 已验证。

**影响**：

该决定固定证据方向，不代表当前已有本地 schema、fixture、运行时或 conformance 结果。

## D-004 · 响应 A-001，统一方案冻结门禁并显式登记父级依赖

**日期**：2026-07-31
**状态**：accepted

**决定**：

采用 A-001 对 F-001 的推荐修正：将 `I-005-003` 和 `I-005-004` 的最晚需要阶段统一为「方案冻结前」，使信息表与 D-002 的“先处理 `I-005-001` 至 `I-005-005`，再冻结”保持一致。采纳 F-002：`I-005-001` 显式依赖 Root `I-PROTO-004` 的 vendor vs pin 工程策略；关闭该项时必须记录所选校验方式和失败边界，但不把 Root 的 non-blocking 项静默升格为本目标独立 required。

**为什么**：

- A-001 已证明原信息表与 D-002 对方案冻结的要求存在矛盾；采用推荐修正可使编排器只有一个可核对的冻结门槛。
- `app-manifest` 的本地校验方式会影响 `I-005-001` 的证据边界；显式关联 Root `I-PROTO-004` 可避免把“已知上游路径”误写成“本地已可校验”。

**影响**：

这只闭合 A-001 的 F-001 文档一致性 finding，不验证或关闭任何 `I-005-*`。R3 仍不得方案冻结、实施或关门；用户在未来拟用独立审计推进方案冻结前，仍须按 P-004.1 决定是否进行同 scope 自审。

## D-005 · 冻结 R3 manifest、导航、路由与 shell 子集

**日期**：2026-07-31
**状态**：accepted
**关联信息项**：`I-005-001` 至 `I-005-005` → `verified`

**决定**：

1. **协议与来源**：R3 采用 `schema-ui-docs@2.7.0`、commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` 作为固定来源。仓库内保留 `docs/schemas/app-manifest.schema.json` 与 app-manifest/app-navigation fixture 副本，并以 `apps/web/src/protocol/upstream/provenance.json` 记录来源、相对路径和 SHA-256；运行时不访问远程协议源。生产 validator 是针对 R3 2.7 host subset 的 fail-fast TypeScript 实现，不宣称完整 JSON Schema 或完整 conformance 支持。
2. **装载与失败边界**：`loadAppManifest()` 只接受 HTTP 成功且能通过 `validateAppManifest()` 的 JSON；HTTP、解析、协议版本、必需 capability、字段和导航结构错误均以 `ManifestError` fail closed，`main.tsx` 渲染 `ManifestFailure`，不渲染不可信页面。
3. **导航映射**：manifest 的 `navigation.top`、`navigation.sidebar`、`navigation.user` 分别投影到 header、desktop sidebar/mobile nav 和 user navigation。`pageRef` 解析到注册 page 的 route；`url` 保留为站内 URL；group 只保留可见 child，空 group 被移除；`visibleWhen`/`permissions.view` 只做当前上下文的结构过滤，不构成 R4 身份或权限实现。
4. **路由与 active**：D4a 依次按 literal 数量、route 模板长度、声明顺序选择匹配；URL query 不参与 route match，空或尾斜杠 segment 不被规范化。`/` 在有非参数 `homePageRef` 时用 `history.replaceState` 进入 home；已知 deep link 优先；站内链接用 `pushState`；未知路径保留并显示 fallback，可返回 home。参数化 `pageRef` 只有当前路径能提供绑定参数时才生成具体 href，但仍可作为 active 匹配目标。
5. **shell 边界**：R3 固定 app identity/header、top/user navigation、responsive mobile navigation、desktop sidebar、main page surface 和 route fallback。page schema renderer、业务页面、真实身份/权限上下文、R4 权限产品化、R5 全量 renderer、完整主题系统和完整协议支持均留在后续目标。
6. **上游 fixture 对照**：测试执行 35/37 个 app-manifest cases 与 16/16 个 app-navigation cases；4 个 negotiation 与 1 个 decoupled-version case 使用明确标注的 fixture-only adapter，不改变生产 host API。仅排除两条 upstream M1 schema error-envelope case，因为上游使用 `CAPABILITY_REQUIRED` 聚合错误，而 R3 fail-fast host 使用 `MISSING_REQUIRED_CAPABILITY`；排除原因由机器断言和测试台账固定。

**为什么**：

- 这些结论直接对应 `I-005-001` 至 `I-005-005`，并能由实现、固定 artifact、fixture 测试和 App 集成测试复核。
- 2.7 exact host 与 fixture-only negotiation 分离，避免将旧版本协商样例误写成生产兼容承诺。
- 明确 empty/unknown/fallback、参数 href、空上下文和 R4/R5 边界，避免实现细节在关门时被误读为更大的产品承诺。

**未选方案**：

- **运行时远程拉取 schema 或 fixture**：不可复现且会使固定来源之外的网络状态影响验证；采用仓库内 pinned copy + hash。
- **把 2.5 negotiation 当作生产 loader 能力**：R3 生产 host 保持 exact `2.7`，协商样例仅在 fixture adapter 中验证。
- **为参数化 pageRef 伪造默认参数**：会生成不可证明的 URL；无绑定时不提供 href。
- **把空 navigation context 解释为 R4 鉴权**：R3 只定义结构过滤和可注入测试 context，真实身份/权限来源留给 R4。

**影响**：

本决策冻结 R3 实施与验收边界，允许响应 A-004 F-001；不修改 Root `I-PROTO-004`、`I-PROTO-002` 或 `I-PROTO-003`，也不直接将 GOAL-005 标为 `done`。实施事实和关门证据分别记录在 `02-execution.md` 与 `03-audit.md`。
