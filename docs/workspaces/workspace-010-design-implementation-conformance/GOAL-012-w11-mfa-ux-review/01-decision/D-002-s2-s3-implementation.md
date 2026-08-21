---
id: D-002-s2-s3-implementation
doc: decision-entry
goal: GOAL-012-w11-mfa-ux-review
date: 2026-08-15
status: accepted
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# D-002 · S2/S3 实施决策

## 决策

1. **M-01 二维码渲染（accepted）**：引入 npm 依赖 qrcode-generator@1.5.2（MIT、零依赖、纯 JS）；新增 `apps/web/src/components/qr-code.tsx`（SVG 渲染模块矩阵，无 canvas 依赖，jsdom 可测）；mfa-manager 绑定流程展示二维码 + 手动密钥兜底。npm audit 仅既有传递依赖 nanoid 高危（与本波无关，已记录待后续处理）。
2. **M-02/M-03 修复方案（accepted，按 D-001 I-001 方向）**：
   - 后端：自服务端点（confirm/disable/rotate）ErrMFAInvalid 改映射 400；登录二步验证保持 401（新增 writeSelfServiceMFAError 与 writeMFAError 分轨）。
   - 前端：解绑成功后设 sessionStorage["mfa.disabledNotice"] → 本地 logout() → 登录页一次性提示横幅（可关闭）；错误码留页重填。
3. **U-01/U-02 动态选项（accepted，本地扩展）**：
   - 表单控件新增 `optionsSource` + `optionsMapping`（select/radio/checkboxGroup 通用；响应 {items:[...]} 或裸数组；加载中回退静态 options；非法源/失败 fail-closed 空集；允许 query 串，禁止 scheme/host/#）。
   - 后端新增只读目录端点 GET /api/permissions 与 GET /api/menu-items（roles.read 门禁，admin.roles provider 注册）：数据源与授权校验同源（permissions/menu_items 表），UI 可选项不会出现后端拒绝的键。
   - users.json 角色分配：textarea 逗号分隔 → checkboxGroup 动态选项（/api/roles?pageSize=1000，value=key/label=name）。
   - roles.json 权限/菜单：硬编码选项 → 动态选项（/api/permissions、/api/menu-items）。
   - **范围取舍**：权限「按模块分组 + 全选」矩阵推迟（P2 遗留，见 U-02 期望的完整形态）；本波以「全量可勾选」为达成标准（解除新模块无法授权阻断）。
4. **协议影响（accepted）**：optionsSource 为本地扩展，协议 pin v2.8.0 不变；渲染层字段解析白名单已同步放行（render.ts），扩展按 fail-closed 原则对未知端点零生效。

## 未选方案

- 前端硬编码全部权限/角色选项：与后端 reconciled 目录脱钩，新模块每次都要改 schema。
- 目录端点做成公开无鉴权：权限名非机密但按最小暴露原则保留 roles.read 门禁。
- 用 schema-ui-docs 上游控件替代：上游无动态选项控件，扩展以本地能力落地并在文档标注。
