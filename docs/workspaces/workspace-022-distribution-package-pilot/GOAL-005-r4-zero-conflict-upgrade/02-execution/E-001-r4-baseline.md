---
status: active
created: 2026-08-29
updated: 2026-08-29
parent: GOAL-005-r4-zero-conflict-upgrade
version: 0.1.0
---

# E-001 · S1 样本设计与 V1 基线冻结（2026-08-29）

## 样本定案（D-001）

| 样本 | 变更 | 类型 |
|------|------|------|
| A1 | `kernel.JoinKeys`（贡献键拼接） | A 层 additive |
| E1 | `normalizePageID`（页面 id 规范化） | 协议面 additive |
| M1 | 新迁移（site_settings updated_at 索引） | 全局台账（双方言） |

选型理由（I-004）：kernel/protocol additive = 「经版本化可平滑」的代表；新迁移 = 判据 #3 显式要求且是 fork-merge 最痛面（台账冲突）——二者组合最真实代表升级压力。

## V1 基线（commit `8686b3fd` · 62 条迁移）

golden-consumer `go run`（SQLite fresh）：

```
kernel=2.0.0 profile=admin dialect=sqlite fresh=true contrib{routes=10 pages=2 perms=3 nav=1 frag=1} users_module=admin.users
```

golden-web 三探针（probe / probe-render / token-check）：PASS（记录于 GOAL-004 E-002/E-003 基线输出）。