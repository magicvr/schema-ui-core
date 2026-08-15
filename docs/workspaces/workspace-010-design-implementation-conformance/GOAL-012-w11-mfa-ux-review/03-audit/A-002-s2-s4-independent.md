---
id: A-002-s2-s4-independent
doc: audit-entry
record_id: A-002
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
scope: S2/S3/S4 实施（M-01～M-03 + U-01～U-07；含 S2 认证/MFA security 门禁）
verdict: conditional
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# A-002 · 独立交叉审计 · S2/S3/S4 实施（2026-08-15）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high）
- **类型** / **scope**：execution-facts · S2 MFA 缺陷修复 + S3 UX P0 + S4 UX P1（M-01～M-03、U-01～U-07）；S2 认证语义为 security 门禁
- **verdict**：**conditional**（安全实施面无 high required；I-004 到期未闭合，阻断无条件 S5 关门）
- **工作区**：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；`canonical_scope` 已校验；`shared_materials_catalog: none`；`primary_plan` = VP-010）

## 范围与区间

- **covered**：GOAL-012 `00-meta` / `01-decision.md` / D-001 / D-002 / E-001～E-004 / `03-audit.md`；`mfa.go` 分轨映射；`auth-client.ts` / `AuthContext.tsx` / `mfa-manager.tsx` / `qr-code.tsx` / `LoginPage.tsx`；`rbac_catalog.go` + `roles_repository.go` 目录；`roles/provider.go` + `kernel/profile.go`；`form-controls.ts/.tsx` `optionsSource`；`render.ts` 白名单 / `render.tsx` FeedbackRegion + `searchFormSubmit`；`schema-table.tsx` 行操作/分页；`data-table.tsx` 空状态；users/roles/datadictionary/filelibrary/scheduledtasks/recyclebin schema；对应单测。
- **独立复跑（2026-08-15）**：
  - `apps/api` `go test ./internal/handler/ ./internal/modules/mfa/ ./internal/modules/roles/ ./internal/modules/authsession/ -count=1` **exit 0**
  - `apps/web` 切片 6 文件 **73/73 PASS**（mfa-manager 5、qr-code 3、LoginPage 12、dynamic-options 4、schema-table 27、schema-crud 22）
- **excluded**：未复跑 Web 全量 1002 / `tsc` / `go test ./...` 全包（采信 E-004 叙述 + 本轮切片绿）；未跑 Playwright / 活栈扫码点验；未改 status / progress / goal-tree / 方案正文 / 00-meta / D-* / E-*。
- **P-005**：I-001～I-003 closed（D-001/D-002 + E-002）；**I-004 required · 最晚阶段 S4 · 仍 open**（无 D-003、无 residual）→ F-001。
- **共享资料**：`none`，无引用可核。
- **编号**：索引此前无 A 条目；用户书面指定落盘 `A-002`（空洞可保留）。

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 自服务 ErrMFAInvalid 401→400；登录 verify 保持 401 | `mfa.go` `writeSelfServiceMFAError` L266–278（confirm/disable/rotate）；`writeMFAError` L248–249 仍 401；`TestMFASelfService` L198–224 断言 400；`TestMFALoginTwoStep` L152 仍断言 401 |
| authFetch 仅 401 清会话；400 不登出 | `auth-client.ts` L193–208；`mfa-manager.test.tsx` 错码 confirm/disable 不调用 logout |
| 解绑成功仍全会话吊销 | `mfa.go` L173 `BumpTokenVersionAndRevokeAll`；`TestMFASelfService` L231–232；前端先 `sessionStorage` 再 `logout()`（`mfa-manager.tsx` L163–168） |
| 真实服务仍校验 TOTP/恢复码 | `service.go` `requireActiveSecondFactor` L257–276；Disable/RotateRecovery 均走该门 |
| M-01 二维码：MIT 零运行时依赖、SVG、otpauthURL | `package.json` qrcode-generator@^1.5.2；`node_modules/qrcode-generator/package.json` license MIT、无 runtime deps；`qr-code.tsx` 只画 rect；`mfa-manager.tsx` L202–206 `value={enrollPayload.otpauthURL}` |
| 目录端点 roles.read + 声明路由一致 | `rbac_catalog.go` L31/L51；`provider.go` L47–48 + L70；`profile.go` L163；`provider_test.go` L65 |
| 目录与授权表同源 | `roles_repository.go` L447 `SELECT key, description FROM permissions`；L471 `SELECT id, page_ref FROM menu_items WHERE enabled = 1` |
| optionsSource 非法源/失败 fail-closed | `form-controls.tsx` L60–63 / L93–96；dynamic-options 测试 4 例绿 |
| U-04 搜索绑定 q；U-05 权限/disabledWhen 保留；U-06 分页；U-07 空状态 | `render.tsx` `searchFormSubmit` L1002–1014；`schema-table.tsx` L332–333 / L646–648 / L802–807 / L862–869；`data-table.tsx` L259–268 |
| 无越界：协议 schema 未改；Profile 仅 admin.roles 路由键增量 | `docs/schemas/` 本波未改 optionsSource 对象契约；`profile.go` 仅 roles 路由列表追加两条 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S2 M-01 二维码 | 达成（静区有残余，F-005） | 组件 + 3 用例；otpauthURL 编码 |
| S2 M-02/M-03 错码不登出 / 解绑提示+吊销 | 达成 | 分轨 400 + 前端回归 + 吊销未削弱 |
| S3 U-01/U-02 动态选项 + 目录 | 达成（协议同名不同形见 F-002） | 本地扩展 + 目录端点 + schema |
| S4 U-03～U-07 | 大体达成（I-004 未闭合；回收站搜索缺口；Toast 仍占流） | E-004 + 代码；F-001/F-003/F-004 |
| 回归切片 | 本轮独立复跑绿 | Go 四包 exit 0；Web 73/73 |
| I-004 S4 方案门禁 | **未达成** | 00-meta / 01-decision 仍 open；无 D-003 |

