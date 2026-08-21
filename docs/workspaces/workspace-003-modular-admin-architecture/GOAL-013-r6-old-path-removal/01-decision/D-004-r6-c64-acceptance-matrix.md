---
id: D-004-r6-c64-acceptance-matrix
doc: decision-entry
goal: GOAL-013-r6-old-path-removal
source: orchestrator
date: 2026-08-06
status: accepted
---

# D-004 · R6 C6.4 终态验收与 VP 退出证据矩阵

## 决策

C6.4 使用下列八组机器可判定验收。所有 required 行均通过、VP 退出 #1～#7 均有
工作区 Q2 证据、self 与 Grok Build independent 均无开放 required finding，并经
`/govern` 响应后，R6-I004 才可从 `collecting` 改为 `verified`。

| ID | 验收面 | required 判据 | 主要命令 / 证据 | VP exit |
|----|--------|---------------|-----------------|---------|
| C64-V01 | 源码与旧路径退出 | 生产静态 Manifest、handler Schema fixtures、旧中央 Register/owner/fallback 与 Records 运行面零残留；允许治理历史原文与明确命名的 test fixture | `rg` 分域零命中；生产镜像内无静态 Manifest；Web 测试只读模块 owner schema | #2、#4、#6 |
| C64-V02 | API 完整回归 | 全量测试、vet、build 通过；Profile、图校验、Contribution、生命周期、readyz 与失败码均被动态执行 | `go test ./...`、`go vet ./...`、`go build ./...`；kernel/composition 定向 `-count=1` | #1、#2、#4、#5、#6 |
| C64-V03 | Web 与同一 build | unit、type/build、Chromium E2E 在 `mvp` 与 `admin` 均通过；同一 Web 代码与构建产物不因 Profile 改动 | `npm test`、`npm run build`、两次 `npm run test:e2e`（`APP_PROFILE=mvp/admin`） | #1、#4、#6、#7 |
| C64-V04 | 数据升级与恢复 | fresh、既有版本升级、未知/缺口/checksum fail-closed、升级前快照恢复、进程重启、Profile 降级保留数据、system-data 幂等与用户字段保护全部通过 | store/composition/server 定向测试 + `go test ./...` | #3、#5、#7 |
| C64-V05 | 双 Profile 容器 | 同一 API/Web image 依次以 `mvp`、`admin` 启动；`/readyz`、Manifest source、页面/路由集合、nginx 精确代理、disposable smoke 与容器重启持久化通过 | Compose 两 Profile 隔离 project；profile smoke；`scripts/smoke.sh --disposable` | #1、#4、#5、#7 |
| C64-V06 | custom 与失败路径 | custom + 显式模块可启动；custom 缺模块、未知/缺依赖/冲突/API 不兼容、端口占用、迁移漂移与未就绪均稳定 fail closed | config/kernel/composition/store 定向测试；错误文本只指向真实 env key | #1、#2、#3、#4、#5 |
| C64-V07 | fork 与运维 | 文档从 clean committed clone 可执行；正确说明 `APP_PROFILE`/`APP_MODULES_ENABLED`、本地 env 注入、模块贡献接入、Compose、升级恢复与 Profile 选择；缓存预热后 15 分钟口径内完成一次隔离启动和验收 | 本地 clone/fork 复现记录 + README/QUICKSTART/app README/Compose 一致性扫描 | #1、#3、#4、#6、#7 |
| C64-V08 | 证据与审计 | 每条 exit 含实现路径、动态结果、失败边界与已知限制；self + Grok independent 分号落盘，全部 required 合法闭合后再关门 | C6.4 evidence attachment、GOAL-013 A 条目、Root close-out A 条目 | #1～#7 |

## 关键验收细则

### A. Profile 与同一 build

- `mvp` Manifest 必须包含 core 示例、users、roles，不包含 settings/activity；相应
  route/schema 为 404，但 `core.operationlog` 仍启用。
- `admin` Manifest 必须在同一集合上增加 settings/activity；受保护 route 匿名请求为
  401 而不是 404，公开 Schema 为 200。
- Web Manifest 请求必须经 Vite/Nginx 精确代理到 API，并保留
  `X-Schema-UI-Manifest-Source: api`；API 与 Web 代理响应 bytes 相同。
- `custom` 只在提供 `APP_MODULES_ENABLED` 时成立；至少运行一个满足依赖闭包的
  custom 启动用例和一个缺失配置的 fail-closed 用例。

### B. 数据与恢复

- 迁移目录始终覆盖所有 compiled modules，不随 Profile 改变；验证全局顺序、checksum、
  tombstone、未知已应用迁移、缺中间版本和事务回滚。
