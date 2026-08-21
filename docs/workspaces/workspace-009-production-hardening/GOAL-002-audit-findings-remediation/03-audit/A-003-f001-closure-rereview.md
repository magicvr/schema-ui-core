---
id: A-003
goal: GOAL-002-audit-findings-remediation
title: A-002 F-001 闭合复审（independent）
source: independent
auditor: grok-build · grok-4.5 · high · audit skill
date: 2026-08-10
verdict: pass
parent: null
version: 0.1.0
status: active
created: 2026-08-10
updated: 2026-08-10
---

# A-003 · A-002 F-001 闭合复审（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build · grok-4.5 · high · 执行 `audit` |
| **类型** | finding-closure（复审） |
| **scope** | A-002 **F-001** required 闭合声明（fixed / commit `01b7202`）；不重审 C2–C8 / D1–D8 全盘 |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |
| **代码 HEAD（复审时）** | `d74c878`（含 F-001 修复 commit `01b7202`） |

## 范围与区间

- 只读核对：`apps/api/internal/handler/upload.go`（`activeContentMarkers` / `containsActiveContent` / upload 拒绝路径 / download 头）、`upload_test.go` 夹带用例。
- 本人执行：`go test ./internal/handler/ -run TestUpload -count=1 -v` → **PASS**（`TestUploadEndpointContract` + `TestUploadRejectsHtmlAndForcesAttachment`）。
- 本人探针：本地复现与 `containsActiveContent` + `http.DetectContentType` + `dangerousInlineTypes` 等价判定（非仅读文档）。
- 过程对照：`02-execution/E-002-a002-response.md`、`03-audit/A-002-goal-002-independent.md` F-001。

## 成果（有证据 · 本人核对）

| 主张 | 结论 | 证据 |
|------|------|------|
| 标记启发式存在且挂在上传硬门 | **成立** | `upload.go:52-67` 五标记；`upload.go:159-161`：`dangerousInlineTypes[base] \|\| containsActiveContent(body)` → 415 |
| A-002 F-001 四类失败场景现拒绝 | **成立** | 探针 reject=true：`<svg…>`→text/plain；`<?xml…><svg…>`→text/xml；600 字节垫高 + `<script>`→text/plain；`GIF89a`+`<html><script>`→image/gif |
| 回归测试锁定四场景 | **成立** | `upload_test.go:185-199` 四用例断言 415；`go test -run TestUpload` 本人复跑绿 |
| 下载面兜底仍在 | **成立** | `upload.go:208-214` attachment + CSP sandbox + nosniff（与 F-001 原缓解路径一致） |

## F-001 闭合结论

| 项 | 结论 |
|----|------|
| **原 F-001（required）** | **可 closed（fixed）** |
| **理由** | 原 finding 的可复现失败面（DetectContentType 漏 SVG / XML 前缀 SVG / 512 后夹带 script / GIF 夹带 script）已被全量 body 标记门挡住；专项测试存在且本人复跑通过。E-002「fixed」声明与代码事实一致，**非**假闭合。 |
| **残余** | 标记表仅覆盖字面量大小写两档（非 case-fold），且不含无 `<script`/`<svg` 的主动内容形态；见本意见 **N-001 / N-002**（**recommended**，不重开 F-001 required）。同源执行主边界仍依赖下载头。 |

## Findings（本复审新开 · 最多残余）

### N-001 · recommended · 中 · 标记大小写不完整 → 混合大小写 SVG/Script 可入库

- **文件:行号**：`apps/api/internal/handler/upload.go:52-58`（仅 `<svg`/`<SVG`/`<script`/`<SCRIPT`/`<?xml`）；对比 `upload.go:159-161` 硬门。
- **失败场景（本人探针）**：
  - `<Svg xmlns="…"><Script>alert(1)</Script></Svg>` → detect=`text/plain`，marker=false，**reject=false**（入库）
  - `AAA…×600` + 同上 → 同上
  - `GIF89a` + 同上 → detect=`image/gif`，**reject=false**
  - `<?XML version="1.0"?><Svg xmlns="x">` → **reject=false**
- **确认度**：高（与生产判定路径同构的本地探针）
- **说明**：浏览器对 HTML/SVG 标签名大小写不敏感；当前双字面量不等于 `bytes.EqualFold` 式覆盖。不否定 F-001 原四场景已修。

### N-002 · recommended · 中 · 无 script/svg 标记的主动内容（垫高/多格式）可入库

- **文件:行号**：同上标记表 + `upload.go:150-161`（DetectContentType 仅看约前 512 字节）。
- **失败场景（本人探针）**：
  - `AAA…×600` + `<html><img src=x onerror=alert(1)>` → text/plain，**reject=false**
  - `GIF89a` + `<html><img … onerror=…>` → image/gif，**reject=false**
  - `AAA…×600` + `<html><body onload=alert(1)>` / `javascript:` 链接 → **reject=false**
  - （对照）同内容若落在前 512 字节且被嗅探为 `text/html` 仍会经 `dangerousInlineTypes` 拒绝
- **确认度**：高
- **说明**：原 F-001 夹带用例均带 `<script>`，已被锁；本条为**相邻**残余面，不是 F-001 复开。

无新 **required**。未将 N-001/N-002 升为 required 的原因：原 F-001 主张面已闭合；下载侧 attachment/CSP/nosniff 仍在；残余为启发式完备性，非「拒绝声明完全不实」的同级缺口。

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| （无） | — | 本复审 scope 内 **无新 required** |
| N-001 | recommended | 标记改为大小写不敏感（或规范化后再匹配） |
| N-002 | recommended | 扩展主动内容启发式 / 默认 MIME allow-list / 或书面 residual「入库 best-effort，边界=下载头」 |

## 与既有意见的异同

| 点 | A-002 | A-003（本意见） |
|----|-------|-----------------|
| F-001 | open required（SVG/夹带拒绝失效） | **closed fixed**（原失败场景 + 测试已核实） |
| verdict | conditional（全 scope） | **pass**（仅 F-001 闭合复审） |
| 残余 | — | N-001/N-002 recommended（混合大小写、无 script 事件处理器垫高/GIF） |

## 结论 + 建议下一步

**verdict: pass** — A-002 **F-001** 可作为 **fixed** 合法闭合；编排器可将 F-001 标 closed，**不得**因本意见自动把 GOAL-002 标 `done`（关门仍取决于全目标 required/信息门禁/其余审计响应，本 scope 不覆盖）。

建议 `/govern`：

1. 记录 F-001 → **fixed**（引用 A-003 + commit `01b7202` + 本复审测试证据）。
2. 视风险处理 N-001/N-002（case-fold 标记、扩展启发式或 allow-list，或 accepted-residual 写清边界）。
3. 全目标无开放 required 后再做关门审计 / VP-008 go 重验证（超出本 scope）。

## 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文。响应与放行由 **`/govern`** 处理。
