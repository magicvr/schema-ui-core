---
id: A-002
goal: GOAL-002-audit-findings-remediation
title: GOAL-002 修复交叉独立审计
source: independent
auditor: grok-build · grok-4.5 · high · audit skill
date: 2026-08-10
verdict: conditional
---

# A-002 · GOAL-002 修复交叉独立审计（independent）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build · grok-4.5 · high · 执行 `audit` |
| **类型** | execution-facts / finding-closure 复核 |
| **scope** | GOAL-002 全部 16 项修复（C1–C8 + D1–D8）+ 回归证据；候选 `9c1d0a7`（HEAD） |
| **verdict** | **conditional** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical 已校验；`shared_materials_catalog: none`） |

## 范围与区间

- 代码：`apps/api`（upload/auth/config/composition/handler/settings/migrate 等）+ `apps/web`（auth-client、render、App、form-controls、theme、request-construction）。
- 过程：`02-execution/E-001-remediation.md`、`03-audit/A-001-goal-002-self.md`。
- 审查输入 `raw/audit-20260810-api-web-bug-review.md` 为 gitignored；本意见以当前 HEAD 源码与可复现测试为准，不依赖 raw 原文复述。
- 信息项：I-001 verified（覆盖映射有 E-001）；I-002 D3 accepted-residual（用户裁决留痕，独立审计不推翻 residual 接受本身）。

## 成果（有证据 · 本人核对）

| 项 | 结论 | 证据 |
|----|------|------|
| C1 HTML 拒绝 + 下载加固 | **部分有效** | `upload.go` DetectContentType + `dangerousInlineTypes`；下载 `Content-Disposition: attachment` + `CSP: sandbox` + `nosniff`；`TestUploadRejectsHtmlAndForcesAttachment` 通过。**SVG 拒绝声明不成立**（见 F-001） |
| C2 刷新原子化 | **有效** | `RevokeRefreshToken` 单语句 `WHERE revoked_at IS NULL` + `RowsAffected`→`ErrAlreadyRevoked`；`Refresh` fail-closed；`Logout` 容忍并发；`TestRefreshConcurrentRotationSingleWinner` 5 连跑通过；前端 `inflightRefresh` 单页内去重 |
| C3 APP_ENV fail-closed | **有效** | `Load` 默认 `""`；`ValidateProd` 空 env 拒绝；`main` 先 Validate 再 `resolveJWTSecret`/`resolveSeedHash`；非 `development` 不注入 dev secret；`TestValidateProd` unset 子测通过 |
| C4 Bootstrap 可重试 | **有效** | `NeedsBootstrap`（users COUNT==0）替换 `WasFresh` 门；composition + testsupport 同步；Bootstrap 单事务 |
| C5–C8 前端 | **代码层有效** | runRequest/runBatchRequest catch→`REQUEST_FAILED`；handleSubmit try/finally；search 空串覆盖；permission 仅对已声明 target gate；`routeQuery`→`context.route` |
| D1–D8 | **有效（D3 residual）** | roles null 不清空；登录限流+dummy bcrypt；splitUrl safeDecode；快照毫秒格式；locale/theme `""`→auto；inputNumber 清空→undefined；上传 input 重置；ThemeToggle→setTheme（color-scheme） |
| 回归 | **可复现绿** | 本人：`go test ./...` 全包 ok（含 cmd/server 与 20+ 可测包）；`npx vitest run` **735** 全过；`npx tsc -b` 无错；HEAD=`9c1d0a7` |

## 对照成功标准

| 标准 | 独立结论 |
|------|----------|
| C1–C8 全部修复并回归 | **条件满足**：C2–C8 与 C1 下载面成立；C1「拒绝 SVG」与「DetectContentType 堵死主动内容」主张过满（F-001） |
| D1–D8 修复或 P-004 裁决 | **满足**（D3 accepted-residual 有书面记录） |
| 回归测试覆盖每项 + go/vitest/tsc 全绿 | **部分满足**：基线全绿已复现；C1 SVG / C4 NeedsBootstrap / C5–C8 / 前端 refresh 去重缺专项用例（F-002/F-004/F-005） |
| 共享基架重验证证据落盘、VP-008 go 恢复 | **未在本 scope 关门**（meta 成功标准仍 open；本意见不改 status） |

## Findings

### F-001 · required · 中高 · C1 SVG/主动内容「拒绝」面失效

