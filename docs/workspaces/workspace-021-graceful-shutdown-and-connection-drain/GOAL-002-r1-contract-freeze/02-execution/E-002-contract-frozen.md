---
id: E-002
title: 信息裁决落盘 + 合同 v0.1.0 冻结
date: 2026-08-27
status: done
---

# E-002 · 合同冻结（2026-08-27）

## 事实

1. **C1 信息裁决**：用户（界面裁决）三条 required 全部采纳建议——I-001 中断标记重跑；I-002 默认 10s + `http.shutdown_timeout`/`HTTP_SHUTDOWN_TIMEOUT`（非法值 fail-closed）；I-003 fail-closed 启动期校验（无运行时迁移窗口）。I-004（non-blocking）lead 口径：结构化日志断言，指标不进分母。`D-001` → accepted；四项 → `verified`。
2. **C2 合同冻结**：`01-decision/D-002-contract-freeze.md` 落盘 **优雅停机 / 连接排空合同 v0.1.0**——§0 适用基线 / §1 强制停机顺序 / §2 HTTP drain / §3 退出码 / §4 Job 中断标记重跑 / §5 Store 排空 × 迁移 fail-closed / §6 配置键与默认 / §7 可观测 / §8 R3 验收方式 / 未选方案。未选方案留痕（等完成、非法值回落、运行时迁移排队、指标进分母、分面子预算）。
3. R1 责任文件 = D-002；实施（R2）与验收（R3）以该文件为分母。

## 验证 / 后续

- C3 自审（A-001）与关门见 `03-audit/` 与 00-meta；Root 信息台账回写（E-003）。