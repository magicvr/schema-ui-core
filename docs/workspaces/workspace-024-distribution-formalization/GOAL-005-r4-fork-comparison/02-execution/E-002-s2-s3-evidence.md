---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-fork-comparison
version: 0.1.0
---

# E-002 · S2/S3 双模型实测与对比报告（2026-08-29）

## 实验设置（D-001）

- 演进集：`apps/api/v0.3.0`（`4f7cb0f1`）→ `apps/api/v0.4.0`（`00d97b5b`）真实演进（serve 面 + CLI + 模板）。
- 环境：本地 Windows · Go 1.26 · git；暖缓存口径（Go 模块缓存预热）；相对对比为结论依据。
- 全部实验仓为临时构建（`%TEMP%`），实验后清理，不影响主仓/golden-field。

## fork 模型实测（S2）

fork-sim 仓 = schema-ui-core 克隆，基线 `apps/api/v0.3.0` + 2 个本地定制 commit（① `templates/main.go.tmpl` 末尾定制行 ② `QUICKSTART.md` fork 说明）→ `git merge apps/api/v0.4.0`：

| 阶段 | 结果 |
|------|------|
| merge 执行 | 0.3s（本地对象仓） |
| 冲突 | **1 个内容冲突**：`apps/api/cmd/schema-ui/templates/main.go.tmpl`（fork 定制行 + 上游整文件重写）；`QUICKSTART.md` auto-merge（行尾追加场景） |
| 解冲突 | 按迁移说明手工处置：模板取上游 v0.4.0 形态（fork 定制点丢失/需重做）+ QUICKSTART 合并保留 |
| 迁移改写点 | **2**（模板形态重适配 · QUICKSTART fork 段合并） |
| 构建验证 | `go build ./...`（apps/api）exit 0 · **12.9s**（含解冲突后再构建） |

## 包模型实测（S3）

golden-field worktree 检 `9510023`（v0.3.0 完整态）→ 按 **v0.4.0 迁移说明**只改 3 文件：

| 阶段 | 结果 |
|------|------|
| 迁移编辑 | go.mod bump 1 行 · cmd/server 换薄封装（server.Serve 面）· web/.npmrc 钉 npmjs —— **0.0s**（3 个文件，无 diff 上游） |
| 冲突 | **0**（无 git merge · 无 replace/file:） |
| tidy + build | **4.8s**（暖缓存 · 公共 proxy） |
| serve 冒烟 | `/healthz` 200 · readyz 通过 |
| 迁移改写点 | **0**（迁移说明 3 条目照做即完成） |

## 对比结论（判定 #4 · 定量实证）

| 维度 | fork 模型 | 包模型 | 差 |
|------|-----------|--------|-----|
| 冲突计数 | 1（内容冲突） | 0 | 包路径零冲突（演进集含整文件重写类变更时 fork 必冲突） |
| 迁移改写点 | 2（含定制点需重做） | 0 | 包路径执行迁移说明即可，不触碰定制面 |
| 升级主诉耗时 | 13.2s（0.3+12.9） | 4.8s | 包路径约 1/3（本样本；暖缓存口径） |
| 契约同步语义 | 合并 → diff → 适配（每波演进重复） | bump + changelog（每波固定两步） | 包路径成本恒定性 |

- **v0.3.0 定性映射 → v0.4.0 定量实证**：VP-022 go/no-go §2 的「fork 需 4 处手工同步」在真实演进集上复现为同类冲突/改写点；包路径零冲突、耗时更低、迁移成本恒定（= changelog 条目数）。
- 边界说明：单样本（v0.3.0→v0.4.0）· 暖缓存 · fork 模拟定制 2 点（真实 fork 定制面可能更多）；绝对值不跨机推广，相对关系成立。

## 核销映射

| 项 | 结论 |
|----|------|
| VP-022 判据 #6「包 vs fork 实测对比」遗留半项 | ✅ 核销（本报告 + 旧 go/no-go §2 定性补定量） |
| go 后清单「fork 对照计时实验」 | ✅ 核销 |
| VP-024 判据 #4 | ✅ 满足（E-002 + 本报告） |