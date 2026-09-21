---
doc_type: goal-attachment
id: r4-doc-hygiene-anchors
parent: GOAL-004-r3-industry-comparison
status: staged
created: 2026-09-09
updated: 2026-09-10
version: 0.3.0
---

# R4 文档卫生取证草稿

本文件只登记 **R4 文档卫生**（VP-035 首波表「文档卫生」行 + R3 D-001 裁决 A「R4 仅做文档卫生」）所需的现时锚点。R3 不据此改任何 `docs/architecture/**` 或 `docs/vision/roadmap.md`；R4 执行时须逐条对照当时 HEAD 复核。

代码基线：`52300975`（本次取证工作树干净）。文档侧行号指当前工作树。

## G-001 · `overview.md` 现时节过期

| # | 文档锚点 | 现在写的 | as-built / 现状 | 备注 |
|---|----------|----------|------------------|------|
| 1 | `docs/architecture/overview.md:71` | Charter **`schema-ui-core-admin-foundation@0.2.0`** | Charter 已到 `@0.4.0`（`workspace.md` / VP-035 `vision_ref`） | 需按 `docs/vision/charter.md` 重取版本 |
| 2 | `docs/architecture/overview.md:64` | `web/` = FastAPI Web 应用 | 产品树为 `apps/api`（Go）+ `apps/web`（React/TS），见 `docs/architecture/directory-layout.md` | 「本仓 web/ 冻结参考」表述与实际产品目录不符 |
| 3 | `docs/architecture/overview.md:72–73` | 组合编排停在 VP-001～004 closed；且把 VP-005 写成「当前交付 VP」 | 现行架构分支已到 VP-035（见下表第 6 行）；VP-005 早已不是「当前交付 VP」 | 这一段是「当前阶段（现时）」的核心过期源 |
| 6 | `docs/architecture/overview.md:73` | 把 VP-005 标为当前交付 VP | 现行当前交付 VP 是架构分支的 VP-035 计划 | 与第 3 行同源；VP-005 自己的版本号与本节无关 |
| 4 | `docs/architecture/overview.md:74–79` | 工作区清单只到 `workspace-006`（且 workspace-006 Root 写 `active / 0/5`） | 实际已有 `workspace-001`～`035` | 至少需改为指针式表述，避免逐个复述导致再次过期 |
| 5 | `docs/architecture/overview.md` frontmatter | `updated: 2026-08-08` / `version: 0.10.0` | 2026-09-09 仍为 0.10.0 | 卫生项：更新 `updated`/`version` |

未核对：`docs/architecture/module-architecture.md:9` 的 `vision_ref: schema-ui-core-admin-foundation@0.2.0`——A-002 已判「旧架构记录，不覆盖现行链」。R4 需决定是「保留历史语境并加注」还是「刷新 vision_ref」，两者都在文档卫生范围内，但**不得**据此改写架构结论。

## G-002 · roadmap 现状锚点 / RT-P04 / RT-D02

| # | 文档锚点 | 现在写的 | as-built | 备注 |
|---|----------|----------|----------|------|
| 1 | `docs/vision/roadmap.md:104` | 现状锚点：单进程 + SQLite（`MaxOpenConns=1`）+ 本地盘上传 + 进程内 Job + 内存限流 | SQLite 已用小连接池：`apps/api/internal/store/store.go:29` `sqlitePoolDefault = 4`；内存库仍 1（`:104`–`114` 的 memory 分支） | 需把「`MaxOpenConns=1`」改为「小连接池（默认 4）/内存库 1」 |
| 2 | `docs/vision/roadmap.md:138` | RT-P04 现状锚点 = `MaxOpenConns=1`，状态 **trigger-gated** | 同上：池化已交付；**读写分离 / replica 仍 gated** | 现状锚点过期，但 gated 判定本身仍成立——不得把「已池化」扩写成读写分离/replica 已交付 |
| 3 | `docs/vision/roadmap.md:209` | RT-D02 同行写「进程生命周期有，无明确 drain 合同」+ 状态 **delivered** | VP-021 已交付停机顺序/HTTP drain/Job 语义/双方言排空 | 同一行内「无明确 drain 合同」与 `delivered` 并存，属自相矛盾表述 |
| 4 | `docs/vision/roadmap.md:39,302` | VP-021 已 `closed` v0.3.0，RT-D02 → delivered | 一致 | 无需修改；作为 1/3 项的对账依据 |

## G-003 · cache-redis-seam §2.6 与现行 RateLimiter 端口

| # | 文档锚点 | 现在写的 | 端口 as-built | 备注 |
|---|----------|----------|---------------|------|
| 1 | `docs/architecture/cache-redis-seam-and-track.md:68` | 端口语义：`Allow` 不注册、失败才 `Record`、`Clear` 清桶、`RetryAfterSeconds` 分母 | 现行 `kernel.RateLimiter` 声明 7 个方法：`apps/api/kernel/ratelimit.go:40` `Allow`、`:44` `Record`、`:52` `AllowRecord`、`:65` `Reserve`、`:71` `Cancel`、`:76` `RetryAfterSeconds`、`:79` `Clear` | 文档未提 `AllowRecord`/`Reserve`/`Cancel` 原子三元组，接缝描述不足以指导现行替换（A-002 已复核同一事实） |
| 2 | `docs/architecture/cache-redis-seam-and-track.md:73` | 原子窗口原语只写 `Record = INCR + 首次 EXPIRE`、`Allow = GET`、`Clear = DEL` | 现行还有预留/取消语义：`Reserve`（`kernel/ratelimit.go:65`；实现 `internal/ratelimit/memory.go:174`）、`Cancel`（`:71`；实现 `:187`） | 建议「现在修（文档）」；**不得**替 owner 冻结 Redis 实现细节 |

## 边界声明

- 本文件是 R4 输入，不是 finding、不是审计意见、不关闭任何门禁。
- R4 修改 `docs/architecture/overview.md` 与 `docs/vision/roadmap.md` 属判据 4 与首波表「文档卫生」；`docs/vision/**` 的改动须与路线图草案同一事务或紧随 `/vision` 冻结（VP-035 首波表原文）。
- G-004（`assembly.NewAuthenticator` 公共工厂限制）不在本文件；其分类需在 R3 C3 由用户裁决。
