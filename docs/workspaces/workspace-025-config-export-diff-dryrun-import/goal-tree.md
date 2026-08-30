# 目标树 · workspace-025-config-export-diff-dryrun-import

## 树

```text
GOAL-001-config-export-diff-dryrun-import（active · 3/4 · primary_plan = VP-025-config-export-diff-dryrun-import）
├── GOAL-002-r1-contract-freeze（done · 3/3 · R1 合同冻结：配置包合同 v0.1.0）
├── GOAL-003-r2-export-diff（done · 3/3 · R2 export+diff：configpkg.go · 判据 #1/#2）
├── GOAL-004-r3-dryrun-import（done · 3/3 · R3 dry-run+import：方案 A · 判据 #3/#4）
└── 纲领阶段：R4 证据与关门（GOAL-005 候选 · D-004 设计）
```

## 状态表

| id | title | status | progress | parent |
|----|-------|--------|----------|--------|
| GOAL-001-config-export-diff-dryrun-import | 配置包导出 / diff / dry-run / 导入 | active | 3/4 | null |
| GOAL-002-r1-contract-freeze | R1 合同冻结 | done | 3/3 | GOAL-001-config-export-diff-dryrun-import |
| GOAL-003-r2-export-diff | R2 导出 + diff | done | 3/3 | GOAL-001-config-export-diff-dryrun-import |
| GOAL-004-r3-dryrun-import | R3 dry-run + 导入 | done | 3/3 | GOAL-001-config-export-diff-dryrun-import |

> 子目标按纲领阶段逐项立项后同步本树。progress 仅由 R1～R4 显式检查点等权计算；不放行阶段、不关闭 finding。