### 审计重点逐项

1. **M-02/M-03 认证语义（security）**：分轨正确。自服务无效码是已鉴权用户的业务校验失败（400），`authFetch` 不会 refresh/清会话。登录 `/api/auth/mfa/verify` 仍走 `writeMFAError` → 401，`mfaVerify` 不经 `authFetch`，`TestMFALoginTwoStep` 仍钉 401。Disable 成功路径仍先 `BumpTokenVersionAndRevokeAll` 再 204，吊销语义未削弱。`sessionStorage["mfa.disabledNotice"]` 只存 `"1"`，LoginPage 读后删除并渲染 i18n 文案（不回显存储值）；同 origin、无提权面；XSS 最多伪造一条无特权横幅。恢复码：生产路径 `requireActiveSecondFactor` 仍校验；空码/错码映射 400。未见 400/401 边界绕过。
2. **M-01 二维码**：依赖选择合理（MIT、零运行时依赖、纯 JS）。otpauthURL 作为 `addData` 文本编码为模块矩阵，SVG 只输出数值坐标 rect，不把 URI 插入标记；只读 Input 经 React 转义。注入面不成立。静区见 F-005。
3. **U-01/U-02 optionsSource**：客户端同源 GET，不是服务端 SSRF。正则 `/^\/(?!\/)[^\s\#]*$/` 拒绝 scheme / `//` / `#`，允许 query。失败与非法源 fail-closed 空集成立。目录与 `permissions` / `menu_items` 同源声明属实。协议 pin **未改制品**属实；「上游无动态选项控件」与同名 string vs 协议 object 见 F-002。
4. **目录端点**：`roles.read` 恰当（D-002 最小暴露）。`admin.roles` provider 声明路由与 `profile.go` 一致，`TestRolesProviderRegistersSurfaces` 钉扎，MODULE_API_MISMATCH 防御成立。内容扩展（既有模块路由键增量），非 Profile 默认集 / 模块矩阵 / Manifest 装配变更 → 不触发 go 门闩失效，与 workspace 既有口径一致。
5. **U-04～U-07**：`mode: search` + `targetTable` + 字段 `q` 写入 `tableQuery` 并重置 page=1，清空会覆盖旧 q。行操作前 2 个内联、其余 `RowActionsMenu`，`effectivePermission` 与 `disabledWhen` 与内联同一谓词。pageSize 切换回 page 1；跳页越界忽略。空状态改为图标+文案，schema-crud 仍断言 “No items match.” 无回归。回收站缺口见 F-003。
6. **回归**：本轮切片绿；全量 1002 / 全仓 `go test ./...` 未本轮复跑，不把 E-004 全量数字当作本意见独立证据。
7. **无越界**：未改协议 pin 制品、未改 Manifest 装配、未改 Profile 默认模块集；canonical 范围在 workspace-010。

## Findings

### F-001 · I-004 到期未闭合（S4 方案门禁）

| 字段 | 值 |
|------|-----|
| level | required |
| severity | med |
| status | open |
| 影响门禁 | S5 关门（P-005：required 信息项最晚阶段已过且未 residual） |
| evidence | `00-meta.md` I-004 `open`；`01-decision.md` I-004 `open` / 「待决」；无 `01-decision/D-003-*.md` |
| 描述 | I-004（Toast 方案；搜索/筛选是否扩展协议）为 required，最晚阶段 S4。S4 已实施：Toast = FeedbackRegion 4s 可关；搜索 = 既有 `mode: search`（未扩协议）；筛选未加（后端无 filters 解析）。信息在 E-004 **事实层面已回答**，但未写 D-003、未改 I-004 状态。到达最晚阶段的开放 required 信息项阻断对应门禁。 |
| 建议 | `/govern` 写 D-003（或等价决策节）把 E-004 取舍落盘，将 I-004 标 closed；勿静默推断。 |

