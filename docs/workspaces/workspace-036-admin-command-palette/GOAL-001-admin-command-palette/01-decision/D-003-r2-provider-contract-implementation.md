---
doc_type: goal-decision
id: D-003-r2-provider-contract-implementation
status: accepted
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# D-003 · R2 provider v1 实现方案

- **日期**：2026-09-14
- **状态**：accepted（承接用户确认的 D-002；非新增范围选择）
- **决定**：不修改 pinned AppManifest 或 Go `kernel.Provider`/Registrar。新增前端 `SearchableProvider` v1 注入 seam 与 `SearchableItem` 聚合器；内置 `manifest` provider 从当前运行时 Manifest 的 `projectNavigation` projection 产生页面/导航目标，再从已通过 D-VAL 的认证页面 Schema body 产生直接 toolbar/actionButton 命令。
- **决定**：provider 输出已本地化的稳定 `id`、`kind`、`label`、`keywords`、`group`、`href` 与可选 page action context；聚合按 provider id + item rank + item id 稳定排序。重复 id 不静默覆盖，冲突的双方均排除并产生非敏感错误。
- **决定**：动作只保留 concrete、rowless、非 upload 的 page-level trigger。行级/批量/未绑定参数不生成全局项；执行时携带原 mount trigger 与 owner page id，经 App pending-action handoff 交给现有 `SchemaCrudProvider.invokeAction`，不直接调用 URL/actionRef。
- **理由**：Manifest 公共协议和后端模块贡献契约均没有 provider/action 元数据；在消费适配器内提供扩展 seam 可保持协议兼容，同时让 mvp/admin/demo/custom 的当前 Manifest 自然承接模块聚合。复用 `projectNavigation`、`loadPageDocument`、permission evaluator 与 existing action executor 可避免新建第二套权限或路由语义。
- **未选方案**：
  1. 在 AppManifest 增加 `providers` / `actions` / `profile` 字段：会触发协议 schema、版本协商和上下游兼容审计，本波不采用。
  2. 从 DB menu/permission catalog 或 `APP_PROFILE` 推导分母：会混入 disabled-profile stale rows，且浏览器没有 profile 真相；本波只认当前运行时 Manifest + authenticated Schema。
  3. 扫描所有顶层 `document.actions` 或 DOM/custom component：顶层定义缺少 mount label/permission，custom component 不属于声明式分母；本波只收录直接 trigger。
- **安全响应**：R2/R3 同时统一 `invokeAction` 的 programmatic visible/permission/cascade gate，移动 custom gate 到 action dispatch 前，修正 actionButton node id fallback，拒绝无 mapping 的未绑定 navigate template，并捕获 row navigation 构造异常；后端 401/403 仍为最终授权边界。
- **审计安排**：R2 provider/聚合为 self 事实审；R3 action invocation + App Shell 为 security/跨边界 scope，采用 `cross`：先 self，再按项目级决策调用 grok build（grok-4.6 · reasoning high）independent；未落盘意见不得作为通过依据。
- **后续**：R3 完成 Command Palette UI、shortcut、ARIA/focus、pending action handoff 与双语主题；R4 完成 profile×permission×route 回归、浏览器证据、关门审计。
