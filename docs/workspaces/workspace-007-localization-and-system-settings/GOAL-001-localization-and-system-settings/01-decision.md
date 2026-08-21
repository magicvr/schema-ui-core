---
id: GOAL-001-localization-and-system-settings
doc: decision
status: done
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.3.0
---

# 决策记录 · GOAL-001

## 信息需求与阶段门禁

> 本文件是稳定索引。信息台账（P-005）与 [VP-007](../../../vision/plans/VP-007-localization-and-system-settings.md) `I-L10N-001`～`005` 同 id 对齐；`required` 项未在其“最晚需要阶段”前 verified 或获合规 `accepted-residual` 时阻断对应门禁（S0 冻结 / S1 / S3 / S4）。

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| `I-L10N-001` | required | Schema 驱动页面的字段标签、说明、动作和服务端文档如何本地化，且不创建宽于或冲突于 `schema-ui-docs@v2.7.0` 的私有协议语义 | 多语种方案冻结（S0→S1） | 多语种方案冻结前 | 盘点 v2.7.0 可用 key/text 字段与当前 Renderer；比较前端 key 解析、服务端 locale overlay、上游提案等路径并冻结兼容策略 | **verified** | — | 2026-08-09 用户书面选定「前端 key 解析」；盘点证据见 E-002，冻结见 D-002 §I-L10N-001 |
| `I-L10N-002` | required | 用户显式语种选择首版持久化在浏览器本地还是账号资料；匿名到登录后的合并规则 | 用户控制实施（S1） | 用户控制实施前 | 结合现有 auth/session 与本地 theme 机制形成优先级和迁移矩阵；用户确认后冻结 | **verified** | — | 2026-08-09 用户书面选定「localStorage 单通道」；优先级/合并规则冻结见 D-002 §I-L10N-002 |
| `I-L10N-003` | required | 公开品牌/locale 启动配置是兼容扩展 `/api/branding`，还是建立新的公开 bootstrap 契约；缓存与配置刷新语义如何保持一致 | Settings API 方案冻结（S3） | Settings API 方案冻结前 | 对当前 branding 消费端、Profile 路由和缓存头做差量盘点，冻结兼容路径 | **verified** | — | 2026-08-09 用户书面选定「兼容扩展 `/api/branding`」；字段与缓存语义冻结见 D-002 §I-L10N-003 |
| `I-L10N-004` | required | 当前 `{error,message}` 错误 envelope、重复 `writeError` 与前端直显链路扩展到 locale 协商的真实成本和兼容边界 | 后端提示本地化实施（S4）与 exit 5 | 后端提示本地化实施前 | 盘点用户可见错误码及调用点；用认证/验证/设置错误验证 `Accept-Language`、`Content-Language`、key/params 与 fallback；关闭结论必须选 exit 5 路径 (a) 实施证据或 (b) 用户书面 `accepted-residual`，禁止仅写“成本无界” | **verified（路径 a）** | — | 2026-08-09 用户书面选定 exit 5 路径 (a) 有界服务端 locale 协商；envelope 扩展与编目范围冻结见 D-002 §I-L10N-004。**本 verified ≠ VP-007 exit 5 关闭**：exit 5 证据 = S4 实施（协商/Content-Language/失败回退），见 A-001 F-003 |
| `I-L10N-005` | required | 默认时区的存储、展示和服务器时间语义如何定义，避免把显示时区与持久化时间混为一谈 | Localization 设置实施（S3） | Localization 设置实施前 | 固定 UTC 存储、显示转换、`auto`/指定时区和无效时区失败语义 | **verified** | — | 2026-08-09 用户书面选定「UTC 存储 + 显示转换」；`siteTimezone` 语义冻结见 D-002 §I-L10N-005 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|------|------|------|------|------|
| D-001 | 2026-08-09 | S0 · scaffold 工作区与 Root、纲领路线图 S0–S5、信息门禁登记 | accepted | `01-decision/D-001-scaffold-roadmap-info-gates.md` |
| D-002 | 2026-08-09 | S0 · 差距盘点与契约冻结：关闭 I-L10N-001～005、冻结 F-V029 覆盖表 | accepted | `01-decision/D-002-s0-contract-freeze-info-gates.md` |
| D-003 | 2026-08-09 | 暂时回退 Root 关门状态承接 S6 子目标（GOAL-007） | accepted | `01-decision/D-003-reopen-root-for-s6.md` |
| D-004 | 2026-08-09 | S6 关门，Root 恢复 done（解除临时回退，`7/7`） | accepted | `01-decision/D-004-s6-closeout-root-restored.md` |
