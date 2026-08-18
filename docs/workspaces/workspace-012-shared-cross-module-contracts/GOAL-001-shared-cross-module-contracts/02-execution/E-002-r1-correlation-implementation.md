---
id: E-002-r1-correlation-implementation
doc: execution-entry
status: recorded
created: 2026-08-18
updated: 2026-08-18
parent: GOAL-001-shared-cross-module-contracts
version: 0.1.0
---

# E-002 · R1 correlation / request-id 实现切片与全量验证

## 已发生事实

- R1 子目标已完成方案冻结（`D-001`）并把 `I-001` required 信息标为 `verified`。
- API 新增 `apps/api/internal/requestid`：校验/生成 `X-Request-ID`、写入 context、对每个生产响应回显；handler/auth 错误包络追加 `correlation_id`。
- operationlog 新增迁移 `operation_log_correlation`（version 41），保持既有业务 detail JSON 形状；auth login/refresh/logout 与 settings patch/reset 的真实写路径记录 correlation id。
- Web `ResourceApiError` 从错误体或 `X-Request-ID` 读取 correlation id，并把它放入错误对象和用户可见消息；新增定向测试覆盖 body/header 两条路径。
- 定向验证通过：`go test ./internal/requestid ./internal/server ./internal/modules/operationlog ./internal/handler ./internal/store`；`npm --prefix apps/web test -- --run src/renderer/error-localization.test.tsx`。
- 全量验证通过：在 `apps/api` 执行 `go test ./...`；在 `apps/web` 执行 `npm test -- --run`（72 files / 1069 tests）；执行 `npm run build` 通过。
- 第一个实现 checkpoint 已创建：`e1f211f`，scope 为 R1 API/Web/operationlog 与治理记录。

## 成功标准对照（事实）

| R1 标准 | 状态 | 证据 |
|---|---|---|
| 每个 API 请求响应头有稳定 correlation id | 已验证 | requestid middleware + server/route 定向测试 + API 全量测试 |
| 错误体含 id，前端可展示 | 已验证 | handler/auth 包络、Web ResourceApiError、error-localization 定向测试 |
| 至少一条 operationlog 事件可关联 id | 已验证 | operation_log_correlation version 41、auth integration test、repository test |
| 有测试/验证路径 | 已验证 | 定向与全量命令均通过 |

## 阻塞

无。R1 仍需阶段自审/关门审计，Root R2 尚未放行。

## 下一步（计划）

追加 R1 关门前执行证据与 self close-out 审计；确认无 required finding 后关闭 GOAL-002，并同步 workspace goal-tree/R1 路线标记。
