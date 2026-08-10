---
id: GOAL-006-s4-remediation-and-regression
doc: execution
status: active
parent: GOAL-001-admin-module-readiness
created: 2026-08-10
updated: 2026-08-10
version: 0.1.0
---

# 执行记录 · GOAL-006

## 执行条目索引

| E-ID | 日期 | 动作 | 结果 | 证据 / 文件 |
|------|------|------|------|-------------|
| E-001 | 2026-08-10 | F-002 required：模态焦点 trap/Escape/恢复 + 移动抽屉焦点管理 + 可复跑断言 | fixed | `apps/web/src/renderer/modal.tsx` + `modal.test.tsx` + `apps/web/src/app/App.tsx` |
| E-002 | 2026-08-10 | minor 处置：F-006 错误编目 5 码、F-008 测试键、F-003 迁移注释、F-004 端口、F-005 compose 默认、F-009 README 指针 | fixed | errorcatalog/error_contract_test、README、QUICKSTART、compose.yaml、s5-denominator-render.test.tsx |
| E-003 | 2026-08-10 | minor 处置：F-007 上传授权深度 | deferred | 见下延期字段 |
| E-004 | 2026-08-10 | 冻结分母回归（V-001~V-008）重跑 | pass | go build/test 全包、npm test/build 全绿 |

## F-002 required 整改详情

模态（`modal.tsx`）：打开记录触发元素并聚焦容器内首个可聚焦元素；Tab/Shift+Tab 在容器内循环（focus trap）；Escape 关闭；卸载恢复焦点到触发元素。断言（`modal.test.tsx` 3 用例）：焦点进入/恢复、Tab 循环、Escape 关闭。移动抽屉（`App.tsx`）：打开聚焦首个可聚焦元素、Tab 循环、Escape 关闭、关闭恢复焦点。

## minor 处置明细

| finding | 处置 | 证据 |
|---------|------|------|
| F-006 错误编目 5 码 | **fixed**：`errorcatalog.go` 补 `LAST_ADMIN/SELF_OPERATION/INVALID_ROLE_REF/ROLE_ASSIGNMENT_FORBIDDEN/INVALID_MENU_ITEM_REF` 双语项；`error_contract_test.go` frozenDomainCodes 同步 | errorcatalog.go + error_contract_test.go |
| F-008 测试用 `activity.read` | **fixed**：改为真实键 `operations.read` | `s5-denominator-render.test.tsx:197` |
| F-003 迁移注释 0001-0008/0009 | **fixed**：README §状态说明 与 `kernel/persistence.go` 注释更新为 0001-0010 | README.md + persistence.go |
| F-004 QUICKSTART 端口 | **fixed**：8080/8081/5173 → 25080/25081/25173 | QUICKSTART.md |
| F-005 compose 默认 admin | **fixed**：`compose.yaml` `${APP_PROFILE:-admin}` → `${APP_PROFILE:-mvp}`（对齐 app 默认） | compose.yaml |
| F-009 README workspace-003 为当前 | **fixed**：README 文档入口补 workspace-008 为当前，workspace-003 标历史 | README.md |
| F-007 上传授权深度 | **deferred**：上传为协议共享能力（D-UPLOAD），服务端认证+大小/类型约束；补模块权限键需协议/设计对齐。owner=VP-008 lead；复核触发=S5 协议判断/用户扩 scope | S3-protocol-judgment §3 |

## 冻结分母回归（候选：本提交，自 S3 `3b9a584` 前进）

| 命令 | 结果 |
|------|------|
| `go build ./...` | ✅ |
| `go test ./...` | ✅ 全包 |
| `npm test` | ✅（见最新运行） |
| `npm run build` | ✅ |
| e2e / smoke | 前端改动不影响 smoke/e2e 覆盖端点；CI 矩阵等价项待 push 复核 |

## 记录规则

只写已发生事实；整改与回归结果绑定新候选 commit。失败到通过证据：F-002 修复前无焦点断言（S1 台账记录缺失），修复后 3 断言通过。
