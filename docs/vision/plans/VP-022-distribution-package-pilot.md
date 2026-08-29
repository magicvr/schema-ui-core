---
doc_type: vision-plan
id: VP-022-distribution-package-pilot
title: 分发形态试点 · 构建期包消费路径（对标 dotnet new + NuGet / Spring Boot starters）
status: closed
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-022-distribution-package-pilot
created: 2026-08-29
updated: 2026-08-29
version: 0.4.0
---

# VP-022 · 分发形态试点 · 构建期包消费路径

## 意图

在既有「fork 消费」路径之外，把「**构建期包消费**」做成可验证的最小闭环：本仓（单主线）以 Go 库模块 + npm 包组形态发布 kernel / 标准模块 / Renderer / Shell，下游项目以 `go get` / `pnpm add` 组装自己的组合根与骨架应用；本仓升级时下游仅 bump 版本 + 执行 changelog 迁移说明，**全程无 git merge**，始终维持上游基准。对标 .NET `dotnet new` + NuGet 与 Spring Boot starters 的包体系。

**试点性质**：不改 Charter（成功边界 #1「可 fork」措辞不动）；不废弃 fork 路径（fork 保留为深度定制逃生舱与既有消费者默认路径）；试点只产生证据 + go/no-go 报告，是否把包消费写入 Charter 成功边界留给试点结论再议。

**边界（不进本波）**：G1 单模块粗粒度起步（不拆多模块细版本 / 独立 tag）；CLI 只做方案评估不要求交付；不做运行时基线镜像与运行时模块下载；不修订 VP-004 贡献 playbook（留到 go 后）；不迁移既有 fork 消费者。

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **Go 库消费闭环**：一个空下游仓仅通过 `go get github.com/magicvr/schema-ui-core/apps/api@<tag>` + 自建组合根，可装配 kernel + ≥1 个标准模块，达到与 fork 形态等价的功能基线（启动、Profile、双方言迁移台账、既有测试基线可复用）。
2. **Web 包消费闭环**：空下游 app 仅通过 npm 包组（protocol / renderer / shell / ui）组装，可渲染与主线同一的 schema 页面集；品牌定制仍走 Token 覆盖路径（theme/brand.css），不要求改包源码。
3. **零冲突升级演练 PASS**：上游施加一次真实演进（至少含配置键变更 + 新增迁移 + 依赖更新），下游仅 bump 版本 + 执行 changelog 迁移说明，功能基线回归全绿；全程无 git merge 操作，冲突计数 = 0。
4. **契约冻结面落盘**：kernel 公共 API / 模块契约 / npm 包组的 semver 与 breaking 流程（含 changelog 模板）成文，明确「契约冻结面」与「主线内部自由演进面」的分界。
5. **发布可复现**：脚本/CI 一键产出 Go 版本 tag + npm 包组（复用 `scripts/pre-release-smoke.sh` 思路），golden consumer（示例下游）消费回归通过。
6. **go/no-go 报告**：包路径 vs fork 路径的实测对比（升级耗时、冲突数、契约迁移成本、契约稳定性税），并给出「是否推进 Charter strategic 修订（把构建期包消费写入成功边界）」的建议。

## 激活前置要求

（响应 VRev-048 · V-F084 / V-F085，2026-08-29）

1. **freshness 复核类型**：本 VP 属组合层平台波，激活时适用**平台/架构类 freshness 轻量复核**（类比 VP-013/VP-021 先例），不适用 Admin 业务类完整矩阵；在 `workspace-022` 的 `D-001-workspace-root-establishment.md` 按先例留痕（`consumer_vp` / `last_freshness_review_at` / `next_freshness_review_trigger`）。
2. **go/no-go 触发框架**（关门判定用）：
   - **倾向推进 Charter strategic 修订**：判据 #1～#5 全部达成，且零冲突演练冲突计数 = 0、包路径升级总耗时 ≤ 同演进集的 fork 同步对照基线、golden consumer 回归稳定。
   - **倾向搁置**：任一判据未达成；或暴露无法经版本化解决的契约障碍（如 Go 公共 API 无法冻结、前端 peer 依赖成环/不可解、跨版本迁移台账不可合并）；或 golden consumer 回归不稳定。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned` 0 区；激活时开 `workspace-022-distribution-package-pilot` 并填表 |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-29 | closed | 六条退出判据按**有界试点口径**满足（用户 P-004 书面接受：F-001/F-002/F-003 residual 范围）；independent 双审闭合（grok A-002/Root A-001 响应节）；Charter 0.3.0 strategic 落地（VR-050 · VRev-050 pass） | workspace-022 Root done 5/5（GOAL-002~006 A/E 条目）；go/no-go 报告；grok 审计响应 | go 后清单：origin tag+Go proxy 发布 / 配置键+依赖样本补测 / CI+registry 上传 / 六包细化+d.ts 链路 / PG external 实测 / fork 对照计时 |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-29 | 初创 v0.1.0（用户 2026-08-29 选择「试点先行」；组合层平台波，与三分支正交、与 VP-009/010 正交；对标 dotnet new + NuGet / Spring Boot starters） |
| 2026-08-29 | v0.2.0 · 响应 [VRev-048](reviews/VRev-048-vp022-distribution-package-pilot-intent.md) V-F084/085/086 同批**fixed**：激活前置「freshness 类型 = 平台/架构类轻量复核」+「go/no-go 触发框架」+ `lead_workspace: —` 占位符 |
| 2026-08-29 | v0.3.0 · **激活**（用户指令）：[VRev-049](reviews/VRev-049-vp022-activation.md) self `pass`（0 required；架构类轻量 freshness PASS `fddaf638`→`5c168070` 不暂挂 `go`）；`planned → active`；lead `workspace-022-distribution-package-pilot` 开区（Root 纲领路线图 R1～R5） |
| 2026-08-29 | v0.4.0 · **关闭**（用户 P-004 裁决后）：Root `GOAL-001-distribution-package-pilot` `done 5/5`；独立关门审计（grok A-002 + Root A-001）required 全闭合；Charter 0.3.0 strategic 已随 GO 裁决定稿；go 后清单落款 |