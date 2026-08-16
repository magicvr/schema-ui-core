---
id: GOAL-013-w12-product-surface-intent
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.8.0
---

# 决策记录 · GOAL-013

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 各列表页搜索/筛选字段矩阵 | S2 T-02 | S2 | 用户选择题 + 对照 | **verified** | — | D-003 |
| I-002 | required | 个人中心选项卡信息架构 | S2 T-03 | S2 | 用户选择题 | **verified** | — | D-004 |
| I-003 | required | 「我的钱包」工作区归属与自服务范围 | S2 T-04 | S2 | 用户书面 | **verified** | — | D-005 → GOAL-022 |
| I-004 | required | T-06 是配置面卫生还是 Profile 语义变更 | S2 T-06 | S2 | 用户书面 | **verified** | — | D-007 |
| I-006 | required | 废除 env/.env 的范围（仅模块 vs 全局含密钥） | S3 T-06 | S3 改加载器前 | 用户选择题 | **verified** | — | D-008 |
| I-005 | non-blocking | 顶栏下拉触发器与移动端抽屉并存策略 | S2 T-01 | S2 | 用户选择题 | **verified** | — | D-002 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-16 | 开波：六条意图纳入本波范围；设计未冻结 | accepted（范围） | `01-decision/D-001-intent-inventory.md` |
| D-002 | 2026-08-16 | T-01 顶栏用户区下拉（头像+姓名；抽屉不重复） | accepted | `01-decision/D-002-t01-user-menu.md` |
| D-003 | 2026-08-16 | T-02 列表搜索字段矩阵（关键词改名 + Extra 筛选） | accepted | `01-decision/D-003-t02-list-search.md` |
| D-004 | 2026-08-16 | T-03 个人中心三档 Tabs（资料 / 安全 / 会话） | accepted | `01-decision/D-004-t03-account-tabs.md` |
| D-005 | 2026-08-16 | T-04 移交 workspace-011 GOAL-022；本波不做 | accepted | `01-decision/D-005-t04-handoff-w011.md` |
| D-006 | 2026-08-16 | T-05 回收站时间改为 API 输出 ISO-8601 | accepted | `01-decision/D-006-t05-recycle-time.md` |
| D-007 | 2026-08-16 | T-06 模块启用只认 config.yaml（预设文件或内联） | accepted | `01-decision/D-007-t06-yaml-modules.md` |
| D-008 | 2026-08-16 | I-006 只取消模块启用 env；S3 分批冻结 | accepted | `01-decision/D-008-i006-env-scope.md` |

## 待决问题（P-004）

无。S2 已冻结。实施前无开放 required 信息项。
