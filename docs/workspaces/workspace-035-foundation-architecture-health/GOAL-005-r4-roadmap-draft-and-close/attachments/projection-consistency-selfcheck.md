---
doc_type: goal-attachment
id: projection-consistency-selfcheck
parent: GOAL-005-r4-roadmap-draft-and-close
status: recorded
created: 2026-09-10
updated: 2026-09-10
version: 0.2.0
---

# 投影一致性自检（R4 关门前置）

用途：在声明任何 finding `fixed`、请求独立复核、或关闭目标之前，**先**运行可执行自检，机器化核对跨投影一致性。
背景：R4 的 required finding 全部源自同一失效模式——改了一个位置、漏改同一事实的其它投影（另有编号自指与 version 未递增）。本自检把该核对从「人肉记忆」改为「可重复执行 + 可独立复现」。

## 可执行入口（权威）

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File docs/workspaces/workspace-035-foundation-architecture-health/GOAL-005-r4-roadmap-draft-and-close/attachments/projection-selfcheck.ps1
```

- 脚本：[projection-selfcheck.ps1](projection-selfcheck.ps1)（**只读**；带 UTF-8 BOM，Windows PowerShell 5.1 依赖 BOM 解析中文，请勿去 BOM）
- 退出码：`0` = 全部通过；`1` = 存在失败项
- 覆盖检查：

| # | 检查 | 判定方式 |
|---|------|----------|
| 1 | 审计编号自指 | `03-audit/A-0NN-*.md` 正文不得把「下一次独立复核」写成自身编号 |
| 2 | 未来式语态 | 不得对已发生的响应使用「待 `/govern` 响应闭合后方可宣称」（历史引述须带「历史/原文」限定） |
| 3 | VP-035 当前版本投影 | 扫描 workspace 全部 md（**排除 `03-audit/` 台账**，其正文属历史时点）+ `docs/vision/roadmap.md` + `docs/vision/workspaces.md`；含 `VP-035` 的行中每个 `vX.Y.Z` 必须等于 VP 当前 `version`，除非该行显式标注「激活记录/历史/时点记录/2026-09-0N 激活」 |
| 4 | frontmatter 必备字段与 id | workspace 内所有 md 含 `status`/`created`/`updated`/`parent`/`version`；`00-meta.md` 的 `id` 等于目录名 |
| 5 | 边界守恒 | `git diff --name-only ebe6013c..HEAD -- apps` 为空；`docs/vision/roadmap.md` 的 `trigger-gated` 计数不下降 |

**未脚本化的检查 4′（人工）**：最近若干次提交中「内容变化即 version 递增 / `updated` 等于内容变更日」。每次提交前用 `git log --name-only -4` 列出文件并逐个核对（A-012 F-012、A-014 F-012 均由该方法发现）。

## 2026-09-10 运行记录（A-013 修正后）

命令与原始输出（仓库根）：

```text
> powershell -NoProfile -ExecutionPolicy Bypass -File docs/workspaces/workspace-035-foundation-architecture-health/GOAL-005-r4-roadmap-draft-and-close/attachments/projection-selfcheck.ps1
PASS  1 审计编号自指
PASS  2 未来式语态
      VP-035 当前 version = 0.2.1
PASS  3 VP-035 当前版本投影
PASS  4 frontmatter 必备字段与 id 一致性
      trigger-gated 基线=36 现=37
PASS  5 边界守恒
ALL CHECKS PASS
EXIT=0
```

有效性证据（非空跑）：

| 阶段 | 结果 |
|------|------|
| 首轮运行（修正前） | 检查 3 **FAIL**，命中 `r4-doc-hygiene-anchors.md:23`（VP-005 版本与 VP-035 同行造成歧义）——说明检查 3 非恒 PASS |
| 修正 `roadmap.md:365` 前（A-014 F-011 项） | 检查 3 命中 `roadmap.md:365` 的当前 `v0.2.0`（首次运行输出中亦可复现：Admin 分支行） |
| 检查 4 | 首轮运行曾捕获自检文档自身缺 frontmatter，补齐后 PASS |

## 与 A-014 F-013 的对应

A-014 指出：① 文档引用的入口不存在；② 检查 3 声称覆盖 roadmap 但 `$targets` 未含 roadmap；③ 检查 4 缺可审计的人工输出。本 v0.2.0 已逐条处置：新增真实可执行 `.ps1`（并在本文给出精确命令）、检查 3 改为扫描 workspace + roadmap + workspaces 的**全部** `VP-035` 命中并区分显式历史记录、检查 4′ 明确人工方法并给出 `git log --name-only -4` 命令、并附原始 stdout/退出码与「非恒 PASS」证据。
