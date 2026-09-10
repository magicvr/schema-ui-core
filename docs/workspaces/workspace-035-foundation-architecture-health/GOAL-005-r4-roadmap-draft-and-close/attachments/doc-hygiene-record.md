---
doc_type: goal-attachment
id: doc-hygiene-record
parent: GOAL-005-r4-roadmap-draft-and-close
status: recorded
created: 2026-09-10
updated: 2026-09-10
version: 0.1.0
---

# R4 文档卫生执行记录（G-001 / G-003 + G-002 / G-005 归属）

## 分工（R4 D-001 §4）

| 项 | 文件 | 执行阶段 | 状态 |
|----|------|----------|------|
| G-001 | `docs/architecture/overview.md` 现时节 | R4 C3（本记录 §1） | **已执行**（v0.11.0） |
| G-003 | `docs/architecture/cache-redis-seam-and-track.md` §2.6 | R4 C3（本记录 §2） | **已执行**（v1.2.0） |
| G-002 | `docs/vision/roadmap.md` 现状锚点 / RT-P04 / RT-D02 | 随 `/vision` editorial 同一事务 | **已执行**（VR-075 / VRev-088，提交 `2b511aa0`） |
| G-005 | `docs/vision/**` 的 mfa-wrap 旧表述 | 随 `/vision` editorial 同一事务 | **已执行**（同上；VP-016 历史原文保留 + 2026-09-10 注记） |

## 1 · G-001 · `overview.md` 现时节

**问题**（R3 C3 分类为「现在修（文档）」）：`docs/architecture/overview.md` 的「当前阶段（现时）」停在早期基线——Charter 版本 `@0.2.0`、产品树写 `web/` FastAPI、组合编排停在 VP-005 active、工作区清单只到 `workspace-006`。

**证据锚点（修正前）**：`overview.md:64`（`web/` FastAPI）、`:71`（Charter `@0.2.0`）、`:72–73`（VP-005 active）、`:74–79`（工作区清单只到 006）。

**执行**（2026-09-10，`docs/architecture/overview.md` v0.11.0）：

| 项 | 修正前 | 修正后 |
|----|--------|--------|
| 仓库布局表 | `web/` = FastAPI Web 应用 | 拆为 `apps/api/`（Go 后端）与 `apps/web/`（React/TS Admin Shell），并加历史说明注明早期 `web/` 参考应用已被取代 |
| 当前阶段标题 | `## 当前阶段（现时）`（无日期） | `## 当前阶段（现时 · 2026-09-10 经 VP-035 R4 复核修正）`，并声明本节改为**指针式**表述以防再次过期 |
| 真相源 | 「含 `workspace-001`～`004`」 | 「现行 `workspace-001`～`035`」 |
| 愿景 | Charter `@0.2.0` | Charter `@0.4.0`（仍唯一 active） |
| 组合编排 | 「VP-001～VP-004 均 closed」 | 指向 `roadmap.md`（现行 v0.80.0）；补一行「架构骨架 A0–A7 已由 VP-013/014/015/016/017/021 交付，唯一未触发项 = A3」 |
| 当前交付 VP | VP-005 `active` v0.4.1（lead workspace-006） | VP-035 `active` v0.2.0（lead workspace-035） |
| 工作区清单 | 逐个复述 `workspace-001`～`006`（含 workspace-006 `active / 0/5` 等陈旧状态） | 删除逐条复述，改为指向 [vision/workspaces.md](../vision/workspaces.md) |
| frontmatter | `updated: 2026-08-08` / `version: 0.10.0` | `updated: 2026-09-10` / `version: 0.11.0` |

**未改**：正文的架构结论（逻辑架构图、治理协议说明、模块扩展入口、演进方向）保持原样；本次只修正「现时」事实与索引指向。

## 2 · G-003 · `cache-redis-seam-and-track.md` §2.6

**问题**：§2.6.1/§2.6.2 只描述旧三元组 `Allow`/`Record`/`Clear`，未覆盖现行端口的原子面 `AllowRecord`/`Reserve`/`Cancel`（修正前该文档全文这三个名字 0 命中）。

**证据锚点**：文档 `:68`、`:73`；端口 `apps/api/kernel/ratelimit.go:40,44,52,65,71,76,79` 与 `:86`–`90`；内存实现 `internal/ratelimit/memory.go:158,174,187`。

**执行**（2026-09-10，文档 v1.2.0）：

| 位置 | 修正后内容 |
|------|-----------|
| §2.6.1 新增一条 | 列出 `kernel.RateLimiter` 现行 **7** 个方法（`Allow`/`Record`/`AllowRecord`/`Reserve`/`Cancel`/`RetryAfterSeconds`/`Clear`）与 `RateLimiterProvider` 工厂，标注 VP-032 来源与「新调用点 SHOULD 使用 `AllowRecord`」，并给出精确锚点 |
| §2.6.2 原语条目扩写 | 旧三元组保持 VP-027 冻结语义；补「原子三元组（VP-032）在同一原语上复合」，并明确 `Reserve`/`Cancel` 的 token 关联结构**触发立项时裁决**（本合同不预裁实现） |
| frontmatter | `updated: 2026-09-10` / `version: 1.2.0` |

**边界遵守**：只补现行端口已声明的方法面与语义边界；未冻结任何 Redis 侧实现细节（Lua/事务形态、ZSET vs 双桶、超时/重试参数均保留「触发立项时裁决」）。

## 3 · 边界

- 只改文档，不改任何端口、Profile、生产代码或路线图状态列。
- 不替 owner 冻结 Redis 实现细节：§2.6 只补「现行端口已声明的原子面」，把 Redis 侧具体实现留待触发后立项。
