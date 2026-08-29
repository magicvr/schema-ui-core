---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-zero-conflict-upgrade
version: 0.1.0
---

# E-002 · 零冲突升级演练事实（2026-08-29）

## 样本集（D-001 定案）与 V2 演进（commit `b0b41405`）

| 样本 | 内容 | 类型 |
|------|------|------|
| A 层 additive | `kernel.JoinKeys(parts ...string) string`（apps/api/kernel/keys.go） | Go 契约面 |
| 协议面 additive | `normalizePageID(id)`（apps/web/src/protocol/app-manifest.ts） | Web 协议面 |
| 新增迁移 0063 | `site_settings_updated_at_index`（admin.settings · CREATE INDEX IF NOT EXISTS · 双方言同句） | 全局台账 |
| 配套 | store 测试夹具 62→63 同步（migrate/operations/restart/identity 四处）+ 台账头常量 | 演进随附 |

changelog 迁移说明 = `attachments/changelog-upgrade-drill-v2.md`（additive 无 breaking；唯一副作用 = updated_at 索引）。

## V1 基线（commit `8686b3fd`，62 条迁移）

- golden-consumer：`kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{routes=10 pages=2 perms=3 nav=1 frag=1}`（exit 0）
- golden-web：protocol probe PASS · render probe PASS（1573 B）· token override PASS

## 下游升级执行（零冲突）

| 仓 | 升级动作 | 结果 |
|----|----------|------|
| golden-consumer | go.mod `require v0.0.1 → v0.0.2`（**1 行 diff**） | `go run` 输出与 V1 基线**逐字节一致**（63 条迁移 fresh apply · exit 0） |
| golden-web | `pnpm install`（file: 快照重拉 = 发布-安装语义） | 旧探针全过（兼容）+ `normalizePageID('  Users ') → 'users'`（V2 能力可用）· **代码 diff = 0 行** |

- **冲突计数 = 0**（gc 仅版本号 1 行；gw 0 行）；**全程无 git merge 命令**（过程事实：仅版本符号变更 + 安装重拉）。
- 主仓全量 `go test ./...` **exit 0**（含 internal/composition ok 33s——**F-003 复核：PG drain harness 全量并发本轮无失败**）。

## 挂账状态（R5 前）

- F-005（PG external 消费）：未核销（环境无本地 PG）；复审触发 = R5 发布回归。
- I-007（协议 pin 漂移 2.9 vs 2.8）：**演练证实代码面（2.9 支持窗 + additive 超集）无冲突**；`/vision` 处理 = R5 前置项。
- F-006（d.ts TS5056）：不变。