### F-002 · optionsSource 与协议字段同名不同形；「上游无控件」不准确

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | med |
| status | open |
| 影响门禁 | 不阻断 S5（pin 制品未改；运行时 D-VAL 不校验字段形状） |
| evidence | `docs/schemas/component-registry.json` L1858–1885（select.`optionsSource` 为 **object**：url / labelField / valueField，since 0.2）；L2176–2180 checkboxGroup「MVP 无 optionsSource」且 `options` required；本地 `form-controls.ts` L100 `optionsSource?: string` + `optionsMapping`；`users.json` L243 / `roles.json` L87 为字符串；D-002 L32「上游无动态选项控件」；D-002 L27「pin v2.8.0」vs `conformance-claim.json` `artifactVersion: 2.9.0` |
| 描述 | 「未改 pin 制品」诚实。但协议 **已有** object 形 `optionsSource`（select）；本地复用同名、改为 string + `optionsMapping`。若日后对页面跑 component-registry L2，users/roles 会因类型/缺 `options` 失败。D-002「上游无动态选项控件」对 select 不成立。pin 版本号写成 2.8.0 与现行 claim 2.9.0 不一致（本波未改 pin）。 |
| 建议 | S5 前在 D-003/响应节更正「上游已有 object 控件、本波刻意用本地 string」；后续若接协议控件，需迁移或改名，避免同名冲突。 |

### F-003 · U-04 回收站有 q 能力但未加搜索表单

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `01-decision.md` U-04 列举 recycle-bin；`recyclebin/service.go` L111 `ListItems(..., q, ...)`；`recycle-bin.json` 无 `mode: search`；E-004 七页清单不含回收站 |
| 描述 | E-004 称「全部 QSearch 后端资源页」。回收站 ListItems 已接 `q`，属 QSearch 能力页，却未加搜索表单。 |
| 建议 | S5 补 schema 或在 D-003 写明「回收站本波刻意不做」及理由。 |

### F-004 · U-03 Toast 仍占文档流

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `render.tsx` FeedbackRegion L1116–1137：`mb-4`、无 `fixed`/`absolute`；挂在 `RenderPageSurface` L2399 正文之前 |
| 描述 | 4s 自动消失 + 可关闭，相对常驻 Alert 已改善。U-03 期望「不占布局」的浮层未完全达到：展示期内仍挤压页面。 |
| 建议 | 后续改为 `fixed`/`absolute` 角标；不阻断本波。 |

### F-005 · QR 静区声明与实现不符

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `qr-code.tsx` L23 注释「4-module default margin」；L52 `viewBox={0 0 count count}` 无额外边距。库实现 `_moduleCount = typeNumber * 4 + 17`（无静区）；库自带 SVG 用 `margin * 2` 扩 viewBox（`qrcode.js` L521） |
| 描述 | 自定义渲染未加静区，仅 1px CSS border。部分扫描器在深色/贴边背景下可能失败。无 HTML 注入面。 |
| 建议 | viewBox/白底外扩 ≥4 模块；或改用库 `createSvgTag({margin:4})`。S5 可选补。 |

### F-006 · rotate 错码 400 路径测试缺口

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `mfa_test.go` `fakeMFAService.RotateRecovery` L77–78 **不校验**码即返回成功；生产 `service.go` L227–229 走 `requireActiveSecondFactor` |
| 描述 | 生产旋转错码仍会 400（`writeSelfServiceMFAError`）。handler 单测未覆盖 rotate 错码。 |
| 建议 | fake 对齐 Disable 校验，并断言 rotate 错码 400。 |

### F-007 · 目录端点缺「无 roles.read → 403」用例

| 字段 | 值 |
|------|-----|
| level | recommended |
| severity | low |
| status | open |
| evidence | `rbac_catalog_test.go`：匿名 401、viewer（默认有 roles.read）200；无无权限身份 403 |
| 描述 | 实现确有 `requirePermission(..., "roles.read")`。测试未钉扎缺权 403。 |
| 建议 | 补一个无 `roles.read` 用户的 403 用例。 |

## 必改项汇总

| ID | 级别 | 摘要 |
|----|------|------|
| F-001 | **required** | 闭合 I-004：写 D-003（记录 E-004 既成取舍）并将 I-004 标 closed |

无 security required。F-002～F-007 为 recommended，不阻断 S5。

## 与既有意见的异同

`03-audit` 此前无 A 条目。本条为该目标首条正式独立意见。

## 结论 + 建议给编排器/用户的下一步

S2 认证语义（401→400 分轨、登录 401 保留、解绑全会话吊销、错码不登出）**成立**，未见绕过。S3/S4 功能主张大体可核对；协议 pin 制品未改；go 门闩不因本波内容扩展失效。

**不可无条件放行 S5 关门**（F-001 / P-005）。建议：

1. `/govern` 响应本意见：写 D-003 闭合 I-004（F-001 → fixed）；F-002～F-007 按 fixed / accepted-residual / 留待后续波次处置并留痕。
2. I-004 合法闭合后进入 S5 关门编排（self 关门审 + 响应本 independent）。本条已覆盖 S2 security；**除非整改再改认证语义，无需为 F-001 重开 independent**。
3. 全量 Web 1002 / 全仓 `go test ./...` 建议在 S5 再跑一轮并写入 E 条目（本意见只独立复跑切片）。

## 声明

本意见不修改 status / progress / goal-tree / 方案正文 / 00-meta。响应由 `/govern` 处理。
