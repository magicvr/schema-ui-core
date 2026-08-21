---
id: E-016-r6-c64-static-manifest-fixture-checkpoint
doc: execution-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: recorded
---

# E-016 · R6 C6.4 静态 Manifest 与测试 fixture 迁移 checkpoint

## 已发生事实

- 提交 `99784bc`（`test(r6): retire static manifest fixtures`）完成第一个 C6.4 实现切片：
  - 将 `apps/web/public/.well-known/schema-ui/app-manifest.json` 以 R100 rename 迁入
    `apps/web/src/test-fixtures/app-manifest.admin.json`，生产 `public/` 不再携带静态
    Manifest；
  - Web Manifest 纯单测改读明确的 test-only fixture，core Schema 测试改读
    `apps/api/internal/modules/schemarender/schema/`，业务 Schema 继续由 owner module 提供；
  - API handler Schema 测试删除对 Web 静态 Manifest 的反向依赖；Web Dockerfile 从构建后
    删除静态文件改为断言该文件不存在；
  - custom Profile 缺显式模块时的错误文本改为真实环境键
    `APP_MODULES_ENABLED`，并新增 custom + 显式模块的 Fx 启停测试。
- checkpoint 前完成动态验证：
  - `apps/api`: `go test -count=1 ./...`、`go vet ./...`、`go build ./...`，退出码均为 0；
  - `apps/web`: `npm test` 为 24 files / `495/495`，`npm run build` 成功；
  - `git diff --cached --check` 通过，提交仅包含上述 owned paths；任务开始前的
    `account_test.go`、`auth_test.go`、`health_test.go` 换行状态未被暂存。
- 非治理文档的宽搜未发现生产代码或测试仍引用旧 handler Schema fixture 或 Web public
  Manifest；仍有 `QUICKSTART.md` 与 `apps/web/README.md` 陈旧指导，纳入下一 fixed 切片。

## 状态边界

- 本条证明 C64-V01/V02/V03 的基础代码回归已恢复，不等于 C64-V01～V07 全部验收。
- README/QUICKSTART、Compose/CI 双 Profile、E2E、升级恢复、容器与 clean-fork 证据仍开放。
- R6-I004 保持 `collecting`，C6.4 不勾选，GOAL-013 保持 `active / 3/4`，Root 保持
  `active / 5/6`。

## 下一步（计划）

1. 修正 Profile/env/fork 文档，显式透传 Compose Profile，并固化可复用双 Profile smoke。
2. 增强 CI 的 API vet、旧路径扫描与 browser/container 双 Profile matrix。
3. 提交第二个实现切片后执行 D-004 C64-V01～V07 完整证据矩阵。
