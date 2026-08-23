# S1 · 当前状态扫描 · 缺陷/缺漏/漂移台账

> 本台账为 GOAL-003-s1-current-state-scan 的执行证据。按 Root [D-003 §9](../../GOAL-001-admin-module-readiness/01-decision/D-003-s0-denominator-freeze.md) 冻结严重度量尺分类；候选基线 commit `852ee7e`（clean）。每条 finding 记录严重度、类别、证据、影响门禁与关闭路径。

## 分类规则（冻结量尺）

- `blocker`（required）：阻断启动/构建/依赖闭包/Manifest/Schema/fail-closed、破坏认证授权/数据隔离/迁移完整性/协议边界、证据不可复现、全局 protocol-gap。
- `major`：影响跨模块必需能力或退出判据；涉生产边界（容器/迁移/认证授权/数据隔离/失败语义/消费路径）→ required，否则 non-blocking+延期。
- `minor`（non-blocking）：局部质量/UX/文档/可维护性。
- `info`：观察项。

---

## Required 发现（S4 整改范围）

### F-002 · 跨模块共享宿主模态缺少焦点约束/恢复与 Escape；移动抽屉无焦点管理；仓库无焦点断言
- **严重度**：major → **required**（S0 §8 映射：跨模块焦点丢失影响标准模块 → blocker/required 候选）
- **类别**：feature-gap（a11y）
- **证据**：
  - `apps/web/src/renderer/modal.tsx:13-45` — ModalHost 仅 `role="dialog"`/`aria-modal`/`aria-label`/close aria-label；无 focus trap、无 Escape 关闭、无焦点恢复
  - `apps/web/src/app/App.tsx:565-599` — 移动抽屉开/关无焦点移入/恢复
  - 全 `apps/web/src` 无 `focus()`/`Escape`/`restoreFocus` 命中（除 data-table.tsx:107 onKeyDown 停止冒泡）
- **影响门禁**：S0 冻结 §8 跨模块 UI 可访问性下限（模态焦点进入/约束/恢复成立）；证据形式为「可复跑焦点/状态断言 + 人工核对」，当前无断言支撑
- **关闭路径**：S4 实现 ModalHost 焦点 trap/Escape/恢复 + 移动抽屉焦点管理 + 对应 vitest/Playwright 断言；人工核对代表页

---

## Major 发现（deferral 到 S3）

### F-001 · `I-PROTO-FULL-001` 冻结覆盖权威声称「37/37 全绿 / 0 exclude / 320/320 全绿」，与实现（35/37 执行 + 2 exclude）及 S0 冻结（318+2）矛盾
- **严重度**：major（non-blocking for S1；deferral → S3）
- **类别**：governance-drift / test-drift
- **证据**：
  - `docs/workspaces/workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/attachments/I-PROTO-FULL-001-coverage-v2-7-0.md:30,60,80,112` — 「无 exclude / exclude 0 / 37/37 全绿 / 320/320 全绿」
  - `apps/web/src/protocol/upstream-fixtures.test.ts:549-558` — 显式排除 `m1-missing-app-manifest-capability`、`m1-navigation-without-capability`（错误信封 `MISSING_REQUIRED_CAPABILITY` vs 上游 `CAPABILITY_REQUIRED`）
  - Root `D-003-s0-denominator-freeze.md:84-85` — 正确记录「318 执行 + 2 排除」
  - 计数复核：16 文件 320 cases、registry types 24、能力域 12 —— 数字本身成立，但「0 exclude / 320 全绿」声明不可按所述复现
- **影响门禁**：协议面分母（S0 §5）；`I-PROTO-FULL-001` 是冻结覆盖权威（workspace-005，§7 不得静默改写），且被 charter.md/workspaces.md/roadmap.md 引用传播
- **关闭路径**：**fixed**（S3 I-READINESS-003）：workspace-005 v1.0.1 + D-003/E-007；workspace-008 A-003；不得改写本 finding 原始证据段与严重度，现行投影以勘误版为准。

**响应（2026-08-10）**：原 finding 的历史证据与 verdict 保留；正式勘误已落盘，开放 required = 0。

---

## Minor 发现（non-blocking，延期）

### F-003 · 迁移账本注释陈旧（README 0001-0008、kernel 注释 0001..0009、migrate_test 注释 0001-0009/0008；实际 0001-0010）
- **证据**：`README.md:96`；`apps/api/internal/kernel/persistence.go:47`；`apps/api/internal/store/migrate_test.go:110,558`（断言本身 `:124` 为 1..10）；实际台账 `modules/{authsession,settings}/migration/*.go`
- **关闭路径**：文档/注释勘误（S4 轻量修正或随 S1 延期）

### F-004 · QUICKSTART 端口漂移（8080/8081/5173 vs 实际 25080/25081/25173）
- **证据**：`QUICKSTART.md:50,51,58,61,73,78-83` vs `compose.yaml:33,54`、`vite.config.ts:18,33`、`playwright.config.ts:24`、`smoke.sh:31-32`、`config.go:44`
- **关闭路径**：QUICKSTART 勘误（fork 用户按文档端口无法找到服务）