- 验证 fresh bootstrap 与 versioned reconcile 分离；disabled Profile 数据、用户拥有字段、
  operation-log 历史在降级/恢复后保留。
- 至少包含一次真实进程或容器重启持久化；单元级快照恢复不能替代该项。

### C. 旧路径和测试 fixture

- Web 测试不得再引用已删除的
  `apps/api/internal/handler/fixtures/schema/`；core 页面直接读取
  `apps/api/internal/modules/schemarender/schema/`，业务页面读取各 owner module。
- 允许保留明确位于 test-only 目录的 Admin Manifest fixture，用于前端纯单测；禁止继续
  放在 `apps/web/public/`，以免 Vite/静态服务器形成可回退生产路径。
- API 的 Profile 聚合测试负责证明 runtime Manifest page 与 Schema contribution 一致；
  handler 测试不再反向读取 Web 静态 fixture 作为生产真相。

### D. CI 与本地证据边界

- `.github/workflows/r6-basic-matrix.yml` 应固化 API vet/build、Web unit/build、双 Profile
  browser/container matrix 和旧路径扫描；本地必须执行同等命令。
- 本地绿只证明当前 checkout；不得写成 hosted CI、合并 revision 或部署证据。关门记录要
  明确该边界，并记录实际 commit、工具版本、命令、退出码与关键断言。

## 已知 required 基线修复

2026-08-06 在方案冻结前执行 Web `npm test`：`481/495` 通过、14 失败。失败全部是
C6.3 删除中心 Schema fixtures 后，`representative-pages.integration.test.tsx`、
`representative-pages.test.tsx`、`schema-crud.test.tsx` 仍读取旧 handler fixture 路径。
这是 C64-V01/V03 的 required 修复，不是可接受 residual；修复前不得执行 C6.4 审计。

另有 required 文档/运行契约漂移：根 README 使用不存在的 `MODULES_ENABLED` 环境键；
QUICKSTART 仍指导写旧 handler fixtures、Web public 静态 Manifest 与中心 store seed；本地
`.env` 被描述为可自动加载，但 API 只读取进程环境。这些全部走 fixed，不请求 residual。

## 审计与关门顺序

1. 先完成 C64-V01～V07 并形成逐条 evidence package；失败即保持 R6-I004
   `collecting`。
2. GOAL-013 先写 self close-out，再由 Grok Build `grok-4.5` / `high` 执行独立
   `/audit`；independent 只写意见，不改 status/progress。
3. `/govern` 响应全部相关意见；无开放 required 后，才可勾 C6.4、置 R6-I004
   `verified`、关闭 GOAL-013 并同步 goal-tree。
4. Root 另做完整 close-out self + Grok independent；`progress: 6/6` 不自动推导
   `done`。只有 Root 门禁通过后才改 Root status。
5. VP-003 `closed` 属 `/vision` 决策层动作，不由 Goal close-out 自动推导。

## 理由

- VP #1～#7 横跨代码、数据、容器、fork 和生产代理，仅跑 `go test ./...` 无法覆盖终态。
- C6.3 已证明 API Contribution/lifecycle，但未重跑 Web，实际产生 14 个旧路径回归；C6.4
  必须把跨应用回归作为 required。
- 将测试 Manifest 移出 `public/`、将 core Schema fixture 指向 owner module，既保留前端
  纯单测速度，又不保留可静默回退的生产装配路径。
- 本地与 CI 共用同一矩阵，减少“文档声称可 fork、自动化只测默认 mvp”的漂移。

## 未选方案

- **只运行现有 workflow 默认路径**：Compose 和 browser 均只验证默认 `mvp`，不能证明
  `admin` 或同一 build。
- **恢复旧 handler fixtures 让 Web 变绿**：重建第二 Schema 真相源，违反 C6.3 与 exit
  #4/#6。
- **保留 `public` 静态 Manifest，仅依赖 Dockerfile 删除**：开发/其他静态部署仍可能
  回退，不能证明唯一生产路径。
- **只引用历史阶段审计**：阶段 pass 不是七条终态的本轮动态证据，也不能关闭
  R6-I004。

## 影响与后续

- 先修 Web fixture ownership、静态 Manifest test fixture、真实 env key 与 fork 文档，
  并增强 workflow/Compose Profile matrix；形成独立 Git checkpoint。
- 再执行 C64-V01～V07，记录实际结果和 evidence package；D-004 本身不改变 C6.4、
  `progress: 3/4`、R6-I004 或任何 status。
