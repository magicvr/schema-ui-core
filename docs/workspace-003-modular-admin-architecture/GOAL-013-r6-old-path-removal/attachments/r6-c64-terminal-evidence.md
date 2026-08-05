---
id: r6-c64-terminal-evidence
title: R6 C6.4 终态验收与 VP 退出证据
status: recorded
created: 2026-08-06
updated: 2026-08-06
parent: GOAL-013-r6-old-path-removal
version: 0.2.0
---

# R6 C6.4 终态验收与 VP 退出证据

## 候选身份与边界

| 项 | 值 |
|----|----|
| candidate revision | `9409b7176a5a07e60b9b07e3f2e1a2fc07ebf683` |
| 实现提交 | `99784bc`、`88a3840`、`9409b71` |
| Go | `go1.26.0 windows/amd64` |
| Node / npm | `v22.17.0` / `10.9.2` |
| Playwright | `1.62.1` |
| Docker / Compose | `29.6.2` / `v5.3.1` |
| Bash | Git for Windows `5.2.37` |
| 证据性质 | 本地 Windows checkout + Linux containers；不是 Hosted CI、合并、部署或发布证据 |

主 checkout 在任务开始前已有三个仅换行状态的用户改动：
`apps/api/internal/handler/account_test.go`、`auth_test.go`、`health_test.go`。它们未被暂存
或提交；clean-clone 证据用于隔离候选 revision，避免将 dirty checkout 当作终态证明。

## C64-V01～V08

| ID | 结果 | 实现 / 动态证据 | 失败边界与限制 |
|----|------|-----------------|----------------|
| C64-V01 | pass | 生产源码扫描未发现 `MountProviderRoutes`、`RegisterSettings`、`RegisterActivity`、`staticSchemaDocuments`、`schemaDocumentsForPlan`、`compiledMigrations`、`seedRBAC`；active Records runtime 扫描零命中；旧 handler Schema 与 Web public Manifest 路径均不存在。Web Docker image 构建时断言 `dist/.well-known/.../app-manifest.json` 不存在，测试只读 `src/test-fixtures/` 与 owner module Schema。 | `internal/manifest/app-manifest.json` 是 `core.manifest-route` 的 API 聚合基片，不是 Web 静态发布/兜底；`operations_test.go` 的 `/api/records` 仅验证所有方法 404。 |
| C64-V02 | pass | clean clone 执行 `go test -count=1 ./...`、`go vet ./...`、`go build ./...`；另执行 `go test -count=1 ./internal/store ./internal/migration ./internal/modules/authsession/systemdata ./cmd/server` 与 `go test -count=1 ./internal/config ./internal/kernel ./internal/composition ./internal/store`，均退出 0。 | 本地与 Linux container 证据，不声称 GitHub Actions。 |
| C64-V03 | pass | clean clone `npm ci` 后 `npm test` 为 24 files / `495/495`，`npm run build` 成功；`APP_PROFILE=mvp` 与 `admin` 的 Chromium E2E 各 `2/2`，动态断言 Manifest source、页面集合、Schema 200/404、受保护路由 401/404。 | 两次 E2E 使用同一 Web checkout；Profile 只改变 API 运行时集合。 |
| C64-V04 | pass | store/migration/system-data/server tests 覆盖 fresh、R2 升级、pre-v0002 快照恢复、未知/缺口/checksum/partial-baseline fail-closed、真实 OS server 子进程重启、reconcile 幂等、用户字段保护与 disabled Profile 数据保留。候选 image 同卷 `admin → mvp → admin` 后保留用户 `usr-dde842d69705ea60`、`siteTitle=C64 Final Retained Admin` 与 `users.create`/`settings.update` 历史。 | 快照恢复由 Go 动态测试证明；容器回环证明实际 volume/进程重建，不替代灾备演练或远端备份。 |
| C64-V05 | pass | clean clone 构建 API image `sha256:75b987...a013`、Web image `sha256:3b89f8...97bc`；同一对 image 在 `c64finalmvp9409` / `c64finaladmin9409` 两隔离 project 启动，两个 `scripts/smoke.sh --disposable` 均 SM-001～007 全绿，包含 `/readyz`、Manifest bytes/source、Profile routes、seed 与 API restart。 | 所有测试 project/volume 已在取证后删除；镜像 cache 保留。 |
| C64-V06 | pass | config/kernel/composition/store 定向测试覆盖 custom + 显式模块启停、custom 缺配置、未知模块/依赖/环/冲突/capability、Kernel API 不兼容、端口占用、provider mismatch、readyz 503/200、迁移漂移。真实 `APP_PROFILE=custom` 且空 `APP_MODULES_ENABLED` 的 `go run ./cmd/server` 退出 1，并输出 `PROFILE_MODULES_REQUIRED` 与真实 env key。 | custom 成功用例由 Fx `NewApp` 启停测试证明；未声称运行时插件或热插拔。 |
| C64-V07 | pass | 从主仓创建 `--local --no-hardlinks` clean clone，固定 `9409b71`；clone 内完成静态门禁、API full/vet/build、全新 `npm ci`、Web unit/build、双 Profile E2E、admin Compose build/start/disposable smoke；开始与结束 `git status --porcelain` 均为空，总耗时 3.56 分钟。README/QUICKSTART/app README/Compose 均使用真实 Profile/env/owner/升级恢复边界。 | 临时 clone 已送入回收站；时间口径为本机缓存已预热，不代表冷网络 SLA。 |
| C64-V08 | pass | 本附件覆盖每条 exit 的实现、动态结果、失败边界与限制；A-012 self 与 A-013 Grok independent 均 `pass`、required 0，A-014 `/govern` response 已按 fixed 闭合 F-R6-001。 | R6-I004/C6.4/GOAL-013 状态只由 A-014 响应更新；不自动改变 Root 或 VP-003。 |

