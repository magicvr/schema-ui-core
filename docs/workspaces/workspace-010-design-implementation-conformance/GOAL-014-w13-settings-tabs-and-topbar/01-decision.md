---
id: GOAL-014-w13-settings-tabs-and-topbar
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# 决策记录 · GOAL-014

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | non-blocking | 设置页功能单元切分与重置按钮归属 | S2 T-01 | S2 | 用户指令 + as-built 对照 | **verified** | — | D-001 §T-01 |
| I-002 | non-blocking | 移动端断点 | S2 T-02 | S2 | 既有壳层断点约定 | **verified** | — | D-001 §T-02 |
| I-003 | non-blocking | 头像存储复用与清理模型 | S2 T-05 | S2 | 既有 W9 资产存储 as-built | **verified** | — | D-002 §T-05：共享 RasterAssetStore；替换/清空删除旧文件；无启动 GC |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-16 | 开波：四项交互整改纳入本波；设计冻结 | accepted（范围+设计） | `01-decision/D-001-w13-freeze.md` |
| D-002 | 2026-08-16 | 追加：移动端汉堡靠左修正 + T-05 头像上传设计 | accepted（范围+设计） | `01-decision/D-002-w13-followup-avatar.md` |
| D-003 | 2026-08-16 | 追加：T-06 通知中心交互修正设计（点击即读 + 展开详情 + 未读数即时刷新） | accepted（范围+设计） | `01-decision/D-003-t06-notification-interactions.md` |
| D-004 | 2026-08-16 | 追加：T-07 列表筛选即时生效设计（即时控件 vs 提交式文本框；chips 以已提交查询为真相源） | accepted（范围+设计） | `01-decision/D-004-t07-live-filters.md` |

## 待决问题（P-004）

无。追加指令即范围确认；筛选即时生效为渲染层整改（无协议/schema/Go 变更），无冲突意见、无开放 required 信息项。