- **文件:行号**：`apps/api/internal/handler/upload.go:40-44`（`dangerousInlineTypes` 仅 `text/html` / `application/xhtml+xml` / `image/svg+xml`）；`upload.go:127-135`（仅用 `http.DetectContentType` 比对）；对比 `upload.go:180-187`（真正兜底的 attachment/CSP/nosniff）。
- **失败场景**：常见 SVG 载荷**不会**被嗅探为 `image/svg+xml`，因而**不会**被拒绝并成功入库。本人用 `net/http.DetectContentType` 实测：
  - `<svg xmlns=...><script>…</script></svg>` → `text/plain; charset=utf-8`
  - `<?xml…?><svg…>` → `text/xml; charset=utf-8`
  - 前 512 字节填充后的 HTML → `text/plain`（绕过 HTML 拒绝）
  - `GIF89a` + HTML → `image/gif`（多格式夹带）
  E-001 / A-001 写「拒绝 HTML/**SVG**」与代码事实不符。同源导航 XSS 目前主要靠 **attachment + CSP sandbox + nosniff** 缓解，而非入库拒绝。若反向代理剥头、未来客户端用错误 MIME 做 blob/iframe 预览，主动内容仍在磁盘上。
- **确认度**：高（源码 + 本地 DetectContentType 探针 + 现有测试未覆盖 SVG）
- **建议闭合路径**：扩展拒绝（内容启发式 / 把 `text/xml`·可疑 SVG 标记拒掉 / 默认 allow-list 白名单），或书面 **accepted-residual** 明确「入库拒绝为 best-effort，安全边界=下载头」并补 SVG/多格式测试。

### F-002 · recommended · 低 · C1 回归测试未覆盖 SVG/嗅探旁路

- **文件:行号**：`apps/api/internal/handler/upload_test.go:116-185`（仅 HTML 误标 `text/plain` + 下载头断言）。
- **失败场景**：SVG/`text/xml`/填充 HTML 等可入库路径回归不会失败；与 F-001 配套。
- **确认度**：高

### F-003 · recommended · 中 · C2 跨标签页并发刷新仍可误清会话

- **文件:行号**：`apps/web/src/account/auth-client.ts:83-101`（`inflightRefresh` 为模块级单 Promise，**仅同一 JS realm**）；服务端 `auth.go:117-126` + `accounts.go:241-259` 已原子单胜。
- **失败场景**：两标签页同时 access 过期 → 各发 refresh → 服务端 1 胜 1 `ErrTokenRevoked` → 负方 `doRefresh` 失败 `clearTokens()` → 用户以为被登出。双会话签发已堵；原 C2「有效会话误登出」在跨标签仍可复现。
- **确认度**：高（代码路径核对；未做浏览器双标签 e2e）

### F-004 · recommended · 低 · C4 NeedsBootstrap 无专项测试

- **文件:行号**：`apps/api/internal/modules/authsession/systemdata/bootstrap.go:18-31`；`composition.go:112-121`。
- **失败场景**：grep 无 `NeedsBootstrap` 测试；仅靠间接系统测试。未来若改回 `WasFresh` 门，未必立刻红。
- **确认度**：高

### F-005 · recommended · 中 · C5–C8 无专项回归用例

- **文件:行号**：修复在 `apps/web/src/renderer/render.tsx`（约 278–304、339–348、634、761–772、1150–1181）、`apps/web/src/app/App.tsx`（约 420–422、517–526、714）等；仓库内无针对「网络失败不卡死 / 清空搜索 / 未声明权限放行 / route query 贯通」的独立 assert。
- **失败场景**：全量 vitest 735 绿只能证明无基线回归，不能证明这四项修复点被锁定；回退可能静默通过。
- **确认度**：高

### F-006 · recommended · 低 · D2 限流为进程内 best-effort（文档已知）

- **文件:行号**：`apps/api/internal/handler/rate_limit.go:10-62`；`auth.go:69-76`。
- **失败场景**：多实例不共享计数；不信任 `X-Forwarded-For` 时全站共用代理 IP 可能误伤；`attempts` map 无驱逐可在长期扫库下涨内存。对「单实例拖慢爆破」仍有效，不可当分布式 WAF。
- **确认度**：高（代码意图与注释一致；非阻塞关闭）

## 必改项汇总

| ID | 级别 | 一句话 |
|----|------|--------|
| F-001 | **required** | 纠正 C1「拒绝 SVG」假象：强化入库拒绝或书面 residual + 补测 |
| F-002～F-006 | recommended | 测试补齐 / 跨标签 refresh 协调 / 限流运维边界 |

无其他 required。D3 保持 accepted-residual（不重复开 finding）。

## 与既有意见的异同（A-001 self）

| 点 | A-001 self | A-002 independent |
|----|------------|-------------------|
| verdict | pass（待 independent） | **conditional** |
| C1 | ✅ 含 SVG 拒绝 | HTML+下载头 ✅；**SVG 拒绝 ❌**（F-001） |
| C2–C4, D1–D8 | ✅ | 同意有效（D3 residual；D2 边界 recommended） |
| C5–C8 | ✅（测全绿） | 代码同意；**缺专项测**（F-005） |
| 回归 go/vitest | 声称全绿 | **本人复跑确认**全绿 |
| required 开放 | 无 | **F-001 开放** |

## 结论 + 建议下一步

**verdict: conditional** — 不得无条件关门 / 不得据此单独恢复 VP-008 `go` 消费有效性，直至 F-001 按 fixed 或 accepted-residual/user-overruled 合法闭合。

建议编排器（`/govern`）：

1. 响应 F-001（优先：拒 SVG/`text/xml` 启发式或收紧 allow-list + 测试；或用户书面 residual 并改 E-001/成功标准措辞）。
2. 视风险处理 F-003（BroadcastChannel/localStorage 锁协调 refresh）与 F-005 专项测。
3. F-001 闭合且无新 required 后，再跑关门审计 / 重验证证据。

## 声明

本意见 `source: independent`，**不修改**目标 `status` / `progress` / goal-tree / 方案正文。响应与放行由 **`/govern`** 处理。