## 关键命令与结果

```text
workflow retirement scans                           PASS
bash -n scripts/smoke.sh / workflow YAML parse      PASS
docker compose config (mvp, admin)                  PASS / PASS
go test -count=1 ./... / go vet / go build          PASS / PASS / PASS
npm test / npm run build                            495/495 / PASS
Playwright Chromium mvp / admin                     2/2 / 2/2
candidate container smoke mvp / admin               SM-001..007 PASS / PASS
candidate admin -> mvp -> admin data loop            PASS
clean clone full reproduction                        PASS, 3.56 minutes
```

两个诊断失败均已保留真实解释：

1. 初版 workflow Records 扫描把 `operations_test.go` 的 404 负向测试误报为运行残留；
   `9409b71` 将 active-runtime 扫描限定为非测试代码，并增加退休装配符号扫描，最终候选
   与 clean clone 原样 Bash 门禁均通过。
2. 已写入额外用户的 profile-loop project 再运行 seed-only disposable smoke 时，SM-006
   正确报告 `total=2`；该次不计 pass。双 Profile disposable 全绿证据来自 fresh 隔离
   project，数据回环另以直接 API 断言证明。

## VP-003 退出 #1～#7 映射

| Exit | 工作区 Q2 证据 | 结论 |
|------|---------------|------|
| #1 单主线与 Profile | `docs/workspace-003-modular-admin-architecture/GOAL-013-r6-old-path-removal/attachments/r6-c64-terminal-evidence.md#c64-v01v08` C64-V03/V05/V06/V07 | 同一候选、同一 image/前端 build 支持 mvp/admin；custom 成败边界可核对。 |
| #2 薄内核、组合根与模块契约 | 同附件 C64-V01/V02/V06；GOAL-013 E-006～E-014、A-004～A-011 | 旧装配符号退出；Contribution、图、生命周期、API range 与 ownership 动态通过。 |
| #3 数据生命周期 | 同附件 C64-V04/V06 | 迁移、快照恢复、reconcile、漂移、真实重启与 Profile 回环取证完成。 |
| #4 后端聚合唯一生产路径 | 同附件 C64-V01/V03/V05 | API/Web Manifest bytes 与 source 相同；静态 Web Manifest 不存在；Profile 页面/路由一致。 |
| #5 安全、横切与生命周期 | 同附件 C64-V02/V04/V05/V06 | auth/permission、operation-log、readyz、Start/Ready/Stop/failure 路径由 tests + containers 覆盖。 |
| #6 能力迁移与旧路径退出 | 同附件 C64-V01/V03/V05 | users/roles/settings/activity 与 Schema 页面可用；Records runtime、旧装配与生产静态兜底退出。 |
| #7 可 fork、运维与回归 | 同附件 C64-V03/V04/V05/V07 | clean clone、文档、双 Profile、升级恢复、容器与自动化矩阵取证完成。 |

## 关门边界

- 以上证明 C64-V01～V07 与 VP exit #1～#7 的本地候选证据齐全。
- C64-V08 已由 A-012 self + A-013 Grok independent + A-014 `/govern` response 完成；
  R6-I004 `verified`、C6.4 完成、GOAL-013 `done / 4/4`。
- Root close-out 仍需另做 self + Grok independent；Root `progress: 6/6` 不自动推导
  `status: done`。VP-003 是否 `closed` 属 `/vision`，不在本次 Root 关门范围内自动执行。