### F-005 · `compose.yaml` `APP_PROFILE` 默认 `admin` vs app 默认 `mvp`
- **证据**：`compose.yaml:22`（`${APP_PROFILE:-admin}`）；`config.go:56`（`ProfileMVP`）；`.env.example:17`（`mvp`）
- **关闭路径**：统一默认或文档化差异；fork 用户未显式选 Profile 时两条启动路径得到不同模块面

### F-006 · 错误编目缺失 5 个实际发出的 domain 错误码 → zh-CN 下英文回退
- **证据**：缺失 `LAST_ADMIN`/`SELF_OPERATION`/`INVALID_ROLE_REF`/`ROLE_ASSIGNMENT_FORBIDDEN`/`INVALID_MENU_ITEM_REF`（`handler/users.go:253,283,285,287`、`roles.go:233-243`）；编目表 `errorcatalog/errorcatalog.go:25-73` 无此 5 码；`handler/error_contract_test.go:34-40` 正则只抓字面量
- **关闭路径**：补齐编目键（双语）；web 端对未知 key 已优雅回退，非功能性断链

### F-007 · 上传端点仅认证、无授权权限键门禁（viewer 也可上传/读取）
- **证据**：`handler/upload.go:74-78` — `POST /api/upload`/`GET /api/files/{id}` 仅 `a.Middleware`，无权限键
- **关闭路径**：S3 协议判断（D-UPLOAD 是否要求服务端权限键）；若协议要求则补授权门禁
- **closure（2026-08-23 补注，[workspace-010 GOAL-033](../../../../workspace-010-design-implementation-conformance/GOAL-033-w22-residual-closeout/00-meta.md) B5 复核）**：`accepted-residual` —— 用户 2026-08-10 go 裁决时书面确认维持 deferred、不升 required、不扩 scope（见 [GOAL-007 D-001](../../GOAL-007-s5-admission-audit-and-verdict/01-decision/D-001-s5-go-decision.md)）；owner=VP-008 lead；触发 = 每个后续业务 VP 激活前的消费前 freshness review + 后续协议判断/用户扩 scope；复核结论：触发未发生，状态有效

### F-008 · S5 分母测试使用不存在的权限键 `activity.read`（真实键 `operations.read`）
- **证据**：`apps/web/src/i18n/s5-denominator-render.test.tsx:197`；真实键 `modules/activity/provider.go:43,72`
- **关闭路径**：测试上下文改为真实权限键（当前未破坏测试，但掩盖真实键）

### F-009 · README 将 workspace-003 标为「当前」工作区，活动 lead 实为 workspace-008
- **证据**：`README.md:15`（「当前模块化架构工作区目标树」）；workspace-003 对应 VP-003 已 closed；活动交付 VP 为 VP-008 / lead workspace-008
- **关闭路径**：README 文档入口更新或补引 workspace-008

---

## Info 观察

### F-010 · D-PERM 角色授权模型不一致（users 分配用 `roles.assign` 权限键，roles 授权变更用 admin 角色字符串）
- **证据**：`handler/users.go:255`（roles.assign）；`handler/roles.go:116,127,154,170`（检查 actor.Roles 含 admin）
- **说明**：更严格而非放空；可能为 I-011-004 设计选择；S2 澄清是否统一权限键

### F-011 · S0 已登记观察（无 `.nvmrc`/`engines`，Node 22 仅 CI+Dockerfile 固定）
- **证据**：Root `D-003-s0-denominator-freeze.md:29`；GOAL-002 A-001
- **说明**：确认存在，非新发现

---

## 未发现问题的类别（扫描确认）

- API 构建/启动/依赖闭包/fail-closed：`go test/vet/build` 通过；`config.go:86-109`、`main.go:30-45` fail-closed 成立；kernel 模块图/贡献校验 fail-closed
- 数据隔离（多租户）：本基架单租户，N/A
- 迁移完整性：`store/migrate.go` ledger 前缀/checksum/快照/完整性检查完整；0001-0010 编译台账连续
- 模块贡献契约 M1–M6：users/roles/settings/activity 全部满足（见 02-execution 模块检查表）
- 页面分母：admin 12 / mvp 10 与 S0 冻结及 test-fixtures 一致
- i18n 键覆盖：manifest.title.*/nav.* 等键在 en-US/zh-CN 均存在
- 权限门禁（D-PERM 服务端）：requirePermission fail-closed；users/roles/settings 资源权限键正确

---

## 汇总

| 严重度 | 数量 | findings |
|--------|------|----------|
| blocker | 0 | — |
| major → required（S4 整改） | 1 | F-002 |
| major（deferral → S3） | 1 | F-001 |
| minor（non-blocking 延期） | 7 | F-003~F-009 |
| info | 2 | F-010, F-011 |
| **合计** | **11** | 全部已分类，无未分类项 |
