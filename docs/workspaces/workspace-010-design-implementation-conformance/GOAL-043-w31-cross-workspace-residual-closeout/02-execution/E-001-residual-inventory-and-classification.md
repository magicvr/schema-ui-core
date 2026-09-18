---
id: E-001-residual-inventory-and-classification
doc: execution-entry
status: recorded
goal_id: GOAL-043-w31-cross-workspace-residual-closeout
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-design-implementation-conformance
version: 1.0.0
---

# E-001 · 残余清点与分类（C1）

## 事实

2026-09-18，扫描 workspace-037 关门后的全部开放项与 gated 清单，得到下表（`I-043-001` verified）。

| # | 编号 | 它到底在说什么 | 性质 | 影响 | 触发 / 依据 |
|---|------|----------------|------|------|-------------|
| 1 | `GOAL-009 A-001 F-001`（low·recommended） | 列表视觉 e2e 只断言「`.dark` 下 `--control` 存在且≠浅色值」，未断言暗色下**开关的计算背景**等于该值——浅色路径已钉住 token 驱动关系，暗色路径无独立断言 | 已交付范围内的测试覆盖残余 | 低 | 无（可直接修复） |
| 2 | `GOAL-009 A-001 F-002`（low·recommended） | 浏览器守卫只跑 `/roles`（唯一同时具备 search form 与 toolbar 的页面），覆盖面窄于「通用列表页」的合同表述 | 同上 | 低 | 无（可直接修复） |
| 3 | `GOAL-005 A-002 F-002` / `R5-I-005`（low·recommended / non-blocking deferred） | R4 已补 maintenance、timeout/offline 与资源读取回归，但 **Host 终态（`HostFailureScreen`）与普通 resource 反馈没有直接对照断言**——两张分类表各自被测，冗余文案或单侧失去显式处理不会被发现 | 同上 | 低 | 无（可直接修复） |
| 4 | `GOAL-008 A-002 F-002`（low·recommended） | 全仓 `tsc` 简写未逐条裁定：354 行提到 `tsc` 而形态不可从文本唯一确定（45 行指 `npm run build`（= `tsc -b && vite build`，形态确定）、40 行字面 `tsc --noEmit`（已更正或已加勘误）、269 行叙述式如「tsc clean / tsc 0」） | 文档证据形态残余 | 无（可执行面已被守卫锁死） | 历史记录被再次引用为类型检查证据时 |
| 5 | `V-F124`（vision·recommended） | 要求在激活/R1 前把首波页面分母、状态字段、Profile/权限覆盖、Saved View 持久化边界、dirty-state/反馈类型做成机器可核对矩阵，**以免执行阶段滑向共享视图、实体搜索或第二套基础设施** | 愿景层记账（实质已由 R1 交付） | 无 | `/vision` 复核闭合 |
| 6 | `I-037-005`（non-blocking·deferred） | 跨用户共享视图、最近使用/收藏、协作权限是否进入后续波次 | 范围边界型悬置决策（非缺陷） | 无 | 真实协作需求出现 → `/vision` |
| 7 | 实体全文检索 / `RT-X01`/`RT-X02` | `RT-X01` = 专用搜索引擎路线、`RT-X02` = DB 全文检索路线（架构触发项） | 未推进 / trigger-gated 能力 | — | 实体级搜索需求 + 规模证据 → `/vision` |
| 8 | 批量结果中心 | roadmap「体验增强」清单里的未立项项 | 同上 | — | 新立 VP |
| 9 | 组织·部门·岗位 + `org` 数据权限 | 基架能力剩余 #2 | 同上（2026-08-29 用户书面降权） | — | 多组织/多团队 fork 消费或真实多组织管理需求 |
| 10 | 新业务域 | 业务域分支 | 同上 | — | 真实业务需求（无触发不预开第二域） |
| 11 | Redis / MQ / 多实例（+ 第二持久化栈） | 架构 A3 / `RT-Q03`，架构骨架唯一未触发项 | 同上 | — | 多实例部署或 C 端业务域模块正式接入同进程 |

## 分类结论

- **可就地修复（1～3）**：都是"已交付 + 已验证过，但测试未固化"的覆盖面问题，改动限于测试与为测试可观测性所需的最小导出。
- **量化收口（4）**：可执行面已由 `GOAL-008`/`GOAL-010` 的守卫与 CI 门禁锁死；文档侧按 bounded residual 处理并登记触发条件，不做逐条考古。
- **提请愿景层闭合（5）**：实质要求已由 R1 交付，属记账未回填；由 `/vision` 判断，本目标只提供证据。
- **统一登记、暂不推进（6～11）**：能力/范围类，保留各自触发条件与责任人。

C1 完成。C2 见 `E-002`/`E-003`。
