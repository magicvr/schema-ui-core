---
id: GOAL-010-w9-branding-asset-upload
doc: decision
status: active
parent: GOAL-001-design-implementation-conformance
created: 2026-08-15
updated: 2026-08-15
version: 0.1.0
---

# 决策记录 · GOAL-010

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | URL 输入彻底移除；旧 URL 值仅兼容读取不迁移 | 方案 | S1 | 用户 2026-08-15 书面确认 | closed | — | D-001 |
| I-002 | required | 处理参数（512px/64px/q82/4MiB）+ x/image 依赖 | 方案 | S1 | 用户 2026-08-15 书面确认 | closed | — | D-001 |
| I-003 | required | SVG 不支持（安全拒收） | 方案 | S1 | 用户 2026-08-15 书面确认 | closed | — | D-001 |
| I-004 | required | 清理语义（替换/重置即删 + 启动 GC 孤儿） | 方案 | S1 | 用户 2026-08-15 书面确认 | closed | — | D-001 |
| I-005 | required | 参数进 config.yaml（env 覆盖） | 方案 | S1 | 用户 2026-08-15 书面确认 | closed | — | D-001 |
| I-006 | required | 关门 cross 审计（self + independent） | 关门 | S6 | 用户 2026-08-15 书面确认；provider 执行时确认 | closed | provider 待 S6 | D-001 |
| I-007 | required | 消费点完整清单 | 实施 | S4 | 创建时代码扫描（D-001 附录 A） | closed | — | D-001 |
| I-008 | required | UploadField 移除已设图片交互能力 | 实施 | S4 | 检查/扩展 form-controls.tsx | closed | — | E-002（扩展 UploadField 移除按钮） |
| I-009 | non-blocking | 公开 GET 缓存/安全头细节 | 实施 | S2 | D-001 设计内定（沿用上传仓安全基线） | closed | — | E-002/E-003（nosniff + CSP sandbox + immutable 实测） |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-15 | 品牌图标上传方案冻结（六项用户裁决） | accepted | `01-decision/D-001-brand-asset-upload.md` |
