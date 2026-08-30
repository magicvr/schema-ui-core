---
status: accepted
created: 2026-08-30
updated: 2026-08-30
parent: GOAL-001-config-export-diff-dryrun-import
version: 0.1.0
---

# D-004 · R4 关门设计（2026-08-30 · 前置于 R3 完成）

R4 = 证据与关门（VP-025 判据 #6 · Root 纲领 R4）。本条目冻结关门路径（继 R1/R2/R3 完成之后执行；R3 C2 import 仍需 I-025-004 裁决）：

1. **前置**：R3 全关（GOAL-004 done 3/3 · 判据 #3/#4 交付）。
2. **证据矩阵**：Root `attachments/r4-evidence-matrix.md`——六条判据 ↔ 阶段证据链接（GOAL-002 合同 / GOAL-003 export+diff / GOAL-004 dry-run+import / 快测与冒烟记录 / 全量回归），每项标 verified 或缺漏。
3. **越界核账**：红线清单逐项核对（Profile 默认集 / 模块矩阵 / Manifest 装配零触碰 · 迁移台账零变更 · 密钥 fail-closed · 热加载未引入 · 管理面未做）；`git log --stat` 佐证。
4. **关门审计（双审）**：
   - **A-001 self**：全链一致性（合同 ↔ 实现 ↔ 判据 ↔ 信息台账）· verdict pass 基线；
   - **A-002 independent · grok build**（项目级路径 `docs/architecture/independent-audit-execution.md`：grok 4.6 · 思考强度 high · 直接 `/audit` 至本 Root）：`source: independent` 落盘 `03-audit/A-002-*.md`；grok 不可用/无可核对输出时 independent 门禁保持未满足，不得冒充。
   - 合并响应全部意见（P-003）；required 三路径闭合；冲突按 P-004 由用户裁决。
5. **vision 层关门**：VRev-055（self · `/vision` 或按需 independent）覆盖 VP-025 关门前审视（退出判据六条 · 区证据链 · 无开放 VRev required）；**VP-025 `active → closed`（v0.3.0）须用户书面确认**；残余（若有）按有界口径点名 workspace/goal 并登记 Root `D-001` 或 GOAL 内。
6. **同步**：VP-025 关门记录 + roadmap/workspaces 行同步 + goal-tree/workspace.md + checkpoint commit。
7. **不进入**：不改 Charter；不重开 VP-007/023/024；不把配置包操作化外溢为配置中心（架构分支无触发）。

**门槛声明**：以上设计不替代 I-025-004 裁决；R3 C2 未完成前 R4 不启动。