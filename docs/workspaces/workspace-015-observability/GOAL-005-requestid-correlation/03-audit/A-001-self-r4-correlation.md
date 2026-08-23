---
id: GOAL-005-requestid-correlation
doc: audit-entry
record_id: A-001
source: self
verdict: pass
scope: R4 与 request-id 关联（I-005 闭合 + span 属性/baggage 注入）
created: 2026-08-22
updated: 2026-08-22
version: 1.0.0
parent: GOAL-001-observability
---

## A-001 · 自审：R4 request-id 关联（source: self）

- **日期**：2026-08-22
- **scope**：GOAL-005 全部交付物——D-001（I-005 闭合）、serverSpan 扩展、Wrap 接线、测试（checkpoint `8b52f2d` / `bc5e196`）
- **verdict**：**pass**（开放 required findings = 0）

### 核对成果

1. **I-005 闭合核对**：属性名冻结 `correlation.request_id`（D-001 §2）；baggage 键 `request-id` + 注入方式有实现与实测；关联判据（§1）由「真实 requestid 中间件链 + X-Request-ID → 属性相等」测试锁定。
2. **VP 退出判据 2 后半**：「能与现有 request-id / correlation 关联」现在有可核对测试证据；R5 只需把这判据纳入双路径证据。
3. **不变式保持**：metrics 标签白名单未动（R1 §4）；requestid 包零改动；HTTP 响应头行为未变；旁路语义（注入失败静默跳过 span/baggage）。
4. **回归**：全仓 `go test ./...` 无 FAIL；vet 干净。

### 偏差

无。实施范围与 D-001 一致；无出站调用所以 baggage 只做惰性就绪（由同 ctx 提取断言验证），符合 §2 边界。

### Findings

| 编号 | 级别 | 内容 | 状态 |
|------|------|------|------|
| N-007 | note | baggage 值受 RFC 字符限制；requestid 允许的 `:`/`.`/`_`/`-` 全部落在允许集内，但若未来 requestid 放宽字符集需复查 `NewMember` 校验 | open-note（不阻断，未来触发器） |
| N-008 | recommendation | R5 双路径证据的 trace 侧核对样例直接采用 D-001 §1 判据式（header id == span 属性）；N-004/N-006 继续有效 | open-note（指向 R5） |

### 结论

GOAL-005 四项成功标准全部满足且有证据链（D-001 → E-001/E-002 → 测试 → commit `8b52f2d`/`bc5e196`）。无未闭合 required finding；可关门（status: done, progress 4/4）。