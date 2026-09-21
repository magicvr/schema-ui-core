---
id: A-017-response-to-v0.7.0-preflight
goal_id: GOAL-001-nav-group-collapsible
doc: audit-entry
source: self
auditor: current-session
type: response
scope: response to A-016 and A-015 on PR #16 v0.7.0 release preflight
date: 2026-09-21
verdict: pass
created: 2026-09-21
updated: 2026-09-21
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# A-017 · 响应 v0.7.0 发布预检（2026-09-21）

- **source**：self（编排器响应）
- **auditor**：current-session
- **类型** / **scope**：response · 汇总 A-015 self 与 [A-016 independent](A-016-release-v0.7.0-preflight-independent.md)，核对 PR #16 合并前门禁并采纳合并后 F-001 的发布顺序
- **verdict**：**pass**（合并前 open required = 0；本记录不代表已合并或已发布）

## 范围与区间

工作区：`workspace-034-nav-group-collapsible`；目标：`GOAL-001-nav-group-collapsible`。对齐链、信息项与 Root 状态见本目标 meta、workspace 页面和 A-016。本响应只处理 PR #16 的 v0.7.0 preflight 意见，不改目标状态、progress 或 goal-tree。

## 成果与门禁

| 检查项 | 状态 | 证据 |
|---|---|---|
| A-015 F-S-001：内部 npm 包名错误映射 | **closed · fixed；A-016 独立确认** | 映射实现/测试、CI mapping step 与 dry-run 公开名；详见 `03-audit/A-015-release-v0.7.0-preflight-self.md`、`03-audit/A-016-release-v0.7.0-preflight-independent.md` |
| PR #16 合并前 required 门禁 | **满足预审；合并前仍须检查本次文档提交后的最新 head CI** | A-016 验证 head `33b07b68` 的 run `35570762993` 9/9 success；本响应与执行记录将新增审计事实，推送后 CI 必须在新 head 重新全绿 |
| A-016 F-001：npm 版本需先于 Go tag 可见 | **adopted · open until executed** | `schema-ui create` 当前钉到 protocol `0.2.16`、lib `0.1.15`、renderer `0.3.14`、ui `0.1.12`；在合并后先发布四包，逐包 `npm view` 确认 registry 可见并做消费校验，再创建 annotated `apps/api/v0.7.0` tag 与 GitHub Release |
| 外部发布动作 | **未发生** | 截至本响应，PR #16 仍 OPEN；四包尚未发布；tag 与 GitHub Release 尚未创建；详见 A-016 和 E-015 |

## 仍开放项

- A-016 F-001 保持 `recommended / med` 开放，直到四个新 npm 版本在 registry 可见并先于 tag 被验证；这是后续发布执行的顺序门禁，不阻断 PR 合并。
- 当前无开放 required finding，也无需用户对 residual / overrule 作裁决。

## 结论

接受 A-016 对合并前候选的 `pass` 与 A-015 包名修复 `fixed` 的独立确认。采纳 npm 先于 tag 的顺序。PR merge、main CI、npm 上传、tag 与 GitHub Release 均须待各自实际证据产生后再记为完成。本响应不提前放行尚未执行的合并后发布门禁。
