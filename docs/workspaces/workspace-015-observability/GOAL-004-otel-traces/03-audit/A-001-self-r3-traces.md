---
id: GOAL-004-otel-traces
doc: audit-entry
record_id: A-001
source: self
verdict: pass
scope: R3 OTel traces 接入（config 三键 + obs.Tracing + composition 接线 + live 冒烟）
created: 2026-08-21
updated: 2026-08-22
version: 1.0.0
parent: GOAL-001-observability
---

## A-001 · 自审：R3 traces 接入（source: self）

- **日期**：2026-08-22
- **scope**：GOAL-004 全部交付物——D-001（I-002 闭合）、config 三键与校验、obs.Tracing、composition 接线、测试与 live 冒烟（checkpoint `0470307` / `2ab4ec4`）
- **verdict**：**pass**（开放 required findings = 0）

### 核对成果

1. **I-002 闭合核对**：协议（OTLP/HTTP）——实现用 `otlptracehttp`，无 gRPC；采样（ParentBased+ratio，缺省 1.0）——`WithSampler` 实证；no-op（未配置 endpoint 不得挡 mvp/dev）——disabled 纯 no-op 单测 + config fail-closed + live 冒烟负路径（不可达 endpoint 启动/服务正常）。
2. **合同符合性（对照 D-001）**：三键名/env/默认值逐条一致；§5 attrs 白名单（method/route/url.path/status）与 span name 有测试断言；≥500 → codes.Error 有断言；W3C traceparent join 有断言；§3「导出失败仅告警不致命」由 live 冒烟实证（2 次 WARN + 服务全程 200）。
3. **不变式保持**：R2 metrics 面行为零回归（全仓 `go test ./...` 无 FAIL；metrics 系列断言仍在）；no-op 缺省零开销；traces 不进 readyz（与 readiness gate 无耦合）；Owning 拦截点单一（未引入第二套 instrumentation 源）。
4. **VP 对齐**：退出判据 2 的「OTLP 可导出 + HTTP 至少可出 span」已由 sink 实证 + span 形状测试满足；「未配置 endpoint 不得挡 mvp/dev」由 live 冒烟满足。

### 偏差

无。实施范围与 D-001 一致；collector 鉴权 headers 按 §2 记录为「后续增量」，未加。

### Findings

| 编号 | 级别 | 内容 | 状态 |
|------|------|------|------|
| N-005 | note | `main.go` 的 fx 日志与 slog JSON 分别落 stderr/stdout——冒烟取证两者都要看；无涉及变更 | open-note（不阻断） |
| N-006 | recommendation | R4 关联入 span 时，建议用 W3C `baggage` 传递 `request-id` 的 key 名（D-001 已预留提取），并在 R4 决策里冻结属性名 `correlation.request_id`；N-004（R5 固化冒烟脚本）继续有效 | open-note（指向 R4/R5） |

### 结论

GOAL-004 四项成功标准全部满足且有证据链（D-001 → E-001/E-002 → 测试/live 冒烟 → commit `0470307`/`2ab4ec4`）。无未闭合 required finding；可关门（status: done, progress 4/4）。N-006 作为输入带入 R4 立